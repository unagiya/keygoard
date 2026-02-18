package keycode

import "testing"

func TestIsSystemControl(t *testing.T) {
	tests := []struct {
		kc   Keycode
		want bool
	}{
		{KC_NKRO_TOGGLE, true},
		{KC_NKRO_ON, true},
		{KC_NKRO_OFF, true},
		{Keycode(0x700F), true},
		{Keycode(0x7010), false}, // マクロ範囲
		{Keycode(0x6FFF), false}, // ゲームパッド範囲
		{KC_NO, false},
	}
	for _, tt := range tests {
		if got := tt.kc.IsSystemControl(); got != tt.want {
			t.Errorf("Keycode(0x%04X).IsSystemControl() = %v, want %v", uint16(tt.kc), got, tt.want)
		}
	}
}

func TestIsMacro(t *testing.T) {
	tests := []struct {
		kc   Keycode
		want bool
	}{
		{Keycode(0x7010), true},
		{Keycode(0x701F), true},
		{Keycode(0x7020), false},
		{Keycode(0x700F), false},
	}
	for _, tt := range tests {
		if got := tt.kc.IsMacro(); got != tt.want {
			t.Errorf("Keycode(0x%04X).IsMacro() = %v, want %v", uint16(tt.kc), got, tt.want)
		}
	}
}

func TestMacroIndex(t *testing.T) {
	tests := []struct {
		kc   Keycode
		want int
	}{
		{Keycode(0x7010), 0},
		{Keycode(0x7015), 5},
		{Keycode(0x701F), 15},
		{Keycode(0x7000), -1}, // システム制御
		{KC_NO, -1},
	}
	for _, tt := range tests {
		if got := tt.kc.MacroIndex(); got != tt.want {
			t.Errorf("Keycode(0x%04X).MacroIndex() = %d, want %d", uint16(tt.kc), got, tt.want)
		}
	}
}
