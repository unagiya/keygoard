package engine

import "github.com/unagiya/keygoard/keycode"

// MaxLayers はキーマップの最大レイヤー数です。
const MaxLayers = 4

// Resolver はアクティブレイヤーからキーコードを解決します。
// 最上位のアクティブレイヤーから下方向に走査し、最初の非 TRNS キーコードを返します。
type Resolver struct {
	keymap *Keymap
	active [MaxLayers]bool
}

// NewResolver は Resolver を生成します。
// レイヤー 0 は常に有効な状態で初期化されます。
func NewResolver(km *Keymap) *Resolver {
	r := &Resolver{keymap: km}
	r.active[0] = true
	return r
}

// Resolve は指定位置のキーコードを解決します。
// 最上位のアクティブレイヤーから下方向に走査し、最初の非 TRNS キーコードを返します。
// すべてのレイヤーが TRNS または非アクティブの場合は None を返します。
func (r *Resolver) Resolve(row, col int) keycode.Keycode {
	for l := MaxLayers - 1; l >= 0; l-- {
		if !r.active[l] {
			continue
		}
		kc := r.keymap.Layers[l][row][col]
		if kc != keycode.TRNS {
			return kc
		}
	}
	return keycode.None
}

// Activate はレイヤーを有効にします。
// レイヤー 0 は常に有効なため、この関数では変更できません。
func (r *Resolver) Activate(layer int) {
	if layer > 0 && layer < MaxLayers {
		r.active[layer] = true
	}
}

// Deactivate はレイヤーを無効にします。
// レイヤー 0 は常に有効なため、この関数では変更できません。
func (r *Resolver) Deactivate(layer int) {
	if layer > 0 && layer < MaxLayers {
		r.active[layer] = false
	}
}

// Toggle はレイヤーの有効/無効を反転します。
// レイヤー 0 は常に有効なため、この関数では変更できません。
func (r *Resolver) Toggle(layer int) {
	if layer > 0 && layer < MaxLayers {
		r.active[layer] = !r.active[layer]
	}
}

// IsActive はレイヤーが有効かどうかを返します。
func (r *Resolver) IsActive(layer int) bool {
	if layer < 0 || layer >= MaxLayers {
		return false
	}
	return r.active[layer]
}
