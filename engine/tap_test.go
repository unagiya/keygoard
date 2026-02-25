package engine

import (
	"testing"

	"github.com/unagiya/keygoard/keycode"
)

func TestTapPress(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	if td.Phase(0, 0) != tapPending {
		t.Errorf("Phase after Press = %d, want tapPending", td.Phase(0, 0))
	}
}

func TestTapReleaseImmediately(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	wasPending, got := td.Release(0, 0)

	if !wasPending {
		t.Error("Release immediately should be tap (wasPending=true)")
	}
	if got != kc {
		t.Errorf("Release keycode = %#x, want %#x", got, kc)
	}
	if td.Phase(0, 0) != tapIdle {
		t.Errorf("Phase after Release = %d, want tapIdle", td.Phase(0, 0))
	}
}

func TestTapReleaseBelowThreshold(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	// 閾値未満まで Advance
	for i := 0; i < tapThreshold-1; i++ {
		td.Advance()
	}

	wasPending, got := td.Release(0, 0)
	if !wasPending {
		t.Error("Release below threshold should be tap")
	}
	if got != kc {
		t.Errorf("Release keycode = %#x, want %#x", got, kc)
	}
}

func TestHoldExceedsThreshold(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	for i := 0; i < tapThreshold; i++ {
		td.Advance()
	}

	timedOut, got := td.CheckTimeout(0, 0)
	if !timedOut {
		t.Error("CheckTimeout should return true after threshold")
	}
	if got != kc {
		t.Errorf("CheckTimeout keycode = %#x, want %#x", got, kc)
	}
	if td.Phase(0, 0) != tapHolding {
		t.Errorf("Phase after timeout = %d, want tapHolding", td.Phase(0, 0))
	}
}

func TestHoldRelease(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	for i := 0; i < tapThreshold; i++ {
		td.Advance()
	}
	td.CheckTimeout(0, 0)

	// ホールド中のリリースは wasPending=false
	wasPending, _ := td.Release(0, 0)
	if wasPending {
		t.Error("Release after hold should not be tap")
	}
	if td.Phase(0, 0) != tapIdle {
		t.Errorf("Phase after hold release = %d, want tapIdle", td.Phase(0, 0))
	}
}

func TestCheckTimeoutBeforeThreshold(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	for i := 0; i < tapThreshold-1; i++ {
		td.Advance()
	}

	timedOut, _ := td.CheckTimeout(0, 0)
	if timedOut {
		t.Error("CheckTimeout should return false before threshold")
	}
	if td.Phase(0, 0) != tapPending {
		t.Errorf("Phase should still be tapPending")
	}
}

func TestMultipleKeysIndependent(t *testing.T) {
	var td TapDetector
	kc1 := keycode.LT(1, keycode.A)
	kc2 := keycode.TT(2)

	td.Press(0, 0, kc1)
	// 100 ticks 後に 2 つ目を押下
	for i := 0; i < 100; i++ {
		td.Advance()
	}
	td.Press(1, 0, kc2)

	// さらに 100 ticks（kc1 は合計 200、kc2 は 100）
	for i := 0; i < 100; i++ {
		td.Advance()
	}

	// kc1 はタイムアウト
	timedOut1, _ := td.CheckTimeout(0, 0)
	if !timedOut1 {
		t.Error("key (0,0) should have timed out")
	}

	// kc2 はまだ pending
	timedOut2, _ := td.CheckTimeout(1, 0)
	if timedOut2 {
		t.Error("key (1,0) should not have timed out yet")
	}
}

func TestReset(t *testing.T) {
	var td TapDetector
	kc := keycode.LT(1, keycode.A)

	td.Press(0, 0, kc)
	td.Advance()
	td.Reset(0, 0)

	if td.Phase(0, 0) != tapIdle {
		t.Errorf("Phase after Reset = %d, want tapIdle", td.Phase(0, 0))
	}
}

func TestAdvanceOnlyAffectsPending(t *testing.T) {
	var td TapDetector

	// idle 状態のキーは Advance で影響を受けない
	td.Advance()
	timedOut, _ := td.CheckTimeout(0, 0)
	if timedOut {
		t.Error("idle key should not time out")
	}
}

func TestTTKeycode(t *testing.T) {
	var td TapDetector
	kc := keycode.TT(1)

	td.Press(0, 0, kc)
	wasPending, got := td.Release(0, 0)

	if !wasPending {
		t.Error("TT tap should be wasPending=true")
	}
	if got != kc {
		t.Errorf("TT Release keycode = %#x, want %#x", got, kc)
	}
}
