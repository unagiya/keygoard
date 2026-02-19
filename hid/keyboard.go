package hid

import (
	"machine"
	"machine/usb/hid"
)

// keyboardReportID はTinyGoのUSB HIDキーボードレポートIDを表す。
// TinyGoのHID記述子ではキーボードはレポートID 1 として登録される。
const keyboardReportID = 0x01

// KeyboardHID はUSBキーボードHID機能をラップする。
type KeyboardHID struct {
	report  KeyboardReport
	buf     *hid.RingBuffer
	waitTxc bool
}

// globalKeyboard はUSB HIDハンドラとして登録されるグローバルインスタンス。
// hid.SetHandler に渡すためシングルトンにする。
var globalKeyboard = &KeyboardHID{
	buf: hid.NewRingBuffer(),
}

func init() {
	// SetHandler がUSBエンドポイントを構成し、ホストがHIDデバイスとして認識するようになる。
	// machine/usb/hid/keyboard パッケージは不要で、直接 hidDevicer を実装して登録する。
	hid.SetHandler(globalKeyboard)
}

// TxHandler はUSB送信完了時に呼ばれるコールバック（hidDevicer インターフェース実装）。
func (k *KeyboardHID) TxHandler() bool {
	k.waitTxc = false
	if b, ok := k.buf.Get(); ok {
		k.waitTxc = true
		hid.SendUSBPacket(b)
		return true
	}
	return false
}

// RxHandler はUSB受信時に呼ばれるコールバック（hidDevicer インターフェース実装）。
// キーボードLED状態の受信に利用可能だが、現在は未使用。
func (k *KeyboardHID) RxHandler(b []byte) bool {
	return false
}

// NewKeyboardHID はキーボードHIDインターフェースを返す。
// USB HIDはシングルトンのため、グローバルインスタンスを返す。
func NewKeyboardHID() *KeyboardHID {
	return globalKeyboard
}

// tx はリングバッファ経由でUSBパケットを送信する。
// USB初期化完了前の呼び出しは無視される。
func (k *KeyboardHID) tx(b []byte) {
	if machine.USBDev.InitEndpointComplete {
		if k.waitTxc {
			k.buf.Put(b)
		} else {
			k.waitTxc = true
			hid.SendUSBPacket(b)
		}
	}
}

// SendReport は現在のキーボードレポートを送信する。
func (k *KeyboardHID) SendReport() error {
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
	k.tx(buf[:])
	return nil
}

// GetReport は現在のキーボードレポートへのポインタを返す。
func (k *KeyboardHID) GetReport() *KeyboardReport {
	return &k.report
}

// Clear はキーボードレポートをクリアする。
func (k *KeyboardHID) Clear() {
	k.report.Clear()
}
