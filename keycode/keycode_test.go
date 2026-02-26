package keycode

import "testing"

func TestKeycodeValues(t *testing.T) {
	// A は HID usage 4 なので 0xF004 のはず
	if A != 0xF004 {
		t.Errorf("A = %#x, want %#x", A, 0xF004)
	}
	// Z は HID usage 29
	if Z != 0xF01D {
		t.Errorf("Z = %#x, want %#x", Z, 0xF01D)
	}
	// ModLeftShift は modifier bitmap 0x02
	if ModLeftShift != 0xE002 {
		t.Errorf("ModLeftShift = %#x, want %#x", ModLeftShift, 0xE002)
	}
	// None は 0
	if None != 0 {
		t.Errorf("None = %#x, want 0", None)
	}
}

func TestKeycodeNoneIsZero(t *testing.T) {
	var kc Keycode
	if kc != None {
		t.Errorf("zero value Keycode should equal None")
	}
}

func TestTRNS(t *testing.T) {
	if TRNS != 0x0001 {
		t.Errorf("TRNS = %#x, want %#x", TRNS, 0x0001)
	}
	if TRNS == None {
		t.Error("TRNS should not equal None")
	}
}

func TestMO(t *testing.T) {
	kc := MO(1)
	if kc != 0xA001 {
		t.Errorf("MO(1) = %#x, want %#x", kc, 0xA001)
	}
	if !kc.IsMO() {
		t.Error("MO(1).IsMO() should be true")
	}
	if kc.IsTG() || kc.IsTT() || kc.IsLT() {
		t.Error("MO(1) should not match TG/TT/LT")
	}
	if !kc.IsLayerAction() {
		t.Error("MO(1).IsLayerAction() should be true")
	}
	if kc.IsTapAction() {
		t.Error("MO(1).IsTapAction() should be false")
	}
	if kc.Layer() != 1 {
		t.Errorf("MO(1).Layer() = %d, want 1", kc.Layer())
	}
}

func TestTG(t *testing.T) {
	kc := TG(2)
	if kc != 0xA012 {
		t.Errorf("TG(2) = %#x, want %#x", kc, 0xA012)
	}
	if !kc.IsTG() {
		t.Error("TG(2).IsTG() should be true")
	}
	if kc.IsMO() || kc.IsTT() || kc.IsLT() {
		t.Error("TG(2) should not match MO/TT/LT")
	}
	if !kc.IsLayerAction() {
		t.Error("TG(2).IsLayerAction() should be true")
	}
	if kc.IsTapAction() {
		t.Error("TG(2).IsTapAction() should be false")
	}
	if kc.Layer() != 2 {
		t.Errorf("TG(2).Layer() = %d, want 2", kc.Layer())
	}
}

func TestTT(t *testing.T) {
	kc := TT(3)
	if kc != 0xA023 {
		t.Errorf("TT(3) = %#x, want %#x", kc, 0xA023)
	}
	if !kc.IsTT() {
		t.Error("TT(3).IsTT() should be true")
	}
	if kc.IsMO() || kc.IsTG() || kc.IsLT() {
		t.Error("TT(3) should not match MO/TG/LT")
	}
	if !kc.IsLayerAction() {
		t.Error("TT(3).IsLayerAction() should be true")
	}
	if !kc.IsTapAction() {
		t.Error("TT(3).IsTapAction() should be true")
	}
	if kc.Layer() != 3 {
		t.Errorf("TT(3).Layer() = %d, want 3", kc.Layer())
	}
}

func TestLT(t *testing.T) {
	// LT(1, A) = 0xB000 | (1 << 8) | 4 = 0xB104
	kc := LT(1, A)
	if kc != 0xB104 {
		t.Errorf("LT(1, A) = %#x, want %#x", kc, 0xB104)
	}
	if !kc.IsLT() {
		t.Error("LT(1, A).IsLT() should be true")
	}
	if kc.IsMO() || kc.IsTG() || kc.IsTT() {
		t.Error("LT(1, A) should not match MO/TG/TT")
	}
	if !kc.IsLayerAction() {
		t.Error("LT(1, A).IsLayerAction() should be true")
	}
	if !kc.IsTapAction() {
		t.Error("LT(1, A).IsTapAction() should be true")
	}
	if kc.Layer() != 1 {
		t.Errorf("LT(1, A).Layer() = %d, want 1", kc.Layer())
	}
	if kc.TapKeycode() != A {
		t.Errorf("LT(1, A).TapKeycode() = %#x, want %#x", kc.TapKeycode(), A)
	}
}

func TestLTLayer0(t *testing.T) {
	kc := LT(0, Space)
	if kc.Layer() != 0 {
		t.Errorf("LT(0, Space).Layer() = %d, want 0", kc.Layer())
	}
	if kc.TapKeycode() != Space {
		t.Errorf("LT(0, Space).TapKeycode() = %#x, want %#x", kc.TapKeycode(), Space)
	}
}

func TestTapKeycodeNonLT(t *testing.T) {
	// LT 以外のキーコードでは None を返す
	if MO(1).TapKeycode() != None {
		t.Error("MO(1).TapKeycode() should return None")
	}
	if A.TapKeycode() != None {
		t.Error("A.TapKeycode() should return None")
	}
}

func TestNormalKeysAreNotLayerAction(t *testing.T) {
	keys := []Keycode{None, TRNS, A, Z, Num0, Enter, ModLeftCtrl, ModRightGUI, F1, F12}
	for _, kc := range keys {
		if kc.IsLayerAction() {
			t.Errorf("%#x.IsLayerAction() should be false", kc)
		}
	}
}
