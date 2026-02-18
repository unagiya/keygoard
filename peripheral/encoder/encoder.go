// Package encoder provides rotary encoder support.
package encoder

import (
	"machine"

	"github.com/unagiya/keygoard/keycode"
)

// Config holds rotary encoder configuration.
type Config struct {
	PinA     machine.Pin     // Encoder A pin
	PinB     machine.Pin     // Encoder B pin
	PinClick machine.Pin     // Click/push button pin (optional, use machine.NoPin if not used)
	CW       keycode.Keycode // Keycode for clockwise rotation
	CCW      keycode.Keycode // Keycode for counter-clockwise rotation
	Click    keycode.Keycode // Keycode for click (optional)
}

// Encoder represents a rotary encoder.
type Encoder struct {
	pinA      machine.Pin
	pinB      machine.Pin
	pinClick  machine.Pin
	lastState uint8
	position  int32
	cw        keycode.Keycode
	ccw       keycode.Keycode
	click     keycode.Keycode
	hasClick  bool
}

// Event represents an encoder event.
type Event uint8

const (
	EventNone  Event = iota
	EventCW          // Clockwise rotation
	EventCCW         // Counter-clockwise rotation
	EventClick       // Button click
)

// New creates a new encoder instance.
func New(cfg *Config) *Encoder {
	e := &Encoder{
		pinA:     cfg.PinA,
		pinB:     cfg.PinB,
		pinClick: cfg.PinClick,
		cw:       cfg.CW,
		ccw:      cfg.CCW,
		click:    cfg.Click,
		hasClick: cfg.PinClick != machine.NoPin,
	}

	// Configure encoder pins as inputs with pull-up
	e.pinA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	e.pinB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure click pin if available
	if e.hasClick {
		e.pinClick.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// Read initial state
	e.lastState = e.readState()

	return e
}

// readState reads the current encoder state (2 bits: A and B).
func (e *Encoder) readState() uint8 {
	a := uint8(0)
	if e.pinA.Get() {
		a = 1
	}
	b := uint8(0)
	if e.pinB.Get() {
		b = 1
	}
	return (a << 1) | b
}

// Update reads the encoder and returns the event.
// Should be called regularly (e.g., every 1ms).
func (e *Encoder) Update() Event {
	currentState := e.readState()

	// Check for rotation using state transition table
	// Gray code: 00 -> 01 -> 11 -> 10 -> 00 (CW)
	//           00 -> 10 -> 11 -> 01 -> 00 (CCW)
	if currentState != e.lastState {
		direction := e.detectDirection(e.lastState, currentState)
		e.lastState = currentState

		if direction > 0 {
			e.position++
			return EventCW
		} else if direction < 0 {
			e.position--
			return EventCCW
		}
	}

	// Check click button
	if e.hasClick && !e.pinClick.Get() {
		return EventClick
	}

	return EventNone
}

// detectDirection detects rotation direction from state transition.
// Returns: 1 for CW, -1 for CCW, 0 for invalid/no change.
func (e *Encoder) detectDirection(lastState, currentState uint8) int {
	// State transition table for quadrature encoder
	// Based on Gray code sequence
	lookup := [16]int{
		0, -1, 1, 0, // 00 -> 00, 01, 10, 11
		1, 0, 0, -1, // 01 -> 00, 01, 10, 11
		-1, 0, 0, 1, // 10 -> 00, 01, 10, 11
		0, 1, -1, 0, // 11 -> 00, 01, 10, 11
	}

	index := (lastState << 2) | currentState
	return lookup[index]
}

// GetPosition returns the current position (accumulated rotation count).
func (e *Encoder) GetPosition() int32 {
	return e.position
}

// GetKeycode returns the keycode for the given event.
func (e *Encoder) GetKeycode(event Event) keycode.Keycode {
	switch event {
	case EventCW:
		return e.cw
	case EventCCW:
		return e.ccw
	case EventClick:
		return e.click
	default:
		return keycode.KC_NO
	}
}

// ResetPosition resets the position counter to zero.
func (e *Encoder) ResetPosition() {
	e.position = 0
}
