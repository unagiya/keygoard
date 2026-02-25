package engine

import (
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// tapThreshold はタップ/ホールド判定の閾値（ticks = ms）です。
// この値未満でリリースされた場合はタップ、以上ホールドされた場合はホールドと判定します。
const tapThreshold = 200

// tapPhase はタップ検出の状態を表します。
type tapPhase uint8

const (
	// tapIdle は待機状態です。
	tapIdle tapPhase = iota
	// tapPending はキー押下直後の判定待ち状態です。
	tapPending
	// tapHolding は閾値を超えてホールドと判定された状態です。
	tapHolding
)

// TapDetector はタップ/ホールド判定を管理します。
// LT / TT キーコードの押下時にタイマーを開始し、
// リリースまでの時間でタップかホールドかを判定します。
// 全フィールドは固定サイズ配列でヒープ割り当てを回避します。
type TapDetector struct {
	phase   [matrix.RowCount][matrix.ColCount]tapPhase
	counter [matrix.RowCount][matrix.ColCount]uint16
	kc      [matrix.RowCount][matrix.ColCount]keycode.Keycode
}

// Press はキー押下時に呼び出します。
// タップ検出が必要なキーコード（LT/TT）の場合、pending 状態に入ります。
func (td *TapDetector) Press(row, col int, kc keycode.Keycode) {
	td.phase[row][col] = tapPending
	td.counter[row][col] = 0
	td.kc[row][col] = kc
}

// Release はキーリリース時に呼び出します。
// pending 状態の場合はタップとして判定し、wasPending=true と元のキーコードを返します。
// holding またはそれ以外の場合は wasPending=false を返します。
func (td *TapDetector) Release(row, col int) (wasPending bool, kc keycode.Keycode) {
	if td.phase[row][col] == tapPending {
		kc = td.kc[row][col]
		td.Reset(row, col)
		return true, kc
	}
	td.Reset(row, col)
	return false, keycode.None
}

// Advance は全 pending キーのカウンタをインクリメントします。
// スキャンサイクルごとに 1 回呼び出します。
func (td *TapDetector) Advance() {
	for r := 0; r < matrix.RowCount; r++ {
		for c := 0; c < matrix.ColCount; c++ {
			if td.phase[r][c] == tapPending {
				td.counter[r][c]++
			}
		}
	}
}

// CheckTimeout は指定位置のタイムアウトを確認します。
// pending 状態で閾値を超えた場合、holding に遷移して timedOut=true を返します。
func (td *TapDetector) CheckTimeout(row, col int) (timedOut bool, kc keycode.Keycode) {
	if td.phase[row][col] == tapPending && td.counter[row][col] >= tapThreshold {
		td.phase[row][col] = tapHolding
		return true, td.kc[row][col]
	}
	return false, keycode.None
}

// Phase は指定位置の現在のフェーズを返します。
func (td *TapDetector) Phase(row, col int) tapPhase {
	return td.phase[row][col]
}

// Reset は指定位置の状態を初期化します。
func (td *TapDetector) Reset(row, col int) {
	td.phase[row][col] = tapIdle
	td.counter[row][col] = 0
	td.kc[row][col] = keycode.None
}
