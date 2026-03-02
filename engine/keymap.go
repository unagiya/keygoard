package engine

import (
	"github.com/unagiya/keygoard/internal/layer"
	"github.com/unagiya/keygoard/internal/matrix"
	"github.com/unagiya/keygoard/keycode"
)

// MaxLayers はキーマップの最大レイヤー数です。
const MaxLayers = layer.MaxLayers

// RowCount はマトリクスの行数です。
const RowCount = matrix.RowCount

// ColCount はマトリクスの列数です。
const ColCount = matrix.ColCount

// Keymap は複数レイヤーのキー割り当てを保持します。
// Layers[0] がベースレイヤーで、番号が大きいほど優先度が高くなります。
// 未使用レイヤーはゼロ値（全キー None）のままにします。
type Keymap struct {
	// Layers はレイヤーごとのキー割り当てです。[layer][row][col] でインデックスします。
	Layers [layer.MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
}
