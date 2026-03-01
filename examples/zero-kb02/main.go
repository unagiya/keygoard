package main

import (
	"image/color"

	// keyboard パッケージの init() が HID ハンドラを登録するためインポートが必要
	_ "machine/usb/hid/keyboard"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/engine/peripheral/encoder"
	"github.com/unagiya/keygoard/engine/peripheral/led"
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

	leds := led.New(&led.Config{
		Pin:   ledPin,
		Count: 12,
		LayerColors: [led.MaxLayerColors]color.RGBA{
			{R: 0, G: 8, B: 0, A: 0},       // Layer 0: 緑
			{R: 0, G: 0, B: 8, A: 0},       // Layer 1: 青
			{R: 8, G: 0, B: 0, A: 0},       // Layer 2: 赤
			{R: 4, G: 4, B: 4, A: 0},       // Layer 3: 白
		},
	})

	kb := engine.New(&engine.Config{
		Scanner:     matrix.New(matrixCols, matrixRows),
		Keymap:      defaultKeymap,
		ProductName: "zero-kb02",
		Peripherals: []engine.Peripheral{enc, leds},
	})
	kb.Run()
}
