package hid

// CompositeHID manages both keyboard and gamepad HID interfaces.
type CompositeHID struct {
	keyboard *UnifiedKeyboardHID
	gamepad  *GamepadHID
}

// NewCompositeHID creates a new composite HID device.
func NewCompositeHID() *CompositeHID {
	return &CompositeHID{
		keyboard: NewUnifiedKeyboardHID(HIDMode6KRO),
		gamepad:  NewGamepadHID(),
	}
}

// NewCompositeHIDWithMode creates a new composite HID device with the specified HID mode.
func NewCompositeHIDWithMode(mode HIDMode) *CompositeHID {
	return &CompositeHID{
		keyboard: NewUnifiedKeyboardHID(mode),
		gamepad:  NewGamepadHID(),
	}
}

// Keyboard returns the keyboard HID interface.
func (c *CompositeHID) Keyboard() *UnifiedKeyboardHID {
	return c.keyboard
}

// Gamepad returns the gamepad HID interface.
func (c *CompositeHID) Gamepad() *GamepadHID {
	return c.gamepad
}

// SendReports sends both keyboard and gamepad reports.
func (c *CompositeHID) SendReports() error {
	// Send keyboard report
	if err := c.keyboard.SendReport(); err != nil {
		return err
	}

	// Send gamepad report
	// Note: This requires proper USB composite device setup
	if err := c.gamepad.SendReport(); err != nil {
		return err
	}

	return nil
}

// Clear clears both keyboard and gamepad reports.
func (c *CompositeHID) Clear() {
	c.keyboard.Clear()
	c.gamepad.Clear()
}
