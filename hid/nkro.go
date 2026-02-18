package hid

// NKROReport はNKRO（Nキーロールオーバー）キーボードHIDレポートを表す。
// ビットマップ方式で最大120キー（0x04-0x77）を同時押し可能。
// レイアウト: Modifier(1byte) + Reserved(1byte) + Keys[15]byte（ビットマップ）
type NKROReport struct {
	Modifier uint8    // モディファイアビットマップ
	Reserved uint8    // 予約（常に0）
	Keys     [15]byte // キービットマップ（120キー: 0x04-0x7B）
}

// Clear はレポートをクリアする。
func (r *NKROReport) Clear() {
	r.Modifier = 0
	r.Reserved = 0
	for i := range r.Keys {
		r.Keys[i] = 0
	}
}

// AddKey はキーをレポートに追加する。
// keycodeはUSB HID Usage ID（0x04-0x7B）。範囲外の場合はfalseを返す。
func (r *NKROReport) AddKey(keycode uint8) bool {
	if keycode < 0x04 || keycode > 0x7B {
		return false
	}
	// 0x04を基準にビットマップに変換
	bit := keycode - 0x04
	byteIdx := bit / 8
	bitIdx := bit % 8
	r.Keys[byteIdx] |= 1 << bitIdx
	return true
}

// RemoveKey はキーをレポートから削除する。
func (r *NKROReport) RemoveKey(keycode uint8) {
	if keycode < 0x04 || keycode > 0x7B {
		return
	}
	bit := keycode - 0x04
	byteIdx := bit / 8
	bitIdx := bit % 8
	r.Keys[byteIdx] &^= 1 << bitIdx
}

// HasKey は指定キーがレポートに含まれているかを返す。
func (r *NKROReport) HasKey(keycode uint8) bool {
	if keycode < 0x04 || keycode > 0x7B {
		return false
	}
	bit := keycode - 0x04
	byteIdx := bit / 8
	bitIdx := bit % 8
	return r.Keys[byteIdx]&(1<<bitIdx) != 0
}

// ToBytes はレポートをバイト列に変換する（17バイト）。
func (r *NKROReport) ToBytes() []byte {
	buf := make([]byte, 17)
	buf[0] = r.Modifier
	buf[1] = r.Reserved
	copy(buf[2:], r.Keys[:])
	return buf
}
