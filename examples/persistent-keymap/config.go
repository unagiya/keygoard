package main

import (
	"machine"

	"github.com/unagiya/keygoard/engine"
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

		// Flash永続化を有効化
		Storage: true,
	}
}
