// Package hid provides USB HID report structures and utilities.
package hid

// KeyboardReport represents a 6KRO keyboard HID report (8 bytes).
type KeyboardReport struct {
	Modifier uint8    // Modifier keys bitmap
	Reserved uint8    // Reserved (always 0)
	Keys     [6]uint8 // Up to 6 simultaneous key presses
}

// GamepadReport represents a gamepad HID report.
// 32 buttons (4 bytes) + 2 axes (2 bytes) = 6 bytes total.
type GamepadReport struct {
	Buttons uint32 // 32 buttons as bit flags
	X       int8   // X axis (-127 to 127)
	Y       int8   // Y axis (-127 to 127)
}

// Clear clears the keyboard report.
func (r *KeyboardReport) Clear() {
	r.Modifier = 0
	r.Reserved = 0
	for i := range r.Keys {
		r.Keys[i] = 0
	}
}

// AddKey adds a key to the report. Returns true if successful.
// Returns false if the report is full (6KRO limit).
func (r *KeyboardReport) AddKey(keycode uint8) bool {
	// Check if key already exists
	for _, k := range r.Keys {
		if k == keycode {
			return true
		}
	}

	// Find empty slot
	for i := range r.Keys {
		if r.Keys[i] == 0 {
			r.Keys[i] = keycode
			return true
		}
	}

	// Report is full
	return false
}

// RemoveKey removes a key from the report.
func (r *KeyboardReport) RemoveKey(keycode uint8) {
	for i := range r.Keys {
		if r.Keys[i] == keycode {
			r.Keys[i] = 0
			break
		}
	}
}

// ToBytes converts the keyboard report to a byte slice.
func (r *KeyboardReport) ToBytes() []byte {
	return []byte{
		r.Modifier,
		r.Reserved,
		r.Keys[0],
		r.Keys[1],
		r.Keys[2],
		r.Keys[3],
		r.Keys[4],
		r.Keys[5],
	}
}

// Clear clears the gamepad report.
func (g *GamepadReport) Clear() {
	g.Buttons = 0
	g.X = 0
	g.Y = 0
}

// SetButton sets a button state (button index 0-31).
func (g *GamepadReport) SetButton(index int, pressed bool) {
	if index < 0 || index >= 32 {
		return
	}

	if pressed {
		g.Buttons |= 1 << uint(index)
	} else {
		g.Buttons &^= 1 << uint(index)
	}
}

// ToBytes converts the gamepad report to a byte slice.
func (g *GamepadReport) ToBytes() []byte {
	return []byte{
		byte(g.Buttons),
		byte(g.Buttons >> 8),
		byte(g.Buttons >> 16),
		byte(g.Buttons >> 24),
		byte(g.X),
		byte(g.Y),
	}
}
