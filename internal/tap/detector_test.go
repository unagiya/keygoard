package tap

import (
	"testing"

	"github.com/unagiya/keygoard/keycode"
)

func TestPress(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	if d.GetPhase(0, 0) != Pending {
		t.Errorf("Phase after Press = %d, want Pending", d.GetPhase(0, 0))
	}
}

func TestReleaseImmediately(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	wasPending, got := d.Release(0, 0)

	if !wasPending {
		t.Error("Release immediately should be tap (wasPending=true)")
	}
	if got != kc {
		t.Errorf("Release keycode = %#x, want %#x", got, kc)
	}
	if d.GetPhase(0, 0) != Idle {
		t.Errorf("Phase after Release = %d, want Idle", d.GetPhase(0, 0))
	}
}

func TestReleaseBelowThreshold(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	// 閾値未満まで Advance
	for i := 0; i < Threshold-1; i++ {
		d.Advance()
	}

	wasPending, got := d.Release(0, 0)
	if !wasPending {
		t.Error("Release below threshold should be tap")
	}
	if got != kc {
		t.Errorf("Release keycode = %#x, want %#x", got, kc)
	}
}

func TestHoldExceedsThreshold(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	for i := 0; i < Threshold; i++ {
		d.Advance()
	}

	timedOut, got := d.CheckTimeout(0, 0)
	if !timedOut {
		t.Error("CheckTimeout should return true after threshold")
	}
	if got != kc {
		t.Errorf("CheckTimeout keycode = %#x, want %#x", got, kc)
	}
	if d.GetPhase(0, 0) != Holding {
		t.Errorf("Phase after timeout = %d, want Holding", d.GetPhase(0, 0))
	}
}

func TestHoldRelease(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	for i := 0; i < Threshold; i++ {
		d.Advance()
	}
	d.CheckTimeout(0, 0)

	// ホールド中のリリースは wasPending=false
	wasPending, _ := d.Release(0, 0)
	if wasPending {
		t.Error("Release after hold should not be tap")
	}
	if d.GetPhase(0, 0) != Idle {
		t.Errorf("Phase after hold release = %d, want Idle", d.GetPhase(0, 0))
	}
}

func TestCheckTimeoutBeforeThreshold(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	for i := 0; i < Threshold-1; i++ {
		d.Advance()
	}

	timedOut, _ := d.CheckTimeout(0, 0)
	if timedOut {
		t.Error("CheckTimeout should return false before threshold")
	}
	if d.GetPhase(0, 0) != Pending {
		t.Errorf("Phase should still be Pending")
	}
}

func TestMultipleKeysIndependent(t *testing.T) {
	var d Detector
	kc1 := keycode.LT(1, keycode.A)
	kc2 := keycode.TT(2)

	d.Press(0, 0, kc1)
	// 100 ticks 後に 2 つ目を押下
	for i := 0; i < 100; i++ {
		d.Advance()
	}
	d.Press(1, 0, kc2)

	// さらに 100 ticks（kc1 は合計 200、kc2 は 100）
	for i := 0; i < 100; i++ {
		d.Advance()
	}

	// kc1 はタイムアウト
	timedOut1, _ := d.CheckTimeout(0, 0)
	if !timedOut1 {
		t.Error("key (0,0) should have timed out")
	}

	// kc2 はまだ pending
	timedOut2, _ := d.CheckTimeout(1, 0)
	if timedOut2 {
		t.Error("key (1,0) should not have timed out yet")
	}
}

func TestReset(t *testing.T) {
	var d Detector
	kc := keycode.LT(1, keycode.A)

	d.Press(0, 0, kc)
	d.Advance()
	d.Reset(0, 0)

	if d.GetPhase(0, 0) != Idle {
		t.Errorf("Phase after Reset = %d, want Idle", d.GetPhase(0, 0))
	}
}

func TestAdvanceOnlyAffectsPending(t *testing.T) {
	var d Detector

	// idle 状態のキーは Advance で影響を受けない
	d.Advance()
	timedOut, _ := d.CheckTimeout(0, 0)
	if timedOut {
		t.Error("idle key should not time out")
	}
}

func TestTTKeycode(t *testing.T) {
	var d Detector
	kc := keycode.TT(1)

	d.Press(0, 0, kc)
	wasPending, got := d.Release(0, 0)

	if !wasPending {
		t.Error("TT tap should be wasPending=true")
	}
	if got != kc {
		t.Errorf("TT Release keycode = %#x, want %#x", got, kc)
	}
}
