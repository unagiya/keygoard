package keycode

// Basic keys (USB HID Usage IDs: 0x04-0xFF)
// Reference: USB HID Usage Tables 1.22, Keyboard/Keypad Page (0x07)

// Letters
const (
	KC_A Keycode = 0x04
	KC_B Keycode = 0x05
	KC_C Keycode = 0x06
	KC_D Keycode = 0x07
	KC_E Keycode = 0x08
	KC_F Keycode = 0x09
	KC_G Keycode = 0x0A
	KC_H Keycode = 0x0B
	KC_I Keycode = 0x0C
	KC_J Keycode = 0x0D
	KC_K Keycode = 0x0E
	KC_L Keycode = 0x0F
	KC_M Keycode = 0x10
	KC_N Keycode = 0x11
	KC_O Keycode = 0x12
	KC_P Keycode = 0x13
	KC_Q Keycode = 0x14
	KC_R Keycode = 0x15
	KC_S Keycode = 0x16
	KC_T Keycode = 0x17
	KC_U Keycode = 0x18
	KC_V Keycode = 0x19
	KC_W Keycode = 0x1A
	KC_X Keycode = 0x1B
	KC_Y Keycode = 0x1C
	KC_Z Keycode = 0x1D
)

// Numbers
const (
	KC_1 Keycode = 0x1E
	KC_2 Keycode = 0x1F
	KC_3 Keycode = 0x20
	KC_4 Keycode = 0x21
	KC_5 Keycode = 0x22
	KC_6 Keycode = 0x23
	KC_7 Keycode = 0x24
	KC_8 Keycode = 0x25
	KC_9 Keycode = 0x26
	KC_0 Keycode = 0x27
)

// Special keys
const (
	KC_ENT  Keycode = 0x28 // Enter
	KC_ESC  Keycode = 0x29 // Escape
	KC_BSPC Keycode = 0x2A // Backspace
	KC_TAB  Keycode = 0x2B // Tab
	KC_SPC  Keycode = 0x2C // Space
	KC_MINS Keycode = 0x2D // - and _
	KC_EQL  Keycode = 0x2E // = and +
	KC_LBRC Keycode = 0x2F // [ and {
	KC_RBRC Keycode = 0x30 // ] and }
	KC_BSLS Keycode = 0x31 // \ and |
	KC_SCLN Keycode = 0x33 // ; and :
	KC_QUOT Keycode = 0x34 // ' and "
	KC_GRV  Keycode = 0x35 // ` and ~
	KC_COMM Keycode = 0x36 // , and <
	KC_DOT  Keycode = 0x37 // . and >
	KC_SLSH Keycode = 0x38 // / and ?
	KC_CAPS Keycode = 0x39 // Caps Lock
)

// Function keys
const (
	KC_F1  Keycode = 0x3A
	KC_F2  Keycode = 0x3B
	KC_F3  Keycode = 0x3C
	KC_F4  Keycode = 0x3D
	KC_F5  Keycode = 0x3E
	KC_F6  Keycode = 0x3F
	KC_F7  Keycode = 0x40
	KC_F8  Keycode = 0x41
	KC_F9  Keycode = 0x42
	KC_F10 Keycode = 0x43
	KC_F11 Keycode = 0x44
	KC_F12 Keycode = 0x45
)

// System keys
const (
	KC_PSCR Keycode = 0x46 // Print Screen
	KC_SLCK Keycode = 0x47 // Scroll Lock
	KC_PAUS Keycode = 0x48 // Pause
	KC_INS  Keycode = 0x49 // Insert
	KC_HOME Keycode = 0x4A // Home
	KC_PGUP Keycode = 0x4B // Page Up
	KC_DEL  Keycode = 0x4C // Delete
	KC_END  Keycode = 0x4D // End
	KC_PGDN Keycode = 0x4E // Page Down
)

// Arrow keys
const (
	KC_RGHT Keycode = 0x4F // Right Arrow
	KC_LEFT Keycode = 0x50 // Left Arrow
	KC_DOWN Keycode = 0x51 // Down Arrow
	KC_UP   Keycode = 0x52 // Up Arrow
)

// Keypad
const (
	KC_NLCK Keycode = 0x53 // Num Lock
	KC_PSLS Keycode = 0x54 // Keypad /
	KC_PAST Keycode = 0x55 // Keypad *
	KC_PMNS Keycode = 0x56 // Keypad -
	KC_PPLS Keycode = 0x57 // Keypad +
	KC_PENT Keycode = 0x58 // Keypad Enter
	KC_P1   Keycode = 0x59 // Keypad 1
	KC_P2   Keycode = 0x5A // Keypad 2
	KC_P3   Keycode = 0x5B // Keypad 3
	KC_P4   Keycode = 0x5C // Keypad 4
	KC_P5   Keycode = 0x5D // Keypad 5
	KC_P6   Keycode = 0x5E // Keypad 6
	KC_P7   Keycode = 0x5F // Keypad 7
	KC_P8   Keycode = 0x60 // Keypad 8
	KC_P9   Keycode = 0x61 // Keypad 9
	KC_P0   Keycode = 0x62 // Keypad 0
	KC_PDOT Keycode = 0x63 // Keypad .
)

// Application keys
const (
	KC_APP Keycode = 0x65 // Application (Menu)
)
