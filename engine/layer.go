package engine

import "github.com/unagiya/keygoard/keycode"

// LayerState manages the keyboard layer state.
type LayerState struct {
	active       uint16         // Active layers bitmap (16 layers max)
	moStack      []uint8        // Momentary layer stack
	toggled      uint16         // Toggled layers bitmap
	ttState      map[uint8]bool // Tap-toggle state
	defaultLayer uint8          // Always 0 for Phase 1
}

// NewLayerState creates a new layer state manager.
func NewLayerState() *LayerState {
	return &LayerState{
		active:       0x0001, // Layer 0 is always active by default
		moStack:      make([]uint8, 0, 8),
		toggled:      0,
		ttState:      make(map[uint8]bool),
		defaultLayer: 0,
	}
}

// Current returns the currently active layer (highest priority).
func (ls *LayerState) Current() uint8 {
	// Scan from highest bit (layer 15) to lowest (layer 0)
	for i := 15; i >= 0; i-- {
		if ls.active&(1<<uint(i)) != 0 {
			return uint8(i)
		}
	}
	return 0
}

// ActivateMO activates a momentary layer.
func (ls *LayerState) ActivateMO(layer uint8) {
	if layer >= 16 {
		return
	}
	ls.active |= 1 << uint(layer)
	ls.moStack = append(ls.moStack, layer)
}

// DeactivateMO deactivates a momentary layer.
func (ls *LayerState) DeactivateMO(layer uint8) {
	if layer >= 16 {
		return
	}

	// Remove from stack
	for i, l := range ls.moStack {
		if l == layer {
			ls.moStack = append(ls.moStack[:i], ls.moStack[i+1:]...)
			break
		}
	}

	// Deactivate if not in stack and not toggled
	if !ls.isInStack(layer) && !ls.isToggled(layer) {
		ls.active &^= 1 << uint(layer)
	}
}

// ToggleTG toggles a layer on/off.
func (ls *LayerState) ToggleTG(layer uint8) {
	if layer >= 16 {
		return
	}

	if ls.isToggled(layer) {
		// Turn off
		ls.toggled &^= 1 << uint(layer)
		if !ls.isInStack(layer) {
			ls.active &^= 1 << uint(layer)
		}
	} else {
		// Turn on
		ls.toggled |= 1 << uint(layer)
		ls.active |= 1 << uint(layer)
	}
}

// HandleTT handles tap-toggle layer logic.
// Returns true if layer was toggled.
func (ls *LayerState) HandleTT(layer uint8, isTap bool) bool {
	if layer >= 16 {
		return false
	}

	if isTap {
		// Toggle on tap
		ls.ToggleTG(layer)
		return true
	} else {
		// Momentary on hold
		ls.ActivateMO(layer)
		return false
	}
}

// HandleLT handles layer-tap logic.
// Returns (shouldSendKey, layer) where shouldSendKey is true if it was a tap.
func (ls *LayerState) HandleLT(layer uint8, key keycode.Keycode, isTap bool) (bool, keycode.Keycode) {
	if layer >= 16 {
		return false, key
	}

	if isTap {
		// Send key on tap
		return true, key
	} else {
		// Activate layer on hold
		ls.ActivateMO(layer)
		return false, key
	}
}

// isInStack checks if a layer is in the momentary stack.
func (ls *LayerState) isInStack(layer uint8) bool {
	for _, l := range ls.moStack {
		if l == layer {
			return true
		}
	}
	return false
}

// isToggled checks if a layer is toggled.
func (ls *LayerState) isToggled(layer uint8) bool {
	return ls.toggled&(1<<uint(layer)) != 0
}

// IsActive checks if a layer is currently active.
func (ls *LayerState) IsActive(layer uint8) bool {
	if layer >= 16 {
		return false
	}
	return ls.active&(1<<uint(layer)) != 0
}

// GetActiveLayers returns a slice of all active layer numbers.
func (ls *LayerState) GetActiveLayers() []uint8 {
	layers := make([]uint8, 0, 4)
	for i := uint8(0); i < 16; i++ {
		if ls.IsActive(i) {
			layers = append(layers, i)
		}
	}
	return layers
}
