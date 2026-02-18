package engine

import (
	"time"

	"github.com/unagiya/keygoard/keycode"
)

// コンボの制限
const (
	MaxCombos    = 16 // 最大コンボ数
	MaxComboKeys = 4  // コンボあたりの最大キー数
)

// デフォルトのコンボウィンドウ
const DefaultComboWindow = 50 * time.Millisecond

// ComboDef はコンボの定義を表す。
type ComboDef struct {
	Keys   [MaxComboKeys]keycode.Keycode // 構成キー（KC_NOで終端）
	Count  uint8                         // 構成キー数
	Output keycode.Keycode               // 成立時に発行するキーコード
}

// comboTrackerState はコンボ追跡の状態を表す。
type comboTrackerState uint8

const (
	comboStateIdle    comboTrackerState = 0 // 待機中
	comboStateWaiting comboTrackerState = 1 // キー入力待ち中
)

// ComboTracker はアクティブなコンボ検出の状態を追跡する。
type ComboTracker struct {
	state     comboTrackerState
	startTime time.Time
	pressed   [MaxComboKeys]keycode.Keycode // 押下中のキー
	count     uint8                         // 押下中のキー数
	candidate [MaxCombos]bool               // 各コンボの候補フラグ
}

// ComboDetector は複数キー同時押しを検出する。
type ComboDetector struct {
	combos     [MaxCombos]ComboDef
	comboCount uint8
	tracker    ComboTracker
	window     time.Duration
}

// NewComboDetector は新しいComboDetectorを作成する。
func NewComboDetector(window time.Duration) *ComboDetector {
	if window == 0 {
		window = DefaultComboWindow
	}
	return &ComboDetector{
		window: window,
	}
}

// RegisterCombo はコンボを登録する。
// keysは2〜4キー。outputは成立時に発行するキーコード。
func (cd *ComboDetector) RegisterCombo(keys []keycode.Keycode, output keycode.Keycode) error {
	if cd.comboCount >= MaxCombos {
		return ErrInvalidConfig
	}
	if len(keys) < 2 || len(keys) > MaxComboKeys {
		return ErrInvalidConfig
	}

	combo := &cd.combos[cd.comboCount]
	combo.Count = uint8(len(keys))
	combo.Output = output
	for i, k := range keys {
		combo.Keys[i] = k
	}
	// 残りをKC_NOで埋める
	for i := len(keys); i < MaxComboKeys; i++ {
		combo.Keys[i] = keycode.KC_NO
	}
	cd.comboCount++
	return nil
}

// ProcessKeyPress はキー押下を処理する。
// 戻り値: (出力キーコード, consumed)。
// consumed=trueの場合、元のキーは消費される（通常処理しない）。
func (cd *ComboDetector) ProcessKeyPress(kc keycode.Keycode) (keycode.Keycode, bool) {
	if cd.comboCount == 0 {
		return keycode.KC_NO, false
	}

	now := time.Now()

	switch cd.tracker.state {
	case comboStateIdle:
		// このキーがいずれかのコンボの開始キーか確認
		hasCandidate := false
		for i := uint8(0); i < cd.comboCount; i++ {
			if cd.comboContainsKey(&cd.combos[i], kc) {
				cd.tracker.candidate[i] = true
				hasCandidate = true
			} else {
				cd.tracker.candidate[i] = false
			}
		}

		if !hasCandidate {
			return keycode.KC_NO, false
		}

		// コンボ検出開始
		cd.tracker.state = comboStateWaiting
		cd.tracker.startTime = now
		cd.tracker.count = 1
		cd.tracker.pressed[0] = kc
		return keycode.KC_NO, true // キーを消費（保留）

	case comboStateWaiting:
		// タイムアウトチェック
		if now.Sub(cd.tracker.startTime) > cd.window {
			// タイムアウト: 保留キーを返却してリセット
			cd.resetTracker()
			return keycode.KC_NO, false
		}

		// 候補コンボを絞り込む
		for i := uint8(0); i < cd.comboCount; i++ {
			if cd.tracker.candidate[i] && !cd.comboContainsKey(&cd.combos[i], kc) {
				cd.tracker.candidate[i] = false
			}
		}

		// キーを保留リストに追加
		if cd.tracker.count < MaxComboKeys {
			cd.tracker.pressed[cd.tracker.count] = kc
			cd.tracker.count++
		}

		// コンボ成立チェック
		for i := uint8(0); i < cd.comboCount; i++ {
			if cd.tracker.candidate[i] && cd.isComboComplete(&cd.combos[i]) {
				output := cd.combos[i].Output
				cd.resetTracker()
				return output, true
			}
		}

		// まだ候補がある場合はキーを消費して待機
		for i := uint8(0); i < cd.comboCount; i++ {
			if cd.tracker.candidate[i] {
				return keycode.KC_NO, true
			}
		}

		// 候補がなくなった: 保留キーを返却
		cd.resetTracker()
		return keycode.KC_NO, false
	}

	return keycode.KC_NO, false
}

// ProcessKeyRelease はキー離上を処理する。
func (cd *ComboDetector) ProcessKeyRelease(kc keycode.Keycode) {
	// コンボ検出中にキーが離された場合、そのキーを保留リストから除外
	if cd.tracker.state == comboStateWaiting {
		for i := uint8(0); i < cd.tracker.count; i++ {
			if cd.tracker.pressed[i] == kc {
				// キーを除去（シフト）
				for j := i; j < cd.tracker.count-1; j++ {
					cd.tracker.pressed[j] = cd.tracker.pressed[j+1]
				}
				cd.tracker.count--
				cd.tracker.pressed[cd.tracker.count] = keycode.KC_NO
				break
			}
		}
		// 保留キーが無くなったらリセット
		if cd.tracker.count == 0 {
			cd.resetTracker()
		}
	}
}

// Update はタイムアウト処理を行い、タイムアウト時に保留キーを返す。
// 戻り値: 保留キー配列と個数。
func (cd *ComboDetector) Update() (pending [MaxComboKeys]keycode.Keycode, count uint8) {
	if cd.tracker.state != comboStateWaiting {
		return pending, 0
	}

	now := time.Now()
	if now.Sub(cd.tracker.startTime) > cd.window {
		// タイムアウト: 保留キーを返却
		count = cd.tracker.count
		for i := uint8(0); i < count; i++ {
			pending[i] = cd.tracker.pressed[i]
		}
		cd.resetTracker()
	}
	return pending, count
}

// comboContainsKey はコンボにキーが含まれるかを確認する。
func (cd *ComboDetector) comboContainsKey(combo *ComboDef, kc keycode.Keycode) bool {
	for i := uint8(0); i < combo.Count; i++ {
		if combo.Keys[i] == kc {
			return true
		}
	}
	return false
}

// isComboComplete は全コンボキーが押下されているかを確認する。
func (cd *ComboDetector) isComboComplete(combo *ComboDef) bool {
	for i := uint8(0); i < combo.Count; i++ {
		found := false
		for j := uint8(0); j < cd.tracker.count; j++ {
			if cd.tracker.pressed[j] == combo.Keys[i] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// resetTracker はトラッカーをリセットする。
func (cd *ComboDetector) resetTracker() {
	cd.tracker.state = comboStateIdle
	cd.tracker.count = 0
	for i := range cd.tracker.pressed {
		cd.tracker.pressed[i] = keycode.KC_NO
	}
	for i := range cd.tracker.candidate {
		cd.tracker.candidate[i] = false
	}
}
