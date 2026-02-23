package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// defaultKeymap は zero-kb02 レイヤー 0 のキー割り当てです（3 行 × 4 列）。
//
// 物理配列（zero-kb02 スイッチ位置）:
//
//	Row 0: [Q] [W] [E] [R]
//	Row 1: [A] [S] [D] [F]
//	Row 2: [Z] [X] [C] [V]
//
// fixme: 実機確認後にキー割り当てを実用的な配列に変更する。
var defaultKeymap = &engine.Keymap{
	Layer0: [matrix.RowCount][matrix.ColCount]keycode.Keycode{
		{keycode.Q, keycode.W, keycode.E, keycode.R},
		{keycode.A, keycode.S, keycode.D, keycode.F},
		{keycode.Z, keycode.X, keycode.C, keycode.V},
	},
}
