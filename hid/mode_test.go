package hid

import "testing"

func TestUnifiedKeyboardHID_DefaultMode(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDMode6KRO)
	if u.GetMode() != HIDMode6KRO {
		t.Error("default mode should be 6KRO")
	}
}

func TestUnifiedKeyboardHID_ModeToggle(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDMode6KRO)

	u.ToggleMode()
	if u.GetMode() != HIDModeNKRO {
		t.Error("mode should be NKRO after toggle")
	}

	u.ToggleMode()
	if u.GetMode() != HIDMode6KRO {
		t.Error("mode should be 6KRO after second toggle")
	}
}

func TestUnifiedKeyboardHID_SetMode(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDMode6KRO)

	u.SetMode(HIDModeNKRO)
	if u.GetMode() != HIDModeNKRO {
		t.Error("mode should be NKRO")
	}

	u.SetMode(HIDMode6KRO)
	if u.GetMode() != HIDMode6KRO {
		t.Error("mode should be 6KRO")
	}
}

func TestUnifiedKeyboardHID_6KRO_AddKey(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDMode6KRO)

	// 6KROモードでは最大6キー
	for i := uint8(0x04); i <= 0x09; i++ {
		if !u.AddKey(i) {
			t.Errorf("AddKey(0x%02X) should succeed in 6KRO mode", i)
		}
	}
	// 7キー目は失敗
	if u.AddKey(0x0A) {
		t.Error("7th key should fail in 6KRO mode")
	}
}

func TestUnifiedKeyboardHID_NKRO_AddKey(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDModeNKRO)

	// NKROモードでは6キー以上OK
	for i := uint8(0x04); i <= 0x0F; i++ {
		if !u.AddKey(i) {
			t.Errorf("AddKey(0x%02X) should succeed in NKRO mode", i)
		}
	}
}

func TestUnifiedKeyboardHID_SetModifier(t *testing.T) {
	// 6KROモード
	u := NewUnifiedKeyboardHID(HIDMode6KRO)
	u.SetModifier(0x01) // LCtrl
	if u.GetReport().Modifier != 0x01 {
		t.Error("6KRO modifier should be set")
	}

	// NKROモード
	u2 := NewUnifiedKeyboardHID(HIDModeNKRO)
	u2.SetModifier(0x01)
	if u2.GetNKROReport().Modifier != 0x01 {
		t.Error("NKRO modifier should be set")
	}
}

func TestUnifiedKeyboardHID_Clear(t *testing.T) {
	u := NewUnifiedKeyboardHID(HIDMode6KRO)
	u.AddKey(0x04)
	u.SetModifier(0x01)

	u.SetMode(HIDModeNKRO)
	u.AddKey(0x05)
	u.SetModifier(0x02)

	u.Clear()
	if u.GetReport().Modifier != 0 {
		t.Error("6KRO report should be cleared")
	}
	if u.GetNKROReport().Modifier != 0 {
		t.Error("NKRO report should be cleared")
	}
}
