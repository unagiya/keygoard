package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/oled"
)

func main() {
	// Configure board with OLED and joystick
	cfg := &engine.BoardConfig{
		// Matrix configuration (3x3 for demo)
		Rows: []machine.Pin{
			machine.GP0, machine.GP1, machine.GP2,
		},
		Cols: []machine.Pin{
			machine.GP3, machine.GP4, machine.GP5,
		},
		MatrixType: engine.COL2ROW,

		// Timing
		ScanInterval:  1 * time.Millisecond,
		DebounceTime:  3 * time.Millisecond,
		USBReportRate: 1 * time.Millisecond,

		// Joystick
		Joystick: &joystick.Config{
			PinX:    machine.ADC0,
			PinY:    machine.ADC1,
			InvertX: false,
			InvertY: false,
		},

		// OLED display
		OLED: &oled.Config{
			I2C:    machine.I2C0,
			Width:  128,
			Height: 32,
		},
	}

	// Keymap demonstrating TT and LT
	km := engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0 (default)
		{
			{keycode.KC_A, keycode.KC_B, keycode.KC_C},
			{keycode.KC_D, keycode.KC_E, keycode.KC_F},
			{keycode.TT(1), keycode.LT(2, keycode.KC_SPC), keycode.KC_ENT},
		},
		// Layer 1 (activated by TT(1) tap-toggle)
		{
			{keycode.KC_1, keycode.KC_2, keycode.KC_3},
			{keycode.KC_4, keycode.KC_5, keycode.KC_6},
			{keycode.KC_TRNS, keycode.KC_7, keycode.KC_8},
		},
		// Layer 2 (activated by LT(2, Space) hold)
		{
			{keycode.KC_F1, keycode.KC_F2, keycode.KC_F3},
			{keycode.KC_F4, keycode.KC_F5, keycode.KC_F6},
			{keycode.KC_F7, keycode.KC_TRNS, keycode.KC_F8},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// Run keyboard
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
