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
