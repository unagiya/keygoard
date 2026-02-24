package main

import (
	// keyboard パッケージの init() が HID ハンドラを登録するためインポートが必要
	_ "machine/usb/hid/keyboard"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/matrix"
)

func main() {
	kb := engine.New(&engine.Config{
		Scanner:     matrix.New(matrixCols, matrixRows),
		Keymap:      defaultKeymap,
		ProductName: "zero-kb02",
	})
	kb.Run()
}
