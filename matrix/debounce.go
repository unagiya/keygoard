// Package matrix はマトリクススキャンとデバウンス処理を実装します。
package matrix

// デバウンスに必要な連続同一読み取り回数
const debounceThr uint8 = 5

// Debouncer はキーごとのデバウンス状態を保持します。
// machine パッケージに依存しないため、標準 Go でテスト可能です。
type Debouncer struct {
	// debounced は現在の確定済みキー状態です。
	debounced [RowCount][ColCount]bool
	// counter は各キーの「前回確定状態と異なる読み取り」の連続カウンタです。
	counter [RowCount][ColCount]uint8
}

// Update は生スキャン結果を受け取り、デバウンス後の状態と変化フラグを返します。
// 同一状態が debounceThr 回連続して読み取れたとき、その状態を確定します。
func (d *Debouncer) Update(raw [RowCount][ColCount]bool) (state [RowCount][ColCount]bool, changed bool) {
	for row := 0; row < RowCount; row++ {
		for col := 0; col < ColCount; col++ {
			if raw[row][col] != d.debounced[row][col] {
				d.counter[row][col]++
				if d.counter[row][col] >= debounceThr {
					d.debounced[row][col] = raw[row][col]
					d.counter[row][col] = 0
					changed = true
				}
			} else {
				d.counter[row][col] = 0
			}
		}
	}
	state = d.debounced
	return
}
