// Package matrix provides matrix scanning functionality for keyboards.
package matrix

import (
	"machine"
	"time"
)

// MatrixType defines the scanning direction.
type MatrixType uint8

const (
	// COL2ROW: columns are inputs (pulled up), rows are outputs (driven low)
	COL2ROW MatrixType = iota
	// ROW2COL: rows are inputs (pulled up), columns are outputs (driven low)
	ROW2COL
)

// Scanner handles matrix scanning and debouncing.
type Scanner struct {
	rows       []machine.Pin
	cols       []machine.Pin
	matrixType MatrixType
	state      [][]bool // current state [row][col]
	rawBuffer  [][]bool // 事前割り当て済みスキャンバッファ
	debouncer  *Debouncer
	rowCount   int
	colCount   int
}

// Config holds the matrix scanner configuration.
type Config struct {
	Rows         []machine.Pin
	Cols         []machine.Pin
	MatrixType   MatrixType
	DebounceTime time.Duration // default: 3ms
}

// NewScanner creates a new matrix scanner.
func NewScanner(cfg *Config) *Scanner {
	rowCount := len(cfg.Rows)
	colCount := len(cfg.Cols)

	// Initialize state matrix
	state := make([][]bool, rowCount)
	for i := range state {
		state[i] = make([]bool, colCount)
	}

	// 事前割り当て済みスキャンバッファ
	rawBuffer := make([][]bool, rowCount)
	for i := range rawBuffer {
		rawBuffer[i] = make([]bool, colCount)
	}

	// Default debounce time
	debounceTime := cfg.DebounceTime
	if debounceTime == 0 {
		debounceTime = 3 * time.Millisecond
	}

	scanner := &Scanner{
		rows:       cfg.Rows,
		cols:       cfg.Cols,
		matrixType: cfg.MatrixType,
		state:      state,
		rawBuffer:  rawBuffer,
		debouncer:  NewDebouncer(rowCount, colCount, debounceTime),
		rowCount:   rowCount,
		colCount:   colCount,
	}

	scanner.init()
	return scanner
}

// init initializes GPIO pins for matrix scanning.
func (s *Scanner) init() {
	switch s.matrixType {
	case COL2ROW:
		// Rows: outputs (driven low when scanning)
		for _, pin := range s.rows {
			pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
			pin.High() // default high (not scanning)
		}
		// Cols: inputs with pull-up
		for _, pin := range s.cols {
			pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}

	case ROW2COL:
		// Cols: outputs (driven low when scanning)
		for _, pin := range s.cols {
			pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
			pin.High() // default high (not scanning)
		}
		// Rows: inputs with pull-up
		for _, pin := range s.rows {
			pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
	}
}

// Scan performs one matrix scan and returns the debounced state.
// Returns a 2D slice [row][col] of key states (true = pressed).
func (s *Scanner) Scan() [][]bool {
	// Read raw matrix state
	raw := s.scanRaw()

	// Apply debouncing
	debounced := s.debouncer.Process(raw)

	// Update internal state
	for r := 0; r < s.rowCount; r++ {
		for c := 0; c < s.colCount; c++ {
			s.state[r][c] = debounced[r][c]
		}
	}

	return s.state
}

// scanRaw performs a raw matrix scan without debouncing.
func (s *Scanner) scanRaw() [][]bool {
	// ゼロクリアして再利用
	for r := 0; r < s.rowCount; r++ {
		for c := 0; c < s.colCount; c++ {
			s.rawBuffer[r][c] = false
		}
	}
	raw := s.rawBuffer

	switch s.matrixType {
	case COL2ROW:
		// Drive each row low and read columns
		for r, rowPin := range s.rows {
			rowPin.Low()                     // drive row low
			time.Sleep(1 * time.Microsecond) // short delay for signal stability

			for c, colPin := range s.cols {
				// Column is pulled up, pressed key pulls it low
				raw[r][c] = !colPin.Get()
			}

			rowPin.High() // restore row to high
		}

	case ROW2COL:
		// Drive each column low and read rows
		for c, colPin := range s.cols {
			colPin.Low()                     // drive column low
			time.Sleep(1 * time.Microsecond) // short delay for signal stability

			for r, rowPin := range s.rows {
				// Row is pulled up, pressed key pulls it low
				raw[r][c] = !rowPin.Get()
			}

			colPin.High() // restore column to high
		}
	}

	return raw
}

// GetState returns the current debounced state of the matrix.
func (s *Scanner) GetState() [][]bool {
	return s.state
}

// GetKeyState returns the state of a specific key at (row, col).
func (s *Scanner) GetKeyState(row, col int) bool {
	if row < 0 || row >= s.rowCount || col < 0 || col >= s.colCount {
		return false
	}
	return s.state[row][col]
}
