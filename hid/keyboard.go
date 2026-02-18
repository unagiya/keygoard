package hid

import (
	"machine/usb/hid"
)

// keyboardReportID はTinyGoのUSB HIDキーボードレポートIDを表す。
const keyboardReportID = 0x02

// KeyboardHID はUSBキーボードHID機能をラップする。
type KeyboardHID struct {
	report KeyboardReport
}

// NewKeyboardHID は新しいキーボードHIDインターフェースを作成する。
func NewKeyboardHID() *KeyboardHID {
	return &KeyboardHID{}
}

// SendReport は現在のキーボードレポートを送信する。
func (k *KeyboardHID) SendReport() error {
	// TinyGoのキーボードHIDレポート形式に合わせたパケットを構築
	// [REPORT_ID, Modifier, Reserved, Key0, Key1, Key2, Key3, Key4, Key5]
	var buf [9]byte
	buf[0] = keyboardReportID
	buf[1] = k.report.Modifier
	buf[2] = k.report.Reserved
	buf[3] = k.report.Keys[0]
	buf[4] = k.report.Keys[1]
	buf[5] = k.report.Keys[2]
	buf[6] = k.report.Keys[3]
	buf[7] = k.report.Keys[4]
	buf[8] = k.report.Keys[5]
	hid.SendUSBPacket(buf[:])
	return nil
}

// GetReport returns a pointer to the current keyboard report.
func (k *KeyboardHID) GetReport() *KeyboardReport {
	return &k.report
}

// Clear clears the keyboard report.
func (k *KeyboardHID) Clear() {
	k.report.Clear()
}
