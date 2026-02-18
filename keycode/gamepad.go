package keycode

// Gamepad buttons (0x6000-0x6FFF)
// Supports up to 32 buttons
const (
	GP_BTN1  Keycode = 0x6000
	GP_BTN2  Keycode = 0x6001
	GP_BTN3  Keycode = 0x6002
	GP_BTN4  Keycode = 0x6003
	GP_BTN5  Keycode = 0x6004
	GP_BTN6  Keycode = 0x6005
	GP_BTN7  Keycode = 0x6006
	GP_BTN8  Keycode = 0x6007
	GP_BTN9  Keycode = 0x6008
	GP_BTN10 Keycode = 0x6009
	GP_BTN11 Keycode = 0x600A
	GP_BTN12 Keycode = 0x600B
	GP_BTN13 Keycode = 0x600C
	GP_BTN14 Keycode = 0x600D
	GP_BTN15 Keycode = 0x600E
	GP_BTN16 Keycode = 0x600F
	GP_BTN17 Keycode = 0x6010
	GP_BTN18 Keycode = 0x6011
	GP_BTN19 Keycode = 0x6012
	GP_BTN20 Keycode = 0x6013
	GP_BTN21 Keycode = 0x6014
	GP_BTN22 Keycode = 0x6015
	GP_BTN23 Keycode = 0x6016
	GP_BTN24 Keycode = 0x6017
	GP_BTN25 Keycode = 0x6018
	GP_BTN26 Keycode = 0x6019
	GP_BTN27 Keycode = 0x601A
	GP_BTN28 Keycode = 0x601B
	GP_BTN29 Keycode = 0x601C
	GP_BTN30 Keycode = 0x601D
	GP_BTN31 Keycode = 0x601E
	GP_BTN32 Keycode = 0x601F
)

// GamepadButtonIndex returns the button index (0-31) for a gamepad button keycode.
// Returns -1 if the keycode is not a gamepad button.
func GamepadButtonIndex(k Keycode) int {
	if !k.IsGamepadButton() {
		return -1
	}
	return int(k - GP_BTN1)
}

// GamepadButtonMask returns a 32-bit mask with the corresponding button bit set.
// Returns 0 if the keycode is not a gamepad button.
func GamepadButtonMask(k Keycode) uint32 {
	idx := GamepadButtonIndex(k)
	if idx < 0 || idx >= 32 {
		return 0
	}
	return 1 << uint(idx)
}
