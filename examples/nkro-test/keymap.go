package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はキーマップを作成する。
// Layer 0にNKRO切り替えキーを配置し、6KRO/NKROの切り替えを試せる。
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: アルファベット + NKRO切り替え
		{
			{keycode.KC_A, keycode.KC_B, keycode.KC_C},
			{keycode.KC_D, keycode.KC_E, keycode.KC_F},
			{keycode.KC_G, keycode.KC_H, keycode.KC_I},
			{keycode.KC_NKRO_TOGGLE, keycode.KC_NKRO_ON, keycode.KC_NKRO_OFF},
		},
	})
}
