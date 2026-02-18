package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はデフォルトキーマップを作成する。
// Flash保存によりキーマップや設定が永続化される。
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: 数字キー + レイヤー切替
		{
			{keycode.KC_1, keycode.KC_2, keycode.KC_3},
			{keycode.KC_4, keycode.KC_5, keycode.KC_6},
			{keycode.KC_7, keycode.KC_8, keycode.KC_9},
			{keycode.KC_0, keycode.MO(1), keycode.KC_ENT},
		},
		// Layer 1: ファンクションキー
		{
			{keycode.KC_F1, keycode.KC_F2, keycode.KC_F3},
			{keycode.KC_F4, keycode.KC_F5, keycode.KC_F6},
			{keycode.KC_F7, keycode.KC_F8, keycode.KC_F9},
			{keycode.KC_F10, keycode.KC_TRNS, keycode.KC_F12},
		},
	})
}
