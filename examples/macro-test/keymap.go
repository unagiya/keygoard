package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はキーマップを作成する。
// KC_MACRO0にコピー＆ペーストマクロを割り当てる。
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: 基本キー + マクロキー
		{
			{keycode.KC_A, keycode.KC_B, keycode.KC_C},
			{keycode.KC_D, keycode.KC_E, keycode.KC_F},
			{keycode.KC_G, keycode.KC_H, keycode.KC_I},
			{keycode.KC_MACRO0, keycode.KC_ENT, keycode.KC_BSPC},
		},
	})
}
