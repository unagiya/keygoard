// Package engine provides the core keyboard engine that integrates all components.
package engine

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/hid"
	"github.com/unagiya/keygoard/matrix"
	"github.com/unagiya/keygoard/peripheral/encoder"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/led"
	"github.com/unagiya/keygoard/peripheral/oled"
)

// MatrixType re-exports matrix.MatrixType for convenience.
type MatrixType = matrix.MatrixType

const (
	COL2ROW = matrix.COL2ROW
	ROW2COL = matrix.ROW2COL
)

// BoardConfig holds the complete keyboard configuration.
type BoardConfig struct {
	// Matrix configuration
	Rows       []machine.Pin
	Cols       []machine.Pin
	MatrixType MatrixType

	// Timing
	ScanInterval  time.Duration // Default: 1ms
	DebounceTime  time.Duration // Default: 3ms
	USBReportRate time.Duration // Default: 1ms (1000Hz)

	// Peripherals (optional)
	Joystick *joystick.Config
	OLED     *oled.Config
	Encoders []*encoder.Config // Up to 2 encoders (Phase 2)
	LED      *led.Config       // RGB LED (Phase 2)

	// Split keyboard (optional)
	Split *SplitConfig

	// HID mode (Phase 4)
	DefaultHIDMode hid.HIDMode // デフォルトHIDモード（0=6KRO, 1=NKRO）

	// Macro definitions (Phase 4)
	Macros []MacroDef

	// Combo definitions (Phase 4)
	Combos      []ComboDef
	ComboWindow time.Duration // デフォルト: 50ms

	// Storage (Phase 4)
	Storage bool // Flash永続化を有効にする

	// Debug enables debug logging via println() (USB CDC serial output)
	Debug bool
}

// LEDSyncMode はLED同期モードを表す（Phase 3 M2）。
type LEDSyncMode uint8

const (
	LEDSyncUnified     LEDSyncMode = 0 // マスターが両側のLEDを制御（デフォルト）
	LEDSyncIndependent LEDSyncMode = 1 // 各側が独立してエフェクト実行
)

// SplitConfig holds split keyboard configuration.
type SplitConfig struct {
	// Mode
	IsMaster bool // true for master, false for slave

	// UART
	UART     *machine.UART
	BaudRate uint32 // Default: 460800

	// Master side configuration
	SlaveRows     int           // Number of rows on slave side
	SlaveCols     int           // Number of columns on slave side
	SlaveLEDCount int           // スレーブ側のLED数（Phase 3 M2）
	Timeout       time.Duration // Default: 5ms

	// Slave side configuration
	SendInterval time.Duration // Default: 1ms

	// LED sync mode（Phase 3 M2）
	LEDSyncMode LEDSyncMode
}

// Validate validates the board configuration.
func (c *BoardConfig) Validate() error {
	if len(c.Rows) == 0 || len(c.Cols) == 0 {
		return ErrInvalidConfig
	}

	// Set defaults
	if c.ScanInterval == 0 {
		c.ScanInterval = 1 * time.Millisecond
	}
	if c.DebounceTime == 0 {
		c.DebounceTime = 3 * time.Millisecond
	}
	if c.USBReportRate == 0 {
		c.USBReportRate = 1 * time.Millisecond
	}

	// Validate split configuration
	if c.Split != nil {
		if c.Split.UART == nil {
			return ErrInvalidConfig
		}
		if c.Split.BaudRate == 0 {
			c.Split.BaudRate = 460800
		}
		if c.Split.IsMaster {
			if c.Split.SlaveRows == 0 || c.Split.SlaveCols == 0 {
				return ErrInvalidConfig
			}
			if c.Split.Timeout == 0 {
				c.Split.Timeout = 5 * time.Millisecond
			}
		} else {
			if c.Split.SendInterval == 0 {
				c.Split.SendInterval = 1 * time.Millisecond
			}
		}
	}

	return nil
}

// GetMatrixSize returns the matrix dimensions.
func (c *BoardConfig) GetMatrixSize() (rows, cols int) {
	return len(c.Rows), len(c.Cols)
}

// HasJoystick returns true if joystick is configured.
func (c *BoardConfig) HasJoystick() bool {
	return c.Joystick != nil
}

// HasOLED returns true if OLED is configured.
func (c *BoardConfig) HasOLED() bool {
	return c.OLED != nil
}

// HasSplit returns true if split keyboard is configured.
func (c *BoardConfig) HasSplit() bool {
	return c.Split != nil
}

// IsSplitMaster returns true if this is the master side of a split keyboard.
func (c *BoardConfig) IsSplitMaster() bool {
	return c.HasSplit() && c.Split.IsMaster
}

// IsSplitSlave returns true if this is the slave side of a split keyboard.
func (c *BoardConfig) IsSplitSlave() bool {
	return c.HasSplit() && !c.Split.IsMaster
}

// HasEncoders returns true if encoders are configured.
func (c *BoardConfig) HasEncoders() bool {
	return len(c.Encoders) > 0
}

// HasLED returns true if LED is configured.
func (c *BoardConfig) HasLED() bool {
	return c.LED != nil
}

// HasMacros returns true if macros are configured.
func (c *BoardConfig) HasMacros() bool {
	return len(c.Macros) > 0
}

// HasCombos returns true if combos are configured.
func (c *BoardConfig) HasCombos() bool {
	return len(c.Combos) > 0
}

// HasStorage returns true if flash storage is enabled.
func (c *BoardConfig) HasStorage() bool {
	return c.Storage
}

// MacroDef はマクロ定義（設定用）。
type MacroDef struct {
	ID    uint8
	Steps []MacroStep
}
