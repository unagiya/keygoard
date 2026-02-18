package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/oled"
)

// Example configuration for a simple gamepad keyboard
// This example assumes:
// - 4 rows x 6 columns matrix
// - Joystick on GP26 (X) and GP27 (Y)
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
		machine.GP8,
		machine.GP9,
	},
	MatrixType: engine.COL2ROW,

	// Timing - Gaming optimized
	ScanInterval:  1 * time.Millisecond, // 1000Hz scan rate
	DebounceTime:  3 * time.Millisecond, // Fast 3ms debounce
	USBReportRate: 1 * time.Millisecond, // 1000Hz USB report rate

	// Joystick configuration
	Joystick: &joystick.Config{
		PinX:    machine.ADC0, // GP26
		PinY:    machine.ADC1, // GP27
		InvertX: false,
		InvertY: true, // Invert Y axis for natural up/down
	},

	// OLED configuration
	OLED: &oled.Config{
		I2C:    machine.I2C0,
		Width:  128,
		Height: 32,
	},
}
