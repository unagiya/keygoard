package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/encoder"
	"github.com/unagiya/keygoard/peripheral/oled"
)

// Example configuration for testing rotary encoders
// This example assumes:
// - 4 rows x 4 columns matrix
// - 2 rotary encoders with click buttons
// - OLED display on I2C0
var boardConfig = engine.BoardConfig{
	// Matrix configuration
	Rows: []machine.Pin{
		machine.GP0,
		machine.GP1,
		machine.GP2,
		machine.GP3,
	},
	Cols: []machine.Pin{
		machine.GP4,
		machine.GP5,
		machine.GP6,
		machine.GP7,
	},
	MatrixType: engine.COL2ROW,

	// Timing
	ScanInterval:  1 * time.Millisecond,
	DebounceTime:  3 * time.Millisecond,
	USBReportRate: 1 * time.Millisecond,

	// Rotary Encoders (Phase 2)
	Encoders: []*encoder.Config{
		// Encoder 1: Volume control
		{
			PinA:     machine.GP8,  // Encoder A
			PinB:     machine.GP9,  // Encoder B
			PinClick: machine.GP10, // Click button
			CW:       keycode.KC_PPLS, // Volume Up (Keypad +)
			CCW:      keycode.KC_PMNS, // Volume Down (Keypad -)
			Click:    keycode.KC_PENT, // Mute (Keypad Enter)
		},
		// Encoder 2: Navigation
		{
			PinA:     machine.GP11, // Encoder A
			PinB:     machine.GP12, // Encoder B
			PinClick: machine.GP13, // Click button
			CW:       keycode.KC_DOWN, // Down arrow
			CCW:      keycode.KC_UP,   // Up arrow
			Click:    keycode.KC_ENT,  // Enter
		},
	},

	// OLED configuration
	OLED: &oled.Config{
		I2C:    machine.I2C0,
		Width:  128,
		Height: 32,
	},
}
