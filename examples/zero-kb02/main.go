package main

import (
	"time"

	// keyboard パッケージの init() が HID ハンドラを登録するためインポートが必要
	_ "machine/usb/hid/keyboard"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/matrix"
)

func main() {
	scanner := matrix.New(matrixCols, matrixRows)
	kb := engine.New(scanner, defaultKeymap)
	kb.Init()

	for {
		kb.Tick()
		// fixme: スキャン間隔は実機計測後に最適値に調整する
		time.Sleep(1 * time.Millisecond)
	}
}
