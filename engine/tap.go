package engine

import (
	"time"

	"github.com/unagiya/keygoard/keycode"
)

// TapConfig holds tap detection configuration.
type TapConfig struct {
	TappingTerm time.Duration // Time window for tap (default: 200ms)
	HoldTerm    time.Duration // Minimum time for hold (default: 150ms)
}

// DefaultTapConfig returns the default tap configuration.
func DefaultTapConfig() *TapConfig {
	return &TapConfig{
		TappingTerm: 200 * time.Millisecond,
		HoldTerm:    150 * time.Millisecond,
	}
}

// TapState represents the state of a tap-capable key.
type TapState int

const (
	TapStateIdle    TapState = iota
	TapStatePressed          // Key is pressed, waiting to determine tap or hold
	TapStateHolding          // Key is being held (layer activation)
	TapStateTapped           // Key was tapped, waiting for release
)

// TapKey tracks the state of a single tap-capable key.
type TapKey struct {
	keycode   keycode.Keycode // The original keycode (TT or LT)
	row       int
	col       int
	state     TapState
	pressTime time.Time
	tapCount  int // For TT: number of taps (for tap-toggle)
}

// TapDetector manages tap detection for TT and LT keycodes.
type TapDetector struct {
	config       *TapConfig
	keys         map[uint32]*TapKey // key: (row << 16) | col
	actionBuffer [4]TapAction       // Update()用の固定バッファ
}

// NewTapDetector creates a new tap detector.
func NewTapDetector(config *TapConfig) *TapDetector {
	if config == nil {
		config = DefaultTapConfig()
	}

	return &TapDetector{
		config: config,
		keys:   make(map[uint32]*TapKey),
	}
}

// makeKey creates a unique key from row and column.
func makeKey(row, col int) uint32 {
	return (uint32(row) << 16) | uint32(col)
}

// ProcessKey processes a key press/release for tap detection.
// Returns the effective keycode to use and whether the key state changed.
func (td *TapDetector) ProcessKey(kc keycode.Keycode, row, col int, pressed bool) (keycode.Keycode, bool) {
	// Only handle TT and LT keycodes
	if !kc.IsLayerOp() {
		return kc, pressed
	}

	opType, layer, ok := keycode.DecodeLayerOp(kc)
	if !ok {
		return kc, pressed
	}

	// Only TT (0x02) and LT (0x03)
	if opType != 0x02 && opType != 0x03 {
		return kc, pressed
	}

	key := makeKey(row, col)
	tk, exists := td.keys[key]

	if pressed {
		// Key pressed
		if !exists {
			// New press
			tk = &TapKey{
				keycode:   kc,
				row:       row,
				col:       col,
				state:     TapStatePressed,
				pressTime: time.Now(),
				tapCount:  0,
			}
			td.keys[key] = tk
		}

		// Don't send any keycode yet, wait to determine tap or hold
		return keycode.KC_NO, false

	} else {
		// Key released
		if !exists || tk.state == TapStateIdle {
			return keycode.KC_NO, false
		}

		now := time.Now()
		duration := now.Sub(tk.pressTime)

		switch tk.state {
		case TapStatePressed:
			// Released before hold term - it's a tap
			if duration < td.config.HoldTerm {
				return td.handleTap(tk, opType, layer)
			} else {
				// Released after hold - was holding
				tk.state = TapStateIdle
				delete(td.keys, key)
				// Release the layer
				return keycode.MO(layer), false
			}

		case TapStateHolding:
			// Was holding, now released
			tk.state = TapStateIdle
			delete(td.keys, key)
			// Release the layer
			return keycode.MO(layer), false

		case TapStateTapped:
			// Already tapped, now releasing
			tk.state = TapStateIdle
			delete(td.keys, key)
			return keycode.KC_NO, false
		}
	}

	return keycode.KC_NO, false
}

// handleTap handles tap action for TT and LT.
func (td *TapDetector) handleTap(tk *TapKey, opType uint8, layer uint8) (keycode.Keycode, bool) {
	switch opType {
	case 0x02: // TT - Tap Toggle
		tk.tapCount++
		tk.state = TapStateTapped

		// Toggle layer on tap
		return keycode.TG(layer), true

	case 0x03: // LT - Layer Tap
		tk.state = TapStateTapped

		// Extract the tap keycode from LT
		// LT format: 0x0300 | (layer << 8) | tapKey
		tapKey := keycode.Keycode(tk.keycode & 0xFF)
		if tapKey == keycode.KC_NO {
			tapKey = keycode.KC_SPC // Default to space if no tap key specified
		}

		// Send the tap key
		return tapKey, true
	}

	return keycode.KC_NO, false
}

// Update checks for keys that need to transition from pressed to holding.
// Should be called regularly (e.g., every scan).
func (td *TapDetector) Update() []TapAction {
	now := time.Now()
	actions := td.actionBuffer[:0]

	for _, tk := range td.keys {
		if tk.state == TapStatePressed {
			duration := now.Sub(tk.pressTime)

			// If held longer than hold term, activate layer
			if duration >= td.config.HoldTerm {
				tk.state = TapStateHolding

				opType, layer, _ := keycode.DecodeLayerOp(tk.keycode)

				// Activate the layer
				actions = append(actions, TapAction{
					Type:    TapActionActivateLayer,
					Layer:   layer,
					Keycode: keycode.MO(layer),
					OpType:  opType,
				})
			}
		}
	}

	return actions
}

// TapActionType represents the type of tap action.
type TapActionType int

const (
	TapActionActivateLayer TapActionType = iota
	TapActionDeactivateLayer
	TapActionSendKey
)

// TapAction represents an action to be taken based on tap detection.
type TapAction struct {
	Type    TapActionType
	Layer   uint8
	Keycode keycode.Keycode
	OpType  uint8
}
