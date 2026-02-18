// Package keycode provides keyboard key code definitions and utilities.
package keycode

// Keycode represents a key action or code.
// Uses uint16 to support basic keys, modifiers, layers, and gamepad buttons.
type Keycode uint16

// Keycode ranges:
// 0x0000-0x00FF: Basic keys (USB HID Usage IDs)
// 0x0100-0x01FF: Modifiers
// 0x0200-0x02FF: Layer operations
// 0x6000-0x6FFF: Gamepad buttons
// 0x7000-0x700F: System control (NKRO切り替え等)
// 0x7010-0x701F: Macro (MACRO0〜MACRO15)
// 0x7020-0xFFFF: Custom/Future expansion

// IsBasicKey returns true if the keycode is a basic key (0x00-0xFF).
func (k Keycode) IsBasicKey() bool {
	return k >= 0x00 && k <= 0xFF
}

// IsModifier returns true if the keycode is a modifier (0x0100-0x01FF).
func (k Keycode) IsModifier() bool {
	return k >= 0x0100 && k < 0x0200
}

// IsLayerOp returns true if the keycode is a layer operation (0x0200-0x02FF).
func (k Keycode) IsLayerOp() bool {
	return k >= 0x0200 && k < 0x0300
}

// IsGamepadButton returns true if the keycode is a gamepad button (0x6000-0x6FFF).
func (k Keycode) IsGamepadButton() bool {
	return k >= 0x6000 && k < 0x7000
}

// IsSystemControl returns true if the keycode is a system control key (0x7000-0x700F).
func (k Keycode) IsSystemControl() bool {
	return k >= 0x7000 && k < 0x7010
}

// IsMacro returns true if the keycode is a macro key (0x7010-0x701F).
func (k Keycode) IsMacro() bool {
	return k >= 0x7010 && k < 0x7020
}

// MacroIndex returns the macro index (0-15) for a macro keycode.
// Returns -1 if the keycode is not a macro.
func (k Keycode) MacroIndex() int {
	if !k.IsMacro() {
		return -1
	}
	return int(k - 0x7010)
}

// Special keycodes
const (
	KC_NO   Keycode = 0x0000 // No key
	KC_TRNS Keycode = 0x0001 // Transparent (use lower layer)
)
