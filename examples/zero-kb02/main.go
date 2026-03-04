package main

import (
	"image/color"

	// keyboard/mouse パッケージの init() が HID ハンドラを登録するためインポートが必要
	_ "machine/usb/hid/keyboard"
	_ "machine/usb/hid/mouse"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

func main() {
	kb := engine.New(&engine.Config{
		ProductName: "zero-kb02",
		ColPins:     []engine.Pin{5, 6, 7, 8},
		RowPins:     []engine.Pin{9, 10, 11},
		Keymap:      defaultKeymap,
		Encoder: &engine.EncoderConfig{
			PinA:   3,
			PinB:   4,
			KeyCW:  keycode.UpArrow,
			KeyCCW: keycode.DownArrow,
		},
		LED: &engine.LEDConfig{
			Pin:   1,
			Count: 12,
			LayerColors: [engine.MaxLayerColors]color.RGBA{
				{R: 0, G: 0, B: 0}, // Layer 0: 黒
				{R: 0, G: 0, B: 8}, // Layer 1: 青
				{R: 8, G: 0, B: 0}, // Layer 2: 赤
				{R: 4, G: 4, B: 4}, // Layer 3: 白
			},
		},
		OLED: &engine.OLEDConfig{
			Bus:      engine.I2C0,
			SDA:      12,
			SCL:      13,
			Address:  0x3C,
			Width:    128,
			Height:   64,
			Rotation: engine.Rotation180,
		},
		Joystick: &engine.JoystickConfig{
			PinX:         29,
			PinY:         28,
			PinButton:    0,
			EnableButton: true,
			ButtonKey:    keycode.None, // マウス左クリック
			Sensitivity:  5,
			InvertY:      true,
		},
	})
	kb.Run()
}
