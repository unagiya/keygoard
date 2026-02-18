package main

import "github.com/unagiya/keygoard/keycode"

// Encoder test keyboard layout
var keymapLayers = [16][][]keycode.Keycode{
	// Layer 0: Default
	{
		{keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4},
		{keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
		{keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_F},
		{keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.KC_V},
	},

	// Layers 1-15: Empty
	nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
}
