package main

import (
	"machine"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// キーマトリクスの行数・列数
const (
	matrixRows = 4
	matrixCols = 3
)

// マクロ定義: Ctrl+C → 100ms遅延 → Ctrl+V（コピー＆ペースト）
var copyPasteMacro = engine.MacroDef{
	ID: 0,
	Steps: []engine.MacroStep{
		{Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
		{Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_C)},
		{Type: engine.MacroStepKeyUp, Value: uint16(keycode.KC_LCTL)},
		{Type: engine.MacroStepDelay, Value: 100}, // 100ms遅延
		{Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
		{Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_V)},
		{Type: engine.MacroStepKeyUp, Value: uint16(keycode.KC_LCTL)},
		{Type: engine.MacroStepEnd},
	},
}

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

		// マクロ定義を登録
		Macros: []engine.MacroDef{copyPasteMacro},
	}
}
