package main

import "github.com/unagiya/keygoard/keycode"

// Simple gamepad keyboard layout
// Layer 0: WASD + gamepad buttons
// Layer 1: Function keys and arrows
var keymapLayers = [16][][]keycode.Keycode{
	// Layer 0: Gaming layer
	{
		// Row 0: Number keys and gamepad buttons
		{keycode.KC_ESC, keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4, keycode.KC_5},
		// Row 1: WASD cluster
		{keycode.KC_TAB, keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R, keycode.KC_T},
		// Row 2: Home row with gamepad buttons
		{keycode.KC_LCTL, keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_F, keycode.GP_BTN1},
		// Row 3: Bottom row with layer switch
		{keycode.KC_LSFT, keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.MO(1), keycode.KC_SPC},
	},

	// Layer 1: Function layer
	{
		// Row 0: F-keys
		{keycode.KC_GRV, keycode.KC_F1, keycode.KC_F2, keycode.KC_F3, keycode.KC_F4, keycode.KC_F5},
		// Row 1: Arrow keys and media
		{keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_UP, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS},
		// Row 2: Arrow navigation
		{keycode.KC_TRNS, keycode.KC_LEFT, keycode.KC_DOWN, keycode.KC_RGHT, keycode.KC_TRNS, keycode.GP_BTN2},
		// Row 3: Additional controls
		{keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_ENT},
	},

	// Layers 2-15: Empty
	nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
}
