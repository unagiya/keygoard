// Package joystick provides analog joystick input handling.
package joystick

import (
	"machine"
)

// Config holds joystick configuration.
type Config struct {
	PinX    machine.Pin // ADC pin for X axis
	PinY    machine.Pin // ADC pin for Y axis
	InvertX bool        // Invert X axis
	InvertY bool        // Invert Y axis
}

// Joystick represents an analog joystick.
type Joystick struct {
	adcX     machine.ADC
	adcY     machine.ADC
	centerX  uint16 // Center position for X (calibrated at startup)
	centerY  uint16 // Center position for Y (calibrated at startup)
	invertX  bool
	invertY  bool
	deadzone uint16 // Deadzone threshold (±5% of full range)
}

const (
	// ADC resolution: 12-bit (0-4095)
	adcMax uint16 = 4095
	// Deadzone: 5% of full range
	deadzonePercent = 5
)

// New creates a new joystick instance.
func New(cfg *Config) *Joystick {
	// Configure ADC pins
	adcX := machine.ADC{Pin: cfg.PinX}
	adcX.Configure(machine.ADCConfig{})

	adcY := machine.ADC{Pin: cfg.PinY}
	adcY.Configure(machine.ADCConfig{})

	j := &Joystick{
		adcX:     adcX,
		adcY:     adcY,
		invertX:  cfg.InvertX,
		invertY:  cfg.InvertY,
		deadzone: adcMax * deadzonePercent / 100,
	}

	// Calibrate center position at startup
	j.calibrate()

	return j
}

// calibrate reads the current position and sets it as the center.
func (j *Joystick) calibrate() {
	// Take multiple samples and average
	const samples = 8
	var sumX, sumY uint32

	for i := 0; i < samples; i++ {
		sumX += uint32(j.adcX.Get())
		sumY += uint32(j.adcY.Get())
	}

	j.centerX = uint16(sumX / samples)
	j.centerY = uint16(sumY / samples)
}

// Read reads the joystick axes and returns values in the range -127 to 127.
// Returns (x, y) where 0 is center.
func (j *Joystick) Read() (x, y int8) {
	// Oversample for better precision (average 4 samples)
	const samples = 4
	var sumX, sumY uint32

	for i := 0; i < samples; i++ {
		sumX += uint32(j.adcX.Get())
		sumY += uint32(j.adcY.Get())
	}

	rawX := uint16(sumX / samples)
	rawY := uint16(sumY / samples)

	// Convert to signed values relative to center
	x = j.convertAxis(rawX, j.centerX, j.invertX)
	y = j.convertAxis(rawY, j.centerY, j.invertY)

	return x, y
}

// convertAxis converts a raw ADC value to a signed axis value (-127 to 127).
func (j *Joystick) convertAxis(raw, center uint16, invert bool) int8 {
	// Calculate difference from center
	var diff int32
	if raw > center {
		diff = int32(raw - center)
	} else {
		diff = -int32(center - raw)
	}

	// Apply deadzone
	absDiff := diff
	if absDiff < 0 {
		absDiff = -absDiff
	}
	if uint16(absDiff) < j.deadzone {
		return 0
	}

	// Scale to -127 to 127
	// Maximum deviation from center is ~2047 (half of 4095)
	const maxDeviation = int32(adcMax / 2)
	scaled := (diff * 127) / maxDeviation

	// Clamp to -127 to 127
	if scaled > 127 {
		scaled = 127
	} else if scaled < -127 {
		scaled = -127
	}

	result := int8(scaled)

	// Apply inversion
	if invert {
		result = -result
	}

	return result
}

// GetCenter returns the calibrated center position.
func (j *Joystick) GetCenter() (x, y uint16) {
	return j.centerX, j.centerY
}

// SetCenter sets the center position manually (Phase 2 M4).
func (j *Joystick) SetCenter(x, y uint16) {
	j.centerX = x
	j.centerY = y
}

// Recalibrate recalibrates the center position (Phase 2 M4).
// Should be called when the joystick is at rest.
func (j *Joystick) Recalibrate() {
	j.calibrate()
}

// GetRaw returns the raw ADC values without any processing (Phase 2 M4).
// Useful for calibration UI.
func (j *Joystick) GetRaw() (x, y uint16) {
	return j.adcX.Get(), j.adcY.Get()
}

// SetDeadzone sets a custom deadzone value (Phase 2 M4).
func (j *Joystick) SetDeadzone(dz uint16) {
	j.deadzone = dz
}

// GetDeadzone returns the current deadzone value (Phase 2 M4).
func (j *Joystick) GetDeadzone() uint16 {
	return j.deadzone
}
