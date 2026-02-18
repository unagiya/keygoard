package split

import (
	"machine"
	"time"
)

// SlaveConfig holds the configuration for split keyboard slave side.
type SlaveConfig struct {
	UART     *machine.UART
	BaudRate uint32        // Default: 460800
	Interval time.Duration // Send interval, default: 1ms
}

// Slave represents the slave side of a split keyboard.
type Slave struct {
	uart         *machine.UART
	interval     time.Duration
	lastSent     time.Time
	rxBuffer     []byte                              // Receive buffer for master->slave messages (Phase 2 M3)
	txBuffer     [68]byte                            // 送信用事前割り当てバッファ（4ヘッダ+64データ）
	ledCallback  func(colors []byte)                 // Callback for LED sync (Phase 2 M3)
	oledCallback func(syncType byte, payload []byte) // Callback for OLED sync (Phase 2 M3)
	ledSyncMode  uint8                               // LED同期モード（Phase 3 M2）
}

// NewSlave creates a new split keyboard slave.
func NewSlave(cfg *SlaveConfig) (*Slave, error) {
	// Default values
	baudRate := cfg.BaudRate
	if baudRate == 0 {
		baudRate = 460800
	}

	interval := cfg.Interval
	if interval == 0 {
		interval = 1 * time.Millisecond
	}

	// Configure UART
	err := cfg.UART.Configure(machine.UARTConfig{
		BaudRate: baudRate,
	})
	if err != nil {
		return nil, err
	}

	return &Slave{
		uart:     cfg.UART,
		interval: interval,
		lastSent: time.Now(),
		rxBuffer: make([]byte, 256),
	}, nil
}

// SendMatrixState sends the matrix state to the master.
func (s *Slave) SendMatrixState(state [][]bool) error {
	// Check if enough time has passed since last send
	if time.Since(s.lastSent) < s.interval {
		return nil
	}

	msg := EncodeMatrixState(state)
	n := msg.EncodeInto(s.txBuffer[:])

	_, err := s.uart.Write(s.txBuffer[:n])
	if err != nil {
		return err
	}

	s.lastSent = time.Now()
	return nil
}

// SendJoystick sends joystick data to the master.
func (s *Slave) SendJoystick(x, y int8) error {
	msg := EncodeJoystick(x, y)
	n := msg.EncodeInto(s.txBuffer[:])

	_, err := s.uart.Write(s.txBuffer[:n])
	return err
}

// SendMatrixAndJoystick sends both matrix state and joystick data.
// This is more efficient than sending separately.
func (s *Slave) SendMatrixAndJoystick(state [][]bool, joyX, joyY int8) error {
	// Check if enough time has passed since last send
	if time.Since(s.lastSent) < s.interval {
		return nil
	}

	// Send matrix state
	msg1 := EncodeMatrixState(state)
	n := msg1.EncodeInto(s.txBuffer[:])
	_, err := s.uart.Write(s.txBuffer[:n])
	if err != nil {
		return err
	}

	// Send joystick data
	msg2 := EncodeJoystick(joyX, joyY)
	n = msg2.EncodeInto(s.txBuffer[:])
	_, err = s.uart.Write(s.txBuffer[:n])
	if err != nil {
		return err
	}

	s.lastSent = time.Now()
	return nil
}

// Update receives and processes messages from the master (Phase 2 M3).
// Should be called regularly (e.g., every 1ms) to receive LED/OLED sync.
func (s *Slave) Update() error {
	// Check for available data
	available := s.uart.Buffered()
	if available == 0 {
		return nil
	}

	// Read available data
	n, err := s.uart.Read(s.rxBuffer[:available])
	if err != nil {
		return err
	}

	// Try to decode messages
	offset := 0
	for offset < n {
		msg, consumed, err := Decode(s.rxBuffer[offset:n])
		if err == ErrInsufficientData {
			break
		}
		if err != nil {
			// Invalid message, skip one byte and try again
			offset++
			continue
		}

		// Process message
		s.processMessage(msg)

		offset += consumed
	}

	return nil
}

// processMessage processes a received message from the master (Phase 2 M3).
func (s *Slave) processMessage(msg *Message) {
	switch msg.Type {
	case MsgTypeLEDSync:
		// Independentモード時はLED同期をスキップ（Phase 3 M2）
		if s.ledSyncMode == 1 {
			return
		}
		if s.ledCallback != nil {
			colors, err := DecodeLEDSync(msg)
			if err == nil {
				s.ledCallback(colors)
			}
		}

	case MsgTypeLEDMode:
		mode, err := DecodeLEDMode(msg)
		if err == nil {
			s.ledSyncMode = mode
		}

	case MsgTypeOLEDSync:
		if s.oledCallback != nil {
			syncType, payload, err := DecodeOLEDSync(msg)
			if err == nil {
				s.oledCallback(syncType, payload)
			}
		}
	}
}

// SetLEDCallback sets the callback for LED sync messages (Phase 2 M3).
func (s *Slave) SetLEDCallback(cb func(colors []byte)) {
	s.ledCallback = cb
}

// SetOLEDCallback sets the callback for OLED sync messages (Phase 2 M3).
func (s *Slave) SetOLEDCallback(cb func(syncType byte, payload []byte)) {
	s.oledCallback = cb
}

// SendKeyEvent はキーイベントをマスターに送信する（Phase 3 M2）。
func (s *Slave) SendKeyEvent(row, col uint8, pressed bool) error {
	msg := EncodeKeyEvent(row, col, pressed)
	n := msg.EncodeInto(s.txBuffer[:])
	_, err := s.uart.Write(s.txBuffer[:n])
	return err
}

// GetLEDSyncMode はマスターから受信したLED同期モードを返す（Phase 3 M2）。
func (s *Slave) GetLEDSyncMode() uint8 {
	return s.ledSyncMode
}
