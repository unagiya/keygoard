package engine

import (
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// Keymap は複数レイヤーのキー割り当てを保持します。
// Layers[0] がベースレイヤーで、番号が大きいほど優先度が高くなります。
// 未使用レイヤーはゼロ値（全キー None）のままにします。
type Keymap struct {
	// Layers はレイヤーごとのキー割り当てです。[layer][row][col] でインデックスします。
	Layers [MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
}
