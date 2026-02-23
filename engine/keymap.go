package engine

import (
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// Keymap はレイヤー 0 のキー割り当てを保持します。
// Phase 1 では単一レイヤーのみ対応します。
//
// fixme: Phase 2 以降で複数レイヤー（MO/TG）に対応する。
type Keymap struct {
	// Layer0 はレイヤー 0 のキー割り当てです。[row][col] でインデックスします。
	Layer0 [matrix.RowCount][matrix.ColCount]keycode.Keycode
}
