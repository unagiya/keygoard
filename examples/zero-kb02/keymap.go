package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// defaultKeymap は zero-kb02 のキーマップです（3 行 × 4 列）。
//
// Layer 0（ベース）:
//
//	Row 0: [Q] [W] [E] [R]
//	Row 1: [A] [S] [D] [F]
//	Row 2: [Z] [X] [C] [MO(1)]
//
// Layer 1（数字・記号）:
//
//	Row 0: [1]    [2]    [3]    [4]
//	Row 1: [5]    [6]    [7]    [8]
//	Row 2: [9]    [0]    [TRNS] [TRNS]
var defaultKeymap = &engine.Keymap{
	Layers: [engine.MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode{
		// Layer 0: ベースレイヤー
		{
			{keycode.Q, keycode.W, keycode.E, keycode.R},
			{keycode.A, keycode.S, keycode.D, keycode.F},
			{keycode.Z, keycode.X, keycode.C, keycode.MO(1)},
		},
		// Layer 1: 数字レイヤー（MO(1) ホールド中に有効）
		{
			{keycode.Num1, keycode.Num2, keycode.Num3, keycode.Num4},
			{keycode.Num5, keycode.Num6, keycode.Num7, keycode.Num8},
			{keycode.Num9, keycode.Num0, keycode.TRNS, keycode.TRNS},
		},
	},
}