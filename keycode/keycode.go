// Package keycode は HID キーコードの型と定数を定義します。
// machine パッケージに依存しないため、標準 Go でテスト可能です。
package keycode

// Keycode は HID キーコードを表す型です。
// machine/usb/hid/keyboard の Keycode と同じエンコーディングを使うため、
// keyboard.Keycode(kc) でキャストするだけで HID 送信できます。
//
// エンコーディング:
//
//	0x0000        : 無効（キー未割り当て）
//	0xE000-0xE0FF : 修飾キー（左右 Ctrl/Shift/Alt/GUI）
//	0xF000-0xFFFF : 通常キー（HID Usage Page 7）
type Keycode uint16

// None はキー未割り当てを表します。
const None Keycode = 0

// 修飾キー（modifier bitmap | 0xE000）
const (
	ModLeftCtrl   Keycode = 0x01 | 0xE000
	ModLeftShift  Keycode = 0x02 | 0xE000
	ModLeftAlt    Keycode = 0x04 | 0xE000
	ModLeftGUI    Keycode = 0x08 | 0xE000
	ModRightCtrl  Keycode = 0x10 | 0xE000
	ModRightShift Keycode = 0x20 | 0xE000
	ModRightAlt   Keycode = 0x40 | 0xE000
	ModRightGUI   Keycode = 0x80 | 0xE000
)

// 通常キー（HID Usage Page 7 usage code | 0xF000）
const (
	A    Keycode = 4 | 0xF000
	B    Keycode = 5 | 0xF000
	C    Keycode = 6 | 0xF000
	D    Keycode = 7 | 0xF000
	E    Keycode = 8 | 0xF000
	F    Keycode = 9 | 0xF000
	G    Keycode = 10 | 0xF000
	H    Keycode = 11 | 0xF000
	I    Keycode = 12 | 0xF000
	J    Keycode = 13 | 0xF000
	K    Keycode = 14 | 0xF000
	L    Keycode = 15 | 0xF000
	M    Keycode = 16 | 0xF000
	N    Keycode = 17 | 0xF000
	O    Keycode = 18 | 0xF000
	P    Keycode = 19 | 0xF000
	Q    Keycode = 20 | 0xF000
	R    Keycode = 21 | 0xF000
	S    Keycode = 22 | 0xF000
	T    Keycode = 23 | 0xF000
	U    Keycode = 24 | 0xF000
	V    Keycode = 25 | 0xF000
	W    Keycode = 26 | 0xF000
	X    Keycode = 27 | 0xF000
	Y    Keycode = 28 | 0xF000
	Z    Keycode = 29 | 0xF000

	Num1 Keycode = 30 | 0xF000
	Num2 Keycode = 31 | 0xF000
	Num3 Keycode = 32 | 0xF000
	Num4 Keycode = 33 | 0xF000
	Num5 Keycode = 34 | 0xF000
	Num6 Keycode = 35 | 0xF000
	Num7 Keycode = 36 | 0xF000
	Num8 Keycode = 37 | 0xF000
	Num9 Keycode = 38 | 0xF000
	Num0 Keycode = 39 | 0xF000

	Enter     Keycode = 40 | 0xF000
	Escape    Keycode = 41 | 0xF000
	Backspace Keycode = 42 | 0xF000
	Tab       Keycode = 43 | 0xF000
	Space     Keycode = 44 | 0xF000

	Minus     Keycode = 45 | 0xF000
	Equal     Keycode = 46 | 0xF000
	LBracket  Keycode = 47 | 0xF000
	RBracket  Keycode = 48 | 0xF000
	Backslash Keycode = 49 | 0xF000
	Semicolon Keycode = 51 | 0xF000
	Quote     Keycode = 52 | 0xF000
	Grave     Keycode = 53 | 0xF000
	Comma     Keycode = 54 | 0xF000
	Dot       Keycode = 55 | 0xF000
	Slash     Keycode = 56 | 0xF000

	CapsLock Keycode = 57 | 0xF000
	F1       Keycode = 58 | 0xF000
	F2       Keycode = 59 | 0xF000
	F3       Keycode = 60 | 0xF000
	F4       Keycode = 61 | 0xF000
	F5       Keycode = 62 | 0xF000
	F6       Keycode = 63 | 0xF000
	F7       Keycode = 64 | 0xF000
	F8       Keycode = 65 | 0xF000
	F9       Keycode = 66 | 0xF000
	F10      Keycode = 67 | 0xF000
	F11      Keycode = 68 | 0xF000
	F12      Keycode = 69 | 0xF000

	Insert     Keycode = 73 | 0xF000
	Home       Keycode = 74 | 0xF000
	PageUp     Keycode = 75 | 0xF000
	Delete     Keycode = 76 | 0xF000
	End        Keycode = 77 | 0xF000
	PageDown   Keycode = 78 | 0xF000
	RightArrow Keycode = 79 | 0xF000
	LeftArrow  Keycode = 80 | 0xF000
	DownArrow  Keycode = 81 | 0xF000
	UpArrow    Keycode = 82 | 0xF000
)
