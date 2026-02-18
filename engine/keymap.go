package engine

import "github.com/unagiya/keygoard/keycode"

// Keymap represents the keyboard layout with multiple layers.
// Structure: [16 layers][rows][cols]
type Keymap struct {
	layers [16][][]keycode.Keycode
	rows   int
	cols   int
}

// NewKeymap creates a new keymap from a layer definition.
// The input should be a slice of up to 16 layers, where each layer is [rows][cols].
func NewKeymap(layers [16][][]keycode.Keycode) *Keymap {
	// Determine matrix size from first layer
	var rows, cols int
	for i := 0; i < 16; i++ {
		if layers[i] != nil && len(layers[i]) > 0 {
			rows = len(layers[i])
			cols = len(layers[i][0])
			break
		}
	}

	return &Keymap{
		layers: layers,
		rows:   rows,
		cols:   cols,
	}
}

// GetKey returns the keycode at the specified position.
// Returns KC_NO if the position is out of bounds.
func (km *Keymap) GetKey(layer uint8, row, col int) keycode.Keycode {
	if layer >= 16 {
		return keycode.KC_NO
	}

	layerData := km.layers[layer]
	if layerData == nil {
		return keycode.KC_TRNS
	}

	if row < 0 || row >= len(layerData) {
		return keycode.KC_NO
	}

	rowData := layerData[row]
	if col < 0 || col >= len(rowData) {
		return keycode.KC_NO
	}

	return rowData[col]
}

// Validate validates the keymap against the board configuration.
func (km *Keymap) Validate(config *BoardConfig) error {
	// Check that keymap dimensions match matrix size
	configRows, configCols := config.GetMatrixSize()

	// Account for split keyboard
	if config.IsSplitMaster() {
		configRows += config.Split.SlaveRows
	}

	// Validate that at least layer 0 exists and matches dimensions
	if km.layers[0] == nil {
		return ErrInvalidConfig
	}

	if len(km.layers[0]) != configRows {
		return ErrInvalidConfig
	}

	for r := range km.layers[0] {
		if len(km.layers[0][r]) != configCols {
			return ErrInvalidConfig
		}
	}

	return nil
}

// GetMatrixSize returns the keymap matrix dimensions.
func (km *Keymap) GetMatrixSize() (rows, cols int) {
	return km.rows, km.cols
}

// LayerCount は有効なレイヤー数を返す。
func (km *Keymap) LayerCount() int {
	count := 0
	for i := 0; i < 16; i++ {
		if km.layers[i] != nil {
			count = i + 1
		}
	}
	return count
}
