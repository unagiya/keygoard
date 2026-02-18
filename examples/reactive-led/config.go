package main

import (
	"machine"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/peripheral/led"
)

// キーマトリクスの行数・列数
const (
	matrixRows = 4
	matrixCols = 3
	ledCount   = 12 // 4x3 = 12 LEDs
)

// newBoardConfig はボード設定を作成する。
func newBoardConfig() *engine.BoardConfig {
	return &engine.BoardConfig{
		Rows: []machine.Pin{
			machine.GP0, machine.GP1, machine.GP2, machine.GP3,
		},
		Cols: []machine.Pin{
			machine.GP4, machine.GP5, machine.GP6,
		},
		MatrixType: engine.COL2ROW,

		// LED設定（キーとLEDのマッピング付き）
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: ledCount,
			// 各キー(row*MaxCols+col)に対応するLEDインデックス
			// 4行3列のマトリクスで、各キーに1:1でLEDが対応
			KeyToLED: []int8{
				0, 1, 2, // row 0
				3, 4, 5, // row 1
				6, 7, 8, // row 2
				9, 10, 11, // row 3
			},
			MaxCols: matrixCols,
		},
	}
}

// ledPositionsX は各LEDのX座標（列方向）。
var ledPositionsX = []uint8{
	0, 4, 8, // row 0
	0, 4, 8, // row 1
	0, 4, 8, // row 2
	0, 4, 8, // row 3
}

// ledPositionsY は各LEDのY座標（行方向）。
var ledPositionsY = []uint8{
	0, 0, 0, // row 0
	4, 4, 4, // row 1
	8, 8, 8, // row 2
	12, 12, 12, // row 3
}
