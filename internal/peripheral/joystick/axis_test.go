package joystick

import "testing"

func TestMapDelta_Center(t *testing.T) {
	m := NewAxisMapper(3000, 5, false)
	if got := m.MapDelta(32768); got != 0 {
		t.Errorf("center: got %d, want 0", got)
	}
}

func TestMapDelta_DeadZone(t *testing.T) {
	m := NewAxisMapper(3000, 5, false)

	// デッドゾーン内の値は 0
	tests := []uint16{
		32768 - 2999, // center - (deadZone-1)
		32768 + 2999, // center + (deadZone-1)
		32768 - 1,
		32768 + 1,
	}
	for _, raw := range tests {
		if got := m.MapDelta(raw); got != 0 {
			t.Errorf("deadZone raw=%d: got %d, want 0", raw, got)
		}
	}
}

func TestMapDelta_JustOutsideDeadZone(t *testing.T) {
	m := NewAxisMapper(3000, 5, false)

	// デッドゾーン境界直後: effective=1 では int 切り捨てで 0 になるため
	// 十分な偏差（activeRange/maxOutput 以上）が必要
	// activeRange=29768, maxOutput=15 → 最小有効偏差 ≈ 1985
	if got := m.MapDelta(32768 + 5000); got == 0 {
		t.Error("outside deadZone (positive): got 0, want non-zero")
	}
	if got := m.MapDelta(32768 - 5000); got == 0 {
		t.Error("outside deadZone (negative): got 0, want non-zero")
	}
}

func TestMapDelta_MaxTilt(t *testing.T) {
	m := NewAxisMapper(3000, 5, false)

	// 最大傾倒（65535）: deviation=32767, effective=29767, 29767*15/29768=14
	// uint16 の非対称性（正方向 max は 32767）により maxOutput-1 になる
	got := m.MapDelta(65535)
	if got != 14 {
		t.Errorf("max tilt: got %d, want 14", got)
	}

	// 最小傾倒（0）: deviation=-32768, effective=-29768, -29768*15/29768=-15
	got = m.MapDelta(0)
	if got != -15 {
		t.Errorf("min tilt: got %d, want -15", got)
	}
}

func TestMapDelta_Invert(t *testing.T) {
	m := NewAxisMapper(3000, 5, true)

	// 反転: 正方向傾倒 → 負の出力（14 → -14）
	got := m.MapDelta(65535)
	if got != -14 {
		t.Errorf("invert max tilt: got %d, want -14", got)
	}

	// 反転: 負方向傾倒 → 正の出力（-15 → 15）
	got = m.MapDelta(0)
	if got != 15 {
		t.Errorf("invert min tilt: got %d, want 15", got)
	}
}

func TestMapDelta_Sensitivity(t *testing.T) {
	// 感度 1: maxOutput = 3, 最大正方向 = 29767*3/29768 = 2
	m1 := NewAxisMapper(3000, 1, false)
	got1 := m1.MapDelta(65535)
	if got1 != 2 {
		t.Errorf("sensitivity 1: got %d, want 2", got1)
	}

	// 感度 10: maxOutput = 30, 最大正方向 = 29767*30/29768 = 29
	m10 := NewAxisMapper(3000, 10, false)
	got10 := m10.MapDelta(65535)
	if got10 != 29 {
		t.Errorf("sensitivity 10: got %d, want 29", got10)
	}

	// 負方向は正確に maxOutput に到達する
	gotNeg1 := m1.MapDelta(0)
	if gotNeg1 != -3 {
		t.Errorf("sensitivity 1 neg: got %d, want -3", gotNeg1)
	}
	gotNeg10 := m10.MapDelta(0)
	if gotNeg10 != -30 {
		t.Errorf("sensitivity 10 neg: got %d, want -30", gotNeg10)
	}
}

func TestMapDelta_SensitivityClamp(t *testing.T) {
	// 感度 0 → 1 にクランプ（maxOutput=3、負方向で検証）
	m0 := NewAxisMapper(3000, 0, false)
	got0 := m0.MapDelta(0)
	if got0 != -3 {
		t.Errorf("sensitivity 0 (clamped to 1): got %d, want -3", got0)
	}

	// 感度 255 → 10 にクランプ（maxOutput=30、負方向で検証）
	m255 := NewAxisMapper(3000, 255, false)
	got255 := m255.MapDelta(0)
	if got255 != -30 {
		t.Errorf("sensitivity 255 (clamped to 10): got %d, want -30", got255)
	}
}

func TestMapDelta_ZeroDeadZone(t *testing.T) {
	m := NewAxisMapper(0, 5, false)

	// デッドゾーン 0: 中心は 0
	if got := m.MapDelta(32768); got != 0 {
		t.Errorf("zero deadZone center: got %d, want 0", got)
	}

	// 最大正方向: 32767*15/32768 = 14
	got := m.MapDelta(65535)
	if got != 14 {
		t.Errorf("zero deadZone max positive: got %d, want 14", got)
	}

	// 最大負方向: -32768*15/32768 = -15
	gotNeg := m.MapDelta(0)
	if gotNeg != -15 {
		t.Errorf("zero deadZone max negative: got %d, want -15", gotNeg)
	}
}

func TestMapDelta_LargeDeadZone(t *testing.T) {
	// デッドゾーンが中心値と同じ → activeRange = 0 → 常に 0
	m := NewAxisMapper(32768, 5, false)
	if got := m.MapDelta(0); got != 0 {
		t.Errorf("deadZone >= center: got %d, want 0", got)
	}
	if got := m.MapDelta(65535); got != 0 {
		t.Errorf("deadZone >= center: got %d, want 0", got)
	}
}

func TestMapDelta_Symmetry(t *testing.T) {
	m := NewAxisMapper(3000, 5, false)

	// 中心から等距離の正負は符号が逆で絶対値が等しい
	pos := m.MapDelta(32768 + 10000)
	neg := m.MapDelta(32768 - 10000)
	if pos != -neg {
		t.Errorf("symmetry: pos=%d, neg=%d, want pos == -neg", pos, neg)
	}
}
