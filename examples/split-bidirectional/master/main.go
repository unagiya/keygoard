package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/led"
	"github.com/unagiya/keygoard/peripheral/oled"
)

func main() {
	// Master side configuration with LED, OLED, and joystick
	cfg := &engine.BoardConfig{
		// Matrix configuration
		Rows: []machine.Pin{
			machine.GP0, machine.GP1, machine.GP2, machine.GP3,
		},
		Cols: []machine.Pin{
			machine.GP4, machine.GP5, machine.GP6,
		},
		MatrixType: engine.COL2ROW,

		// Timing
		ScanInterval:  1 * time.Millisecond,
		DebounceTime:  3 * time.Millisecond,
		USBReportRate: 1 * time.Millisecond,

		// Joystick (X-axis)
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

		// LED (WS2812)
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: 12,
		},

		// Split keyboard configuration
		Split: &engine.SplitConfig{
			IsMaster:     true,
			UART:         machine.UART1,
			BaudRate:     460800,
			Timeout:      5 * time.Millisecond,
			SendInterval: 1 * time.Millisecond,
			SlaveRows:    4,
			SlaveCols:    3,
		},
	}

	// Keymap
	km := engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0
		{
			// Master side (left)
			{keycode.KC_Q, keycode.KC_W, keycode.KC_E},
			{keycode.KC_A, keycode.KC_S, keycode.KC_D},
			{keycode.KC_Z, keycode.KC_X, keycode.KC_C},
			{keycode.KC_LCTL, keycode.MO(1), keycode.KC_SPC},
			// Slave side (right) - rows 4-7
			{keycode.KC_R, keycode.KC_T, keycode.KC_Y},
			{keycode.KC_F, keycode.KC_G, keycode.KC_H},
			{keycode.KC_V, keycode.KC_B, keycode.KC_N},
			{keycode.KC_ENT, keycode.MO(1), keycode.KC_LSFT},
		},
		// Layer 1
		{
			// Master side
			{keycode.KC_1, keycode.KC_2, keycode.KC_3},
			{keycode.KC_4, keycode.KC_5, keycode.KC_6},
			{keycode.KC_7, keycode.KC_8, keycode.KC_9},
			{keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_0},
			// Slave side
			{keycode.KC_F1, keycode.KC_F2, keycode.KC_F3},
			{keycode.KC_F4, keycode.KC_F5, keycode.KC_F6},
			{keycode.KC_F7, keycode.KC_F8, keycode.KC_F9},
			{keycode.KC_F10, keycode.KC_TRNS, keycode.KC_F11},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// Set LED effect (rainbow on both master and slave)
	kb.SetLEDEffect(led.NewRainbowEffect(5))

	// Run keyboard
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
