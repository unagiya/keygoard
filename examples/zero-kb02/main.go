package main

import (
	// keyboard パッケージの init() が HID ハンドラを登録するためインポートが必要
	_ "machine/usb/hid/keyboard"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/engine/peripheral/encoder"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

func main() {
	enc := encoder.New(&encoder.Config{
		PinA:   encoderPinA,
		PinB:   encoderPinB,
		KeyCW:  keycode.UpArrow,
		KeyCCW: keycode.DownArrow,
	})

	kb := engine.New(&engine.Config{
		Scanner:     matrix.New(matrixCols, matrixRows),
		Keymap:      defaultKeymap,
		ProductName: "zero-kb02",
		Peripherals: []engine.Peripheral{enc},
	})
	kb.Run()
}
