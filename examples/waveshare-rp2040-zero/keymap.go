package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はzero-kb02向けのキーマップを作成する。
//
// 物理レイアウト（3行4列）:
//
//	┌───┬───┬───┬───┐
//	│ Q │ W │ E │ R │ row 0
//	├───┼───┼───┼───┤
//	│ A │ S │ D │ F │ row 1
//	├───┼───┼───┼───┤
//	│ Z │ X │ C │MO1│ row 2
//	└───┴───┴───┴───┘
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// レイヤー0: デフォルト（ゲーミングキー）
		{
			{keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
			{keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_F},
			{keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.MO(1)},
		},
		// レイヤー1: ファンクション
		{
			{keycode.KC_ESC, keycode.KC_F1, keycode.KC_F2, keycode.KC_F3},
			{keycode.KC_TAB, keycode.KC_F4, keycode.KC_F5, keycode.KC_F6},
			{keycode.KC_TRNS, keycode.KC_F7, keycode.KC_F8, keycode.KC_TRNS},
		},
	})
}
