// Package oled provides OLED display functionality using SSD1306.
package oled

import (
	"image/color"
	"machine"
	"tinygo.org/x/drivers/ssd1306"
)

// Config holds OLED configuration.
type Config struct {
	I2C    *machine.I2C // I2C bus
	Width  int16        // Display width (128)
	Height int16        // Display height (32 or 64)
}

// Display wraps the SSD1306 display driver.
type Display struct {
	device *ssd1306.Device
	width  int16
	height int16
}

// New creates a new OLED display instance.
func New(cfg *Config) (*Display, error) {
	// Configure I2C
	err := cfg.I2C.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
	})
	if err != nil {
		return nil, err
	}

	// Create SSD1306 device
	device := ssd1306.NewI2C(cfg.I2C)

	// Configure the device
	device.Configure(ssd1306.Config{
		Width:  cfg.Width,
		Height: cfg.Height,
	})

	// Clear the display
	device.ClearDisplay()

	return &Display{
		device: device,
		width:  cfg.Width,
		height: cfg.Height,
	}, nil
}

// Clear clears the display.
func (d *Display) Clear() {
	d.device.ClearDisplay()
}

// Display sends the buffer to the display.
func (d *Display) Display() error {
	return d.device.Display()
}

// SetPixel sets a pixel at (x, y).
func (d *Display) SetPixel(x, y int16, on bool) {
	var c color.RGBA
	if on {
		c = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	d.device.SetPixel(x, y, c)
}

// DrawText draws text at the specified position.
// This is a simple implementation for Phase 1.
// For more advanced text rendering, we'll enhance in Phase 2.
func (d *Display) DrawText(x, y int16, text string) {
	// Simple 8x8 character rendering
	// For Phase 1, we'll use a basic implementation
	// In Phase 2, we can add proper font rendering
	for i, ch := range text {
		d.drawChar(x+int16(i*8), y, byte(ch))
	}
}

// drawChar draws a single character using a simple 8x8 font.
// This is a placeholder for Phase 1.
func (d *Display) drawChar(x, y int16, ch byte) {
	// Simplified character drawing
	// In a real implementation, we'd use a font bitmap
	// For now, just draw a simple rectangle as placeholder
	for dy := int16(0); dy < 8; dy++ {
		for dx := int16(0); dx < 6; dx++ {
			// Simple pattern based on character code
			if (ch+byte(dx)+byte(dy))%2 == 0 {
				d.SetPixel(x+dx, y+dy, true)
			}
		}
	}
}

// GetSize returns the display dimensions.
func (d *Display) GetSize() (width, height int16) {
	return d.width, d.height
}

// SetText sets text at a specific position (Phase 2 M3).
// row and col are character positions (8x8 font).
func (d *Display) SetText(row, col int, text string) {
	x := int16(col * 8)
	y := int16(row * 8)
	d.DrawText(x, y, text)
	d.Display()
}

// ShowStatus displays keyboard status information.
// This is a convenience method for showing layer, joystick, and lock states.
func (d *Display) ShowStatus(layer uint8, joyX, joyY int8, capsLock, numLock bool) {
	d.Clear()

	// Display layer (top line)
	d.DrawText(0, 0, "Layer:")
	// Simple number display (0-9)
	if layer < 10 {
		d.DrawText(48, 0, string(rune('0'+layer)))
	}

	// Display joystick values (middle)
	d.DrawText(0, 12, "Joy:")
	// X value
	d.DrawText(32, 12, "X:")
	// Y value
	d.DrawText(0, 20, "Y:")

	// Display lock states (bottom)
	if capsLock {
		d.DrawText(0, 24, "CAPS")
	}
	if numLock {
		d.DrawText(40, 24, "NUM")
	}

	d.Display()
}
