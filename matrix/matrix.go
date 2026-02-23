//go:build tinygo

package matrix

import (
	"machine"
	"time"
)

// Scanner は COL2ROW 方式のマトリクスキースキャナーです。
// 列ピンを 1 本ずつ High にして全行ピンを読み取ります。
type Scanner struct {
	cols     [ColCount]machine.Pin
	rows     [RowCount]machine.Pin
	debounce Debouncer
}

// New は Scanner を生成します。
// cols: 列ピン（出力）, rows: 行ピン（入力プルダウン）
func New(cols [ColCount]machine.Pin, rows [RowCount]machine.Pin) *Scanner {
	return &Scanner{
		cols: cols,
		rows: rows,
	}
}

// Init はピンを初期化します。
// 列ピンを出力 Low、行ピンをプルダウン入力に設定します。
func (s *Scanner) Init() {
	for i := range s.cols {
		s.cols[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		s.cols[i].Low()
	}
	for i := range s.rows {
		s.rows[i].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
}

// Scan は 1 回のマトリクススキャンを実行し、デバウンス後のキー状態と変化フラグを返します。
// 変化がない場合 changed = false を返します。
func (s *Scanner) Scan() (state [RowCount][ColCount]bool, changed bool) {
	var raw [RowCount][ColCount]bool

	for col := 0; col < ColCount; col++ {
		s.cols[col].High()
		// fixme: GPIO 安定化待機時間は実機チャタリング計測後に調整する
		time.Sleep(10 * time.Microsecond)
		for row := 0; row < RowCount; row++ {
			raw[row][col] = s.rows[row].Get()
		}
		s.cols[col].Low()
	}

	return s.debounce.Update(raw)
}
