package storage

import "github.com/unagiya/keygoard/keycode"

// SerializeKeymap はキーマップをバイト列にシリアライズする。
// フォーマット: [layers(1)] [rows(1)] [cols(1)] [reserved(1)] [keycode(2bytes) × layers×rows×cols]
func SerializeKeymap(layers [16][][]keycode.Keycode, numLayers, rows, cols int) []byte {
	headerSize := 4
	dataSize := numLayers * rows * cols * 2
	buf := make([]byte, headerSize+dataSize)

	buf[0] = uint8(numLayers)
	buf[1] = uint8(rows)
	buf[2] = uint8(cols)
	buf[3] = 0 // reserved

	idx := headerSize
	for l := 0; l < numLayers; l++ {
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				var kc keycode.Keycode
				if l < len(layers) && layers[l] != nil && r < len(layers[l]) && c < len(layers[l][r]) {
					kc = layers[l][r][c]
				}
				buf[idx] = byte(kc)
				buf[idx+1] = byte(kc >> 8)
				idx += 2
			}
		}
	}

	return buf
}

// DeserializeKeymap はバイト列からキーマップをデシリアライズする。
// 戻り値: layers, numLayers, rows, cols, error
func DeserializeKeymap(data []byte) ([16][][]keycode.Keycode, int, int, int, error) {
	if len(data) < 4 {
		return [16][][]keycode.Keycode{}, 0, 0, 0, ErrInvalidMagic
	}

	numLayers := int(data[0])
	rows := int(data[1])
	cols := int(data[2])

	if numLayers == 0 || rows == 0 || cols == 0 {
		return [16][][]keycode.Keycode{}, 0, 0, 0, ErrInvalidMagic
	}

	headerSize := 4
	expectedSize := headerSize + numLayers*rows*cols*2
	if len(data) < expectedSize {
		return [16][][]keycode.Keycode{}, 0, 0, 0, ErrDataTooLarge
	}

	var layers [16][][]keycode.Keycode
	idx := headerSize
	for l := 0; l < numLayers && l < 16; l++ {
		layers[l] = make([][]keycode.Keycode, rows)
		for r := 0; r < rows; r++ {
			layers[l][r] = make([]keycode.Keycode, cols)
			for c := 0; c < cols; c++ {
				kc := keycode.Keycode(data[idx]) | keycode.Keycode(data[idx+1])<<8
				layers[l][r][c] = kc
				idx += 2
			}
		}
	}

	return layers, numLayers, rows, cols, nil
}
