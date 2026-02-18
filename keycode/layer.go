package keycode

// Layer operations (0x0200-0x02FF)
const (
	layerOpBase Keycode = 0x0200
)

// Layer operation types
const (
	layerOpMO uint8 = 0x00 // Momentary (hold to activate)
	layerOpTG uint8 = 0x01 // Toggle (press to toggle on/off)
	layerOpTT uint8 = 0x02 // Tap Toggle (tap to toggle, hold for momentary)
	layerOpLT uint8 = 0x03 // Layer Tap (tap for key, hold for layer)
)

// MO creates a momentary layer switch keycode.
// While held, the specified layer is activated.
func MO(layer uint8) Keycode {
	return layerOpBase + Keycode(layerOpMO<<8) + Keycode(layer)
}

// TG creates a toggle layer keycode.
// Each press toggles the layer on/off.
func TG(layer uint8) Keycode {
	return layerOpBase + Keycode(uint16(layerOpTG)<<8) + Keycode(layer)
}

// TT creates a tap-toggle layer keycode.
// Tap to toggle, hold for momentary activation.
func TT(layer uint8) Keycode {
	return layerOpBase + Keycode(uint16(layerOpTT)<<8) + Keycode(layer)
}

// LT creates a layer-tap keycode.
// Tap for the key, hold for the layer.
func LT(layer uint8, key Keycode) Keycode {
	// For LT, we encode: 0x02 | (layerOpLT << 8) | (layer << 12)
	// But we need to store the tap key separately
	// For simplicity in Phase 1, we'll use a different encoding
	// LT base: 0x0240, then layer in next 4 bits, key reference in action handler
	return layerOpBase + 0x40 + Keycode(layer)
}

// DecodeLayerOp decodes a layer operation keycode.
// Returns (opType, layer, isLayerOp).
func DecodeLayerOp(k Keycode) (opType uint8, layer uint8, ok bool) {
	if !k.IsLayerOp() {
		return 0, 0, false
	}

	offset := k - layerOpBase

	// Check if it's LT (0x40-0x4F range)
	if offset >= 0x40 && offset < 0x50 {
		return layerOpLT, uint8(offset - 0x40), true
	}

	// Standard layer ops: opType in high byte, layer in low byte
	opType = uint8((offset >> 8) & 0xFF)
	layer = uint8(offset & 0xFF)
	return opType, layer, true
}

// GetLayerOpType returns the layer operation type string.
func GetLayerOpType(opType uint8) string {
	switch opType {
	case layerOpMO:
		return "MO"
	case layerOpTG:
		return "TG"
	case layerOpTT:
		return "TT"
	case layerOpLT:
		return "LT"
	default:
		return "UNKNOWN"
	}
}
