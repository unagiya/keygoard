package hid

// GamepadHID handles gamepad HID reports.
// Note: TinyGo doesn't have built-in gamepad support yet,
// so we'll prepare the structure for manual USB descriptor setup.
type GamepadHID struct {
	report GamepadReport
}

// NewGamepadHID creates a new gamepad HID interface.
func NewGamepadHID() *GamepadHID {
	return &GamepadHID{}
}

// GetReport returns a pointer to the current gamepad report.
func (g *GamepadHID) GetReport() *GamepadReport {
	return &g.report
}

// Clear clears the gamepad report.
func (g *GamepadHID) Clear() {
	g.report.Clear()
}

// SendReport sends the current gamepad report.
// TODO: Implement USB gamepad HID descriptor and sending logic.
// For Phase 1, this is a placeholder that will need low-level USB implementation.
func (g *GamepadHID) SendReport() error {
	// Placeholder: In Phase 1, we'll need to implement custom USB HID descriptor
	// for composite device (keyboard + gamepad)
	// This will require working with TinyGo's USB stack directly
	return nil
}
