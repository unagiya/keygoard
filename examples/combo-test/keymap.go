package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はキーマップを作成する。
// A+B同時押しでEscが出力されるコンボを試す。
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: アルファベットキー
		{
			{keycode.KC_A, keycode.KC_B, keycode.KC_C},
			{keycode.KC_D, keycode.KC_E, keycode.KC_F},
			{keycode.KC_G, keycode.KC_H, keycode.KC_I},
			{keycode.KC_J, keycode.KC_K, keycode.KC_L},
		},
	})
}
