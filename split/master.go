package split

import (
	"machine"
	"time"
)

// MasterConfig holds the configuration for split keyboard master side.
type MasterConfig struct {
	UART      *machine.UART
	BaudRate  uint32        // Default: 460800
	Timeout   time.Duration // Default: 5ms
	SlaveRows int           // Number of rows on slave side
	SlaveCols int           // Number of columns on slave side
}

// Master represents the master side of a split keyboard.
type Master struct {
	uart             *machine.UART
	timeout          time.Duration
	slaveRows        int
	slaveCols        int
	slaveState       [][]bool // Last received slave matrix state
	slaveJoyX        int8
	slaveJoyY        int8
	rxBuffer         []byte
	txBuffer         [68]byte // 送信用事前割り当てバッファ（4ヘッダ+64データ）
	lastRxTime       time.Time
	keyEventCallback func(row, col uint8, pressed bool) // キーイベントコールバック（Phase 3 M2）
}

// NewMaster creates a new split keyboard master.
func NewMaster(cfg *MasterConfig) (*Master, error) {
	// Default values
	baudRate := cfg.BaudRate
	if baudRate == 0 {
		baudRate = 460800
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Millisecond
	}

	// Configure UART
	err := cfg.UART.Configure(machine.UARTConfig{
		BaudRate: baudRate,
	})
	if err != nil {
		return nil, err
	}

	// Initialize slave state
	slaveState := make([][]bool, cfg.SlaveRows)
	for i := range slaveState {
		slaveState[i] = make([]bool, cfg.SlaveCols)
	}

	return &Master{
		uart:       cfg.UART,
		timeout:    timeout,
		slaveRows:  cfg.SlaveRows,
		slaveCols:  cfg.SlaveCols,
		slaveState: slaveState,
		rxBuffer:   make([]byte, 256),
		lastRxTime: time.Now(),
	}, nil
}

// Update receives and processes messages from the slave.
// Should be called regularly (e.g., every 1ms).
func (m *Master) Update() error {
	// Check for available data
	available := m.uart.Buffered()
	if available == 0 {
		// Check timeout
		if time.Since(m.lastRxTime) > m.timeout {
			// Slave not responding, clear slave state
			for r := range m.slaveState {
				for c := range m.slaveState[r] {
					m.slaveState[r][c] = false
				}
			}
			m.slaveJoyX = 0
			m.slaveJoyY = 0
		}
		return nil
	}

	// Read available data
	n, err := m.uart.Read(m.rxBuffer[:available])
	if err != nil {
		return err
	}

	// Try to decode messages
	offset := 0
	for offset < n {
		msg, consumed, err := Decode(m.rxBuffer[offset:n])
		if err == ErrInsufficientData {
			break
		}
		if err != nil {
			// Invalid message, skip one byte and try again
			offset++
			continue
		}

		// Process message
		m.processMessage(msg)
		m.lastRxTime = time.Now()

		offset += consumed
	}

	return nil
}

// processMessage processes a received message from the slave.
func (m *Master) processMessage(msg *Message) {
	switch msg.Type {
	case MsgTypeMatrixState:
		// 事前割り当て済みのslaveStateに直接デコード
		err := DecodeMatrixStateInto(msg, m.slaveState)
		_ = err

	case MsgTypeJoystick:
		x, y, err := DecodeJoystick(msg)
		if err == nil {
			m.slaveJoyX = x
			m.slaveJoyY = y
		}

	case MsgTypeKeyEvent:
		row, col, pressed, err := DecodeKeyEvent(msg)
		if err == nil && m.keyEventCallback != nil {
			m.keyEventCallback(row, col, pressed)
		}
	}
}

// GetSlaveState returns the current slave matrix state.
func (m *Master) GetSlaveState() [][]bool {
	return m.slaveState
}

// GetSlaveJoystick returns the current slave joystick values.
func (m *Master) GetSlaveJoystick() (x, y int8) {
	return m.slaveJoyX, m.slaveJoyY
}

// IsSlaveConnected returns true if the slave has sent data recently.
func (m *Master) IsSlaveConnected() bool {
	return time.Since(m.lastRxTime) <= m.timeout
}

// SendLEDSync sends LED color data to the slave (Phase 2 M3).
// colors should be RGB byte array (R, G, B, R, G, B, ...)
func (m *Master) SendLEDSync(colors []byte) error {
	msg := EncodeLEDSync(colors)
	n := msg.EncodeInto(m.txBuffer[:])
	_, err := m.uart.Write(m.txBuffer[:n])
	return err
}

// SendOLEDSync sends OLED synchronization data to the slave (Phase 2 M3).
func (m *Master) SendOLEDSync(syncType byte, payload []byte) error {
	msg := EncodeOLEDSync(syncType, payload)
	n := msg.EncodeInto(m.txBuffer[:])
	_, err := m.uart.Write(m.txBuffer[:n])
	return err
}

// SetKeyEventCallback はスレーブからのキーイベント受信時のコールバックを設定する（Phase 3 M2）。
func (m *Master) SetKeyEventCallback(cb func(row, col uint8, pressed bool)) {
	m.keyEventCallback = cb
}

// SendLEDMode はLED同期モードをスレーブに送信する（Phase 3 M2）。
func (m *Master) SendLEDMode(mode uint8) error {
	msg := EncodeLEDMode(mode)
	n := msg.EncodeInto(m.txBuffer[:])
	_, err := m.uart.Write(m.txBuffer[:n])
	return err
}

// SendOLEDClear sends a command to clear the slave's OLED display (Phase 2 M3).
func (m *Master) SendOLEDClear() error {
	return m.SendOLEDSync(OLEDSyncClear, nil)
}

// SendOLEDText sends text to display at a specific position on slave's OLED (Phase 2 M3).
func (m *Master) SendOLEDText(row, col uint8, text string) error {
	textBytes := []byte(text)
	if len(textBytes) > MaxDataSize-3 {
		textBytes = textBytes[:MaxDataSize-3]
	}

	payload := make([]byte, 3+len(textBytes))
	payload[0] = row
	payload[1] = col
	payload[2] = uint8(len(textBytes))
	copy(payload[3:], textBytes)

	return m.SendOLEDSync(OLEDSyncText, payload)
}
