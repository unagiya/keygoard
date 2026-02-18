package keycode

// Modifiers (0x0100-0x01FF)
const (
	KC_LCTL Keycode = 0x0100 // Left Control
	KC_LSFT Keycode = 0x0101 // Left Shift
	KC_LALT Keycode = 0x0102 // Left Alt
	KC_LGUI Keycode = 0x0103 // Left GUI (Windows/Command)
	KC_RCTL Keycode = 0x0104 // Right Control
	KC_RSFT Keycode = 0x0105 // Right Shift
	KC_RALT Keycode = 0x0106 // Right Alt
	KC_RGUI Keycode = 0x0107 // Right GUI (Windows/Command)
)

// Modifier bit masks for HID report
const (
	MOD_LCTL uint8 = 0x01
	MOD_LSFT uint8 = 0x02
	MOD_LALT uint8 = 0x04
	MOD_LGUI uint8 = 0x08
	MOD_RCTL uint8 = 0x10
	MOD_RSFT uint8 = 0x20
	MOD_RALT uint8 = 0x40
	MOD_RGUI uint8 = 0x80
)

// ModifierMask converts a modifier keycode to its bit mask.
// Returns 0 if the keycode is not a modifier.
func ModifierMask(k Keycode) uint8 {
	switch k {
	case KC_LCTL:
		return MOD_LCTL
	case KC_LSFT:
		return MOD_LSFT
	case KC_LALT:
		return MOD_LALT
	case KC_LGUI:
		return MOD_LGUI
	case KC_RCTL:
		return MOD_RCTL
	case KC_RSFT:
		return MOD_RSFT
	case KC_RALT:
		return MOD_RALT
	case KC_RGUI:
		return MOD_RGUI
	default:
		return 0
	}
}
