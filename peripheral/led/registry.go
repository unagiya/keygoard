package led

// maxEffects はレジストリに登録可能なエフェクトの最大数。
const maxEffects = 16

// EffectID はレジストリ内のエフェクトを識別するID。
type EffectID uint8

// Registry はエフェクトの登録・切替を管理する。
type Registry struct {
	effects [maxEffects]Effect
	count   uint8
	current uint8
}

// NewRegistry は新しいRegistryを作成する。
func NewRegistry() *Registry {
	return &Registry{}
}

// Register はエフェクトをレジストリに登録し、IDを返す。
// 最大数を超えた場合は最後のIDを返す。
func (r *Registry) Register(effect Effect) EffectID {
	if r.count >= maxEffects {
		return EffectID(r.count - 1)
	}
	id := EffectID(r.count)
	r.effects[r.count] = effect
	r.count++
	return id
}

// SetCurrent は指定されたIDのエフェクトを現在のエフェクトに設定する。
func (r *Registry) SetCurrent(id EffectID) {
	if uint8(id) < r.count {
		r.current = uint8(id)
	}
}

// Current は現在のエフェクトを返す。
// エフェクトが未登録の場合はnilを返す。
func (r *Registry) Current() Effect {
	if r.count == 0 {
		return nil
	}
	return r.effects[r.current]
}

// CurrentID は現在のエフェクトIDを返す。
func (r *Registry) CurrentID() EffectID {
	return EffectID(r.current)
}

// Next は次のエフェクトに切り替え、そのエフェクトを返す。
// 末尾の場合は先頭に戻る。
func (r *Registry) Next() Effect {
	if r.count == 0 {
		return nil
	}
	r.current = (r.current + 1) % r.count
	return r.effects[r.current]
}

// Previous は前のエフェクトに切り替え、そのエフェクトを返す。
// 先頭の場合は末尾に移動する。
func (r *Registry) Previous() Effect {
	if r.count == 0 {
		return nil
	}
	if r.current == 0 {
		r.current = r.count - 1
	} else {
		r.current--
	}
	return r.effects[r.current]
}

// Count は登録されているエフェクト数を返す。
func (r *Registry) Count() int {
	return int(r.count)
}
