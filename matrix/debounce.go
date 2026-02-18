package matrix

import "time"

// Debouncer implements eager debouncing.
// Keys are registered immediately on press, but must remain stable for the debounce time.
type Debouncer struct {
	lastState    [][]bool      // last stable state
	lastChangeAt [][]time.Time // timestamp of last state change
	resultBuffer [][]bool      // 事前割り当て済み結果バッファ
	debounceTime time.Duration
	rowCount     int
	colCount     int
}

// NewDebouncer creates a new debouncer.
func NewDebouncer(rowCount, colCount int, debounceTime time.Duration) *Debouncer {
	lastState := make([][]bool, rowCount)
	lastChangeAt := make([][]time.Time, rowCount)
	resultBuffer := make([][]bool, rowCount)

	for i := 0; i < rowCount; i++ {
		lastState[i] = make([]bool, colCount)
		lastChangeAt[i] = make([]time.Time, colCount)
		resultBuffer[i] = make([]bool, colCount)
	}

	return &Debouncer{
		lastState:    lastState,
		lastChangeAt: lastChangeAt,
		resultBuffer: resultBuffer,
		debounceTime: debounceTime,
		rowCount:     rowCount,
		colCount:     colCount,
	}
}

// Process applies debouncing to the raw matrix state.
// Returns the debounced state.
func (d *Debouncer) Process(raw [][]bool) [][]bool {
	now := time.Now()
	result := d.resultBuffer

	for r := 0; r < d.rowCount; r++ {
		for c := 0; c < d.colCount; c++ {
			rawState := raw[r][c]
			lastState := d.lastState[r][c]

			// If state changed, record the time and update immediately (eager)
			if rawState != lastState {
				d.lastChangeAt[r][c] = now
				d.lastState[r][c] = rawState
				result[r][c] = rawState
			} else {
				// State hasn't changed, check if debounce time has passed
				elapsed := now.Sub(d.lastChangeAt[r][c])
				if elapsed >= d.debounceTime {
					// Stable state confirmed
					result[r][c] = lastState
				} else {
					// Still within debounce window, use last stable state
					result[r][c] = lastState
				}
			}
		}
	}

	return result
}
