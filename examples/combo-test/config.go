package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// キーマトリクスの行数・列数
const (
	matrixRows = 4
	matrixCols = 3
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

		// コンボ定義: A+B同時押し → Esc
		Combos: []engine.ComboDef{
			{
				Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_A, keycode.KC_B},
				Count:  2,
				Output: keycode.KC_ESC,
			},
		},
		ComboWindow: 50 * time.Millisecond,
	}
}
