package tap

import (
	"github.com/unagiya/keygoard/internal/matrix"
	"github.com/unagiya/keygoard/keycode"
)

// Threshold はタップ/ホールド判定の閾値（ticks = ms）です。
// この値未満でリリースされた場合はタップ、以上ホールドされた場合はホールドと判定します。
const Threshold = 200

// Phase はタップ検出の状態を表します。
type Phase uint8

const (
	// Idle は待機状態です。
	Idle Phase = iota
	// Pending はキー押下直後の判定待ち状態です。
	Pending
	// Holding は閾値を超えてホールドと判定された状態です。
	Holding
)

// Detector はタップ/ホールド判定を管理します。
// LT / TT キーコードの押下時にタイマーを開始し、
// リリースまでの時間でタップかホールドかを判定します。
// 全フィールドは固定サイズ配列でヒープ割り当てを回避します。
type Detector struct {
	phase   [matrix.RowCount][matrix.ColCount]Phase
	counter [matrix.RowCount][matrix.ColCount]uint16
	kc      [matrix.RowCount][matrix.ColCount]keycode.Keycode
}

// Press はキー押下時に呼び出します。
// タップ検出が必要なキーコード（LT/TT）の場合、pending 状態に入ります。
func (d *Detector) Press(row, col int, kc keycode.Keycode) {
	d.phase[row][col] = Pending
	d.counter[row][col] = 0
	d.kc[row][col] = kc
}

// Release はキーリリース時に呼び出します。
// pending 状態の場合はタップとして判定し、wasPending=true と元のキーコードを返します。
// holding またはそれ以外の場合は wasPending=false を返します。
func (d *Detector) Release(row, col int) (wasPending bool, kc keycode.Keycode) {
	if d.phase[row][col] == Pending {
		kc = d.kc[row][col]
		d.Reset(row, col)
		return true, kc
	}
	d.Reset(row, col)
	return false, keycode.None
}

// Advance は全 pending キーのカウンタをインクリメントします。
// スキャンサイクルごとに 1 回呼び出します。
func (d *Detector) Advance() {
	for r := 0; r < matrix.RowCount; r++ {
		for c := 0; c < matrix.ColCount; c++ {
			if d.phase[r][c] == Pending {
				d.counter[r][c]++
			}
		}
	}
}

// CheckTimeout は指定位置のタイムアウトを確認します。
// pending 状態で閾値を超えた場合、holding に遷移して timedOut=true を返します。
func (d *Detector) CheckTimeout(row, col int) (timedOut bool, kc keycode.Keycode) {
	if d.phase[row][col] == Pending && d.counter[row][col] >= Threshold {
		d.phase[row][col] = Holding
		return true, d.kc[row][col]
	}
	return false, keycode.None
}

// GetPhase は指定位置の現在のフェーズを返します。
func (d *Detector) GetPhase(row, col int) Phase {
	return d.phase[row][col]
}

// Reset は指定位置の状態を初期化します。
func (d *Detector) Reset(row, col int) {
	d.phase[row][col] = Idle
	d.counter[row][col] = 0
	d.kc[row][col] = keycode.None
}
