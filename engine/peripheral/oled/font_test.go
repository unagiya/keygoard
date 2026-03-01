package oled

import "testing"

func TestGlyphIndex(t *testing.T) {
	tests := []struct {
		char byte
		want int
	}{
		{'L', 0},
		{'0', 1},
		{'1', 2},
		{'2', 3},
		{'3', 4},
		{'X', -1},
	}
	for _, tt := range tests {
		if got := glyphIndex(tt.char); got != tt.want {
			t.Errorf("glyphIndex(%q) = %d, want %d", tt.char, got, tt.want)
		}
	}
}

func TestFontConstants(t *testing.T) {
	if charWidth != 8 {
		t.Errorf("charWidth = %d, want 8", charWidth)
	}
	if charHeight != 8 {
		t.Errorf("charHeight = %d, want 8", charHeight)
	}
	if fontScale != 3 {
		t.Errorf("fontScale = %d, want 3", fontScale)
	}
}

func TestGlyphs(t *testing.T) {
	if len(glyphs) != 5 {
		t.Fatalf("len(glyphs) = %d, want 5", len(glyphs))
	}

	// 各グリフが charHeight 行であることを確認
	for i, g := range glyphs {
		if len(g) != charHeight {
			t.Errorf("glyphs[%d] has %d rows, want %d", i, len(g), charHeight)
		}
	}

	// 空のグリフがないことを確認
	for i, g := range glyphs {
		allZero := true
		for _, row := range g {
			if row != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			t.Errorf("glyphs[%d] is empty", i)
		}
	}
}
