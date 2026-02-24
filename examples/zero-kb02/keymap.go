package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// defaultKeymap は zero-kb02 のキーマップです（3 行 × 4 列）。
//
// 物理配列（zero-kb02 スイッチ位置）:
//
//	Row 0: [Q] [W] [E] [R]
//	Row 1: [A] [S] [D] [F]
//	Row 2: [Z] [X] [C] [V]
//
// fixme: Phase 3 以降でレイヤー 1 以降を活用する。
var defaultKeymap = &engine.Keymap{
	Layers: [engine.MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode{
		// Layer 0: ベースレイヤー
		{
			{keycode.Q, keycode.W, keycode.E, keycode.R},
			{keycode.A, keycode.S, keycode.D, keycode.F},
			{keycode.Z, keycode.X, keycode.C, keycode.V},
		},
	},
}
