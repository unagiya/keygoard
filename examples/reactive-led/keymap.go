package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
)

// newKeymap はキーマップを作成する。
func newKeymap() *engine.Keymap {
	return engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: 数字キー + レイヤー切替
		{
			{keycode.KC_1, keycode.KC_2, keycode.KC_3},
			{keycode.KC_4, keycode.KC_5, keycode.KC_6},
			{keycode.KC_7, keycode.KC_8, keycode.KC_9},
			{keycode.KC_0, keycode.MO(1), keycode.MO(2)},
		},
		// Layer 1: アルファベット
		{
			{keycode.KC_Q, keycode.KC_W, keycode.KC_E},
			{keycode.KC_R, keycode.KC_T, keycode.KC_Y},
			{keycode.KC_U, keycode.KC_I, keycode.KC_O},
			{keycode.KC_P, keycode.KC_TRNS, keycode.KC_NO},
		},
		// Layer 2: アルファベット（続き）
		{
			{keycode.KC_A, keycode.KC_S, keycode.KC_D},
			{keycode.KC_F, keycode.KC_G, keycode.KC_H},
			{keycode.KC_J, keycode.KC_K, keycode.KC_L},
			{keycode.KC_ENT, keycode.KC_NO, keycode.KC_TRNS},
		},
	})
}
