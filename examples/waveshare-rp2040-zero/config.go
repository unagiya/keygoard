package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/encoder"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/led"
	"github.com/unagiya/keygoard/peripheral/oled"
)

// マトリックス定数
const (
	matrixRows = 3
	matrixCols = 4
	ledCount   = 12 // 3x4 = 12個（各キーに1:1対応）
)

// newBoardConfig はzero-kb02向けのボード設定を作成する。
//
// ピン配置:
//   - マトリックス行: GPIO9, GPIO10, GPIO11
//   - マトリックス列: GPIO5, GPIO6, GPIO7, GPIO8
//   - WS2812 LED: GPIO1（12個）
//   - エンコーダ A/B: GPIO3, GPIO4
//   - エンコーダボタン: GPIO2
//   - ジョイスティック X/Y: GPIO29(ADC3), GPIO28(ADC2)
//   - OLED SDA/SCL: GPIO12, GPIO13（I2C0）
func newBoardConfig() *engine.BoardConfig {
	return &engine.BoardConfig{
		// マトリックス: 3行4列
		Rows: []machine.Pin{
			machine.GPIO9, machine.GPIO10, machine.GPIO11,
		},
		Cols: []machine.Pin{
			machine.GPIO5, machine.GPIO6, machine.GPIO7, machine.GPIO8,
		},
		MatrixType: engine.COL2ROW,

		// タイミング（ゲーミング最適化）
		ScanInterval:  1 * time.Millisecond,
		DebounceTime:  3 * time.Millisecond,
		USBReportRate: 1 * time.Millisecond,

		// WS2812 LED: GPIO1（12個、各キーに1:1対応）
		LED: &led.Config{
			Pin:   machine.GPIO1,
			Count: ledCount,
			KeyToLED: []int8{
				0, 1, 2, 3, // row 0
				4, 5, 6, 7, // row 1
				8, 9, 10, 11, // row 2
			},
			MaxCols: matrixCols,
		},

		// ロータリーエンコーダ: GPIO3(A), GPIO4(B), GPIO2(クリック)
		Encoders: []*encoder.Config{
			{
				PinA:     machine.GPIO3,
				PinB:     machine.GPIO4,
				PinClick: machine.GPIO2,
				CW:       keycode.KC_PGUP,
				CCW:      keycode.KC_PGDN,
				Click:    keycode.KC_ENT,
			},
		},

		// アナログジョイスティック: X=GPIO29(ADC3), Y=GPIO28(ADC2)
		Joystick: &joystick.Config{
			PinX:    machine.ADC3,
			PinY:    machine.ADC2,
			InvertX: false,
			InvertY: true,
		},

		// SSD1306 OLED: I2C0、SDA=GPIO12、SCL=GPIO13、128x32
		OLED: &oled.Config{
			I2C:    machine.I2C0,
			Width:  128,
			Height: 32,
			SDA:    machine.GPIO12,
			SCL:    machine.GPIO13,
		},
	}
}

// ledPositionsX は各LEDのX座標（列方向）。
var ledPositionsX = []uint8{
	0, 3, 6, 9, // row 0
	0, 3, 6, 9, // row 1
	0, 3, 6, 9, // row 2
}

// ledPositionsY は各LEDのY座標（行方向）。
var ledPositionsY = []uint8{
	0, 0, 0, 0, // row 0
	4, 4, 4, 4, // row 1
	8, 8, 8, 8, // row 2
}
