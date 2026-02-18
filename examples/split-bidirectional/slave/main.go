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
	// Slave side configuration with LED, OLED, and joystick
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

		// Joystick (Y-axis)
		Joystick: &joystick.Config{
			PinX:    machine.ADC0, // Use as Y-axis in slave
			PinY:    machine.ADC1,
			InvertX: false,
			InvertY: false,
		},

		// OLED display (will receive sync from master)
		OLED: &oled.Config{
			I2C:    machine.I2C0,
			Width:  128,
			Height: 32,
		},

		// LED (WS2812, will receive sync from master)
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: 12,
		},

		// Split keyboard configuration
		Split: &engine.SplitConfig{
			IsMaster:     false,
			UART:         machine.UART1,
			BaudRate:     460800,
			Timeout:      5 * time.Millisecond,
			SendInterval: 1 * time.Millisecond,
			SlaveRows:    4,
			SlaveCols:    3,
		},
	}

	// Keymap (not used on slave, but required for initialization)
	km := engine.NewKeymap([16][][]keycode.Keycode{
		{
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// Run keyboard (slave mode)
	// LED and OLED will be synced from master automatically
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
