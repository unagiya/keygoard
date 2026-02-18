package engine

import (
	"testing"

	"github.com/unagiya/keygoard/keycode"
)

func TestMacroManager_RegisterMacro(t *testing.T) {
	m := NewMacroManager()

	steps := []MacroStep{
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_A)},
		{Type: MacroStepDelay, Value: 50},
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_B)},
	}

	if err := m.RegisterMacro(0, steps); err != nil {
		t.Errorf("RegisterMacro should succeed: %v", err)
	}

	// 範囲外のID
	if err := m.RegisterMacro(16, steps); err == nil {
		t.Error("RegisterMacro with id=16 should fail")
	}

	// ステップ数上限超過
	longSteps := make([]MacroStep, MaxMacroSteps+1)
	if err := m.RegisterMacro(1, longSteps); err == nil {
		t.Error("RegisterMacro with too many steps should fail")
	}
}

func TestMacroManager_Trigger(t *testing.T) {
	m := NewMacroManager()

	steps := []MacroStep{
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_A)},
	}
	m.RegisterMacro(0, steps)

	if m.IsRunning() {
		t.Error("should not be running before trigger")
	}

	m.Trigger(0)
	if !m.IsRunning() {
		t.Error("should be running after trigger")
	}
}

func TestMacroManager_TriggerEmpty(t *testing.T) {
	m := NewMacroManager()

	// 空マクロはトリガーしない
	m.Trigger(0)
	if m.IsRunning() {
		t.Error("should not be running for empty macro")
	}
}

func TestMacroManager_Update_Sequence(t *testing.T) {
	m := NewMacroManager()

	steps := []MacroStep{
		{Type: MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_C)},
		{Type: MacroStepKeyUp, Value: uint16(keycode.KC_LCTL)},
	}
	m.RegisterMacro(0, steps)
	m.Trigger(0)

	// ステップ1: KeyDown LCtrl
	kc, st, done := m.Update()
	if kc != keycode.KC_LCTL || st != MacroStepKeyDown || done {
		t.Errorf("step 1: got kc=0x%04X st=%d done=%v", uint16(kc), st, done)
	}

	// ステップ2: KeyTap C
	kc, st, done = m.Update()
	if kc != keycode.KC_C || st != MacroStepKeyTap || done {
		t.Errorf("step 2: got kc=0x%04X st=%d done=%v", uint16(kc), st, done)
	}

	// ステップ3: KeyUp LCtrl
	kc, st, done = m.Update()
	if kc != keycode.KC_LCTL || st != MacroStepKeyUp || done {
		t.Errorf("step 3: got kc=0x%04X st=%d done=%v", uint16(kc), st, done)
	}

	// 完了
	kc, st, done = m.Update()
	if !done {
		t.Errorf("should be done: got kc=0x%04X st=%d done=%v", uint16(kc), st, done)
	}

	if m.IsRunning() {
		t.Error("should not be running after completion")
	}
}

func TestMacroManager_Update_Delay(t *testing.T) {
	m := NewMacroManager()

	steps := []MacroStep{
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_A)},
		{Type: MacroStepDelay, Value: 3},
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_B)},
	}
	m.RegisterMacro(0, steps)
	m.Trigger(0)

	// ステップ1: KeyTap A
	kc, st, _ := m.Update()
	if kc != keycode.KC_A || st != MacroStepKeyTap {
		t.Errorf("step 1: got kc=0x%04X st=%d", uint16(kc), st)
	}

	// ステップ2: Delay開始
	_, st, _ = m.Update()
	if st != MacroStepDelay {
		t.Errorf("step 2 should be delay, got st=%d", st)
	}

	// 遅延中（3ティック）
	for i := 0; i < 3; i++ {
		_, st, _ = m.Update()
		if st != MacroStepDelay {
			t.Errorf("delay tick %d: should still be delay, got st=%d", i, st)
		}
	}

	// 遅延終了後: KeyTap B
	kc, st, _ = m.Update()
	if kc != keycode.KC_B || st != MacroStepKeyTap {
		t.Errorf("after delay: got kc=0x%04X st=%d", uint16(kc), st)
	}
}

func TestMacroManager_TriggerWhileRunning(t *testing.T) {
	m := NewMacroManager()

	steps := []MacroStep{
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_A)},
		{Type: MacroStepDelay, Value: 100},
		{Type: MacroStepKeyTap, Value: uint16(keycode.KC_B)},
	}
	m.RegisterMacro(0, steps)
	m.RegisterMacro(1, steps)

	m.Trigger(0)

	// 実行中に別のマクロをトリガー→無視される
	m.Update() // KeyTap A
	m.Trigger(1)

	// まだマクロ0が実行中
	if !m.IsRunning() {
		t.Error("should still be running macro 0")
	}
}
