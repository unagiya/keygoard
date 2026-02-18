package storage

import (
	"testing"

	"github.com/unagiya/keygoard/keycode"
)

func TestSerializeDeserializeKeymap(t *testing.T) {
	var layers [16][][]keycode.Keycode

	// 2レイヤー、3行、4列のキーマップ
	for l := 0; l < 2; l++ {
		layers[l] = make([][]keycode.Keycode, 3)
		for r := 0; r < 3; r++ {
			layers[l][r] = make([]keycode.Keycode, 4)
			for c := 0; c < 4; c++ {
				layers[l][r][c] = keycode.Keycode(l*100 + r*10 + c + 4)
			}
		}
	}

	// シリアライズ
	data := SerializeKeymap(layers, 2, 3, 4)

	// デシリアライズ
	result, numLayers, rows, cols, err := DeserializeKeymap(data)
	if err != nil {
		t.Fatalf("DeserializeKeymap: %v", err)
	}

	if numLayers != 2 || rows != 3 || cols != 4 {
		t.Errorf("dimensions: got %d×%d×%d, want 2×3×4", numLayers, rows, cols)
	}

	// データの一致確認
	for l := 0; l < 2; l++ {
		for r := 0; r < 3; r++ {
			for c := 0; c < 4; c++ {
				expected := keycode.Keycode(l*100 + r*10 + c + 4)
				if result[l][r][c] != expected {
					t.Errorf("layer %d [%d][%d]: got 0x%04X, want 0x%04X",
						l, r, c, uint16(result[l][r][c]), uint16(expected))
				}
			}
		}
	}
}

func TestDeserializeKeymap_InvalidData(t *testing.T) {
	// データが短すぎる
	_, _, _, _, err := DeserializeKeymap([]byte{1, 2})
	if err == nil {
		t.Error("should fail with short data")
	}

	// ゼロ値
	_, _, _, _, err = DeserializeKeymap([]byte{0, 0, 0, 0})
	if err == nil {
		t.Error("should fail with zero dimensions")
	}
}

func TestSerializeKeymap_EmptyLayers(t *testing.T) {
	var layers [16][][]keycode.Keycode
	// レイヤー0のみ、1×1
	layers[0] = [][]keycode.Keycode{{keycode.KC_A}}

	data := SerializeKeymap(layers, 1, 1, 1)
	result, numLayers, rows, cols, err := DeserializeKeymap(data)
	if err != nil {
		t.Fatalf("DeserializeKeymap: %v", err)
	}
	if numLayers != 1 || rows != 1 || cols != 1 {
		t.Errorf("dimensions: got %d×%d×%d, want 1×1×1", numLayers, rows, cols)
	}
	if result[0][0][0] != keycode.KC_A {
		t.Errorf("got 0x%04X, want KC_A", uint16(result[0][0][0]))
	}
}
