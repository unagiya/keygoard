package hid

// HIDMode はキーボードHIDのモードを表す。
type HIDMode uint8

const (
	HIDMode6KRO HIDMode = 0 // 6KRO（デフォルト）
	HIDModeNKRO HIDMode = 1 // NKRO
)

// UnifiedKeyboardHID は6KROとNKROを切り替え可能なキーボードHIDラッパー。
type UnifiedKeyboardHID struct {
	sixKRO     *KeyboardHID
	nkroReport NKROReport
	mode       HIDMode
}

// NewUnifiedKeyboardHID は新しいUnifiedKeyboardHIDを作成する。
func NewUnifiedKeyboardHID(mode HIDMode) *UnifiedKeyboardHID {
	return &UnifiedKeyboardHID{
		sixKRO: NewKeyboardHID(),
		mode:   mode,
	}
}

// AddKey はキーをレポートに追加する。現在のモードに応じたレポートに追加。
func (u *UnifiedKeyboardHID) AddKey(kc uint8) bool {
	if u.mode == HIDModeNKRO {
		return u.nkroReport.AddKey(kc)
	}
	return u.sixKRO.report.AddKey(kc)
}

// SetModifier はモディファイアビットを設定する。
func (u *UnifiedKeyboardHID) SetModifier(mask uint8) {
	if u.mode == HIDModeNKRO {
		u.nkroReport.Modifier |= mask
	} else {
		u.sixKRO.report.Modifier |= mask
	}
}

// SendReport は現在のモードに応じたレポートを送信する。
func (u *UnifiedKeyboardHID) SendReport() error {
	if u.mode == HIDModeNKRO {
		// NKROモードでもTinyGoのkeyboard.Deviceを経由して送信
		// NKROレポートをTinyGoの6KROフォーマットに合わせて送信
		// 注: 真のNKROにはカスタムHID記述子が必要だが、
		// TinyGoのUSBスタックの制約により、現時点では6KROデバイスに
		// ビットマップから最大6キーを抽出して送信する。
		// 将来的にTinyGoのカスタムHID記述子対応時に完全NKRO送信に移行。
		return u.sendNKROVia6KRO()
	}
	return u.sixKRO.SendReport()
}

// sendNKROVia6KRO はNKROレポートから最大6キーを抽出して6KROデバイスで送信する。
// NKROの内部状態管理は維持しつつ、USB送信は6KRO互換で行う。
func (u *UnifiedKeyboardHID) sendNKROVia6KRO() error {
	u.sixKRO.report.Clear()
	u.sixKRO.report.Modifier = u.nkroReport.Modifier

	// ビットマップから最大6キーを抽出
	count := 0
	for byteIdx := 0; byteIdx < 15 && count < 6; byteIdx++ {
		if u.nkroReport.Keys[byteIdx] == 0 {
			continue
		}
		for bitIdx := uint8(0); bitIdx < 8 && count < 6; bitIdx++ {
			if u.nkroReport.Keys[byteIdx]&(1<<bitIdx) != 0 {
				keycode := uint8(byteIdx)*8 + bitIdx + 0x04
				u.sixKRO.report.Keys[count] = keycode
				count++
			}
		}
	}

	return u.sixKRO.SendReport()
}

// Clear はレポートをクリアする。
func (u *UnifiedKeyboardHID) Clear() {
	u.sixKRO.Clear()
	u.nkroReport.Clear()
}

// GetReport は6KROレポートへのポインタを返す（後方互換性用）。
func (u *UnifiedKeyboardHID) GetReport() *KeyboardReport {
	return &u.sixKRO.report
}

// GetNKROReport はNKROレポートへのポインタを返す。
func (u *UnifiedKeyboardHID) GetNKROReport() *NKROReport {
	return &u.nkroReport
}

// SetMode はHIDモードを切り替える。
func (u *UnifiedKeyboardHID) SetMode(mode HIDMode) {
	u.mode = mode
}

// GetMode は現在のHIDモードを返す。
func (u *UnifiedKeyboardHID) GetMode() HIDMode {
	return u.mode
}

// ToggleMode はHIDモードをトグルする。
func (u *UnifiedKeyboardHID) ToggleMode() {
	if u.mode == HIDMode6KRO {
		u.mode = HIDModeNKRO
	} else {
		u.mode = HIDMode6KRO
	}
}
