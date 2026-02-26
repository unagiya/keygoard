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
//	0x0001        : 透過（下位レイヤー参照）
//	0xA000-0xA00F : MO(layer) — Momentary レイヤー
//	0xA010-0xA01F : TG(layer) — Toggle レイヤー
//	0xA020-0xA02F : TT(layer) — Tap-Toggle レイヤー
//	0xB000-0xBFFF : LT(layer, kc) — Layer-Tap（上位ニブル=レイヤー, 下位バイト=Usage ID）
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

// --- レイヤーアクション ---

// TRNS は透過キーを表します。
// レイヤー上で TRNS が割り当てられたキーは、下位レイヤーのキーコードを参照します。
const TRNS Keycode = 0x0001

// レイヤーアクションのプレフィックス
const (
	prefixMO Keycode = 0xA000
	prefixTG Keycode = 0xA010
	prefixTT Keycode = 0xA020
	prefixLT Keycode = 0xB000
)

// MO は Momentary レイヤーキーコードを生成します。
// ホールド中のみ指定レイヤーを有効にします。
func MO(layer uint8) Keycode {
	return prefixMO | Keycode(layer)
}

// TG は Toggle レイヤーキーコードを生成します。
// 押下のたびにレイヤーを ON/OFF トグルします。
func TG(layer uint8) Keycode {
	return prefixTG | Keycode(layer)
}

// TT は Tap-Toggle レイヤーキーコードを生成します。
// タップでレイヤーをトグルし、ホールドで MO として動作します。
func TT(layer uint8) Keycode {
	return prefixTT | Keycode(layer)
}

// LT は Layer-Tap キーコードを生成します。
// タップで通常キーコードを送信し、ホールドでレイヤーを有効にします。
// kc には通常キー（0xF0xx）を指定します。下位 8 ビット（HID Usage ID）のみ保持されます。
func LT(layer uint8, kc Keycode) Keycode {
	return prefixLT | Keycode(layer)<<8 | (kc & 0x00FF)
}

// --- レイヤーアクション判定 ---

// IsMO は Momentary レイヤーキーコードかどうかを返します。
func (kc Keycode) IsMO() bool {
	return kc&0xFFF0 == prefixMO
}

// IsTG は Toggle レイヤーキーコードかどうかを返します。
func (kc Keycode) IsTG() bool {
	return kc&0xFFF0 == prefixTG
}

// IsTT は Tap-Toggle レイヤーキーコードかどうかを返します。
func (kc Keycode) IsTT() bool {
	return kc&0xFFF0 == prefixTT
}

// IsLT は Layer-Tap キーコードかどうかを返します。
func (kc Keycode) IsLT() bool {
	return kc&0xF000 == prefixLT
}

// IsLayerAction はレイヤー操作キーコード（MO/TG/TT/LT）かどうかを返します。
func (kc Keycode) IsLayerAction() bool {
	return kc.IsMO() || kc.IsTG() || kc.IsTT() || kc.IsLT()
}

// IsTapAction はタップ検出が必要なキーコード（TT/LT）かどうかを返します。
func (kc Keycode) IsTapAction() bool {
	return kc.IsTT() || kc.IsLT()
}

// --- レイヤーアクション情報抽出 ---

// Layer はレイヤー操作キーコードからレイヤー番号を抽出します。
// レイヤー操作キーコード以外で呼び出した場合の結果は未定義です。
func (kc Keycode) Layer() int {
	if kc.IsLT() {
		return int((kc >> 8) & 0x0F)
	}
	return int(kc & 0x000F)
}

// TapKeycode は LT キーコードからタップ時のキーコードを抽出します。
// LT 以外で呼び出した場合は None を返します。
func (kc Keycode) TapKeycode() Keycode {
	if !kc.IsLT() {
		return None
	}
	usage := kc & 0x00FF
	if usage == 0 {
		return None
	}
	return usage | 0xF000
}
