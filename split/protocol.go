// Package split provides split keyboard communication via UART.
package split

// Message types
const (
	MsgTypeMatrixState byte = 0x01 // Slave -> Master: matrix state
	MsgTypeJoystick    byte = 0x02 // Slave -> Master: joystick data
	MsgTypeKeyEvent    byte = 0x03 // Slave -> Master: key event (Phase 3)
	MsgTypeLEDSync     byte = 0x10 // Master -> Slave: LED sync (Phase 2)
	MsgTypeOLEDSync    byte = 0x11 // Master -> Slave: OLED sync (Phase 2)
	MsgTypeLEDMode     byte = 0x12 // Master -> Slave: LED mode (Phase 3)
)

// Protocol constants
const (
	HeaderByte  byte = 0xFF
	MaxDataSize      = 64 // Maximum data payload size
)

// Message represents a split keyboard protocol message.
type Message struct {
	Type     byte
	Data     []byte
	Checksum byte
}

// Encode encodes a message into a byte slice.
// Format: [Header: 0xFF] [Type: 1byte] [Length: 1byte] [Data: N bytes] [Checksum: 1byte]
func (m *Message) Encode() []byte {
	length := byte(len(m.Data))
	buf := make([]byte, 4+length)

	buf[0] = HeaderByte
	buf[1] = m.Type
	buf[2] = length

	copy(buf[3:], m.Data)

	// Calculate checksum (simple XOR)
	checksum := m.Type ^ length
	for _, b := range m.Data {
		checksum ^= b
	}
	buf[3+length] = checksum

	return buf
}

// EncodeInto は既存バッファにメッセージをエンコードし、使用バイト数を返す。
// bufは少なくとも 4+len(m.Data) バイトの容量が必要。
func (m *Message) EncodeInto(buf []byte) int {
	length := byte(len(m.Data))
	totalSize := 4 + int(length)

	buf[0] = HeaderByte
	buf[1] = m.Type
	buf[2] = length

	copy(buf[3:], m.Data)

	// チェックサム計算（XOR）
	checksum := m.Type ^ length
	for _, b := range m.Data {
		checksum ^= b
	}
	buf[3+length] = checksum

	return totalSize
}

// Decode decodes a message from a byte slice.
// Returns (message, bytesConsumed, error).
func Decode(buf []byte) (*Message, int, error) {
	if len(buf) < 4 {
		return nil, 0, ErrInsufficientData
	}

	// Check header
	if buf[0] != HeaderByte {
		return nil, 0, ErrInvalidHeader
	}

	msgType := buf[1]
	length := int(buf[2])

	// Check if we have enough data
	totalSize := 4 + length
	if len(buf) < totalSize {
		return nil, 0, ErrInsufficientData
	}

	// Extract data
	data := make([]byte, length)
	copy(data, buf[3:3+length])

	// Verify checksum
	receivedChecksum := buf[3+length]
	calculatedChecksum := msgType ^ buf[2]
	for _, b := range data {
		calculatedChecksum ^= b
	}

	if receivedChecksum != calculatedChecksum {
		return nil, totalSize, ErrInvalidChecksum
	}

	msg := &Message{
		Type:     msgType,
		Data:     data,
		Checksum: receivedChecksum,
	}

	return msg, totalSize, nil
}

// EncodeMatrixState encodes a matrix state message.
// Data format: [rows][cols][state bytes...]
// State is packed as bits: 1 bit per key
func EncodeMatrixState(state [][]bool) *Message {
	if len(state) == 0 {
		return &Message{Type: MsgTypeMatrixState, Data: []byte{0, 0}}
	}

	rows := len(state)
	cols := len(state[0])

	// Calculate number of bytes needed for state bits
	totalKeys := rows * cols
	stateBytes := (totalKeys + 7) / 8

	data := make([]byte, 2+stateBytes)
	data[0] = byte(rows)
	data[1] = byte(cols)

	// Pack state into bits
	bitIndex := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if state[r][c] {
				byteIndex := bitIndex / 8
				bitOffset := uint(bitIndex % 8)
				data[2+byteIndex] |= 1 << bitOffset
			}
			bitIndex++
		}
	}

	return &Message{
		Type: MsgTypeMatrixState,
		Data: data,
	}
}

// DecodeMatrixState decodes a matrix state message.
func DecodeMatrixState(msg *Message) ([][]bool, error) {
	if msg.Type != MsgTypeMatrixState {
		return nil, ErrInvalidMessageType
	}

	if len(msg.Data) < 2 {
		return nil, ErrInsufficientData
	}

	rows := int(msg.Data[0])
	cols := int(msg.Data[1])

	if rows == 0 || cols == 0 {
		return make([][]bool, 0), nil
	}

	// Unpack state bits
	state := make([][]bool, rows)
	for i := range state {
		state[i] = make([]bool, cols)
	}

	bitIndex := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			byteIndex := bitIndex / 8
			bitOffset := uint(bitIndex % 8)

			if 2+byteIndex < len(msg.Data) {
				state[r][c] = (msg.Data[2+byteIndex] & (1 << bitOffset)) != 0
			}
			bitIndex++
		}
	}

	return state, nil
}

// DecodeMatrixStateInto はマトリクス状態を既存の2D boolスライスに直接デコードする。
// dstは事前割り当て済みの [][]bool で、行数・列数がメッセージと一致する前提。
func DecodeMatrixStateInto(msg *Message, dst [][]bool) error {
	if msg.Type != MsgTypeMatrixState {
		return ErrInvalidMessageType
	}

	if len(msg.Data) < 2 {
		return ErrInsufficientData
	}

	rows := int(msg.Data[0])
	cols := int(msg.Data[1])

	if rows == 0 || cols == 0 {
		return nil
	}

	// dstのゼロクリアとビットのアンパック
	bitIndex := 0
	for r := 0; r < rows && r < len(dst); r++ {
		for c := 0; c < cols && c < len(dst[r]); c++ {
			byteIndex := bitIndex / 8
			bitOffset := uint(bitIndex % 8)

			if 2+byteIndex < len(msg.Data) {
				dst[r][c] = (msg.Data[2+byteIndex] & (1 << bitOffset)) != 0
			} else {
				dst[r][c] = false
			}
			bitIndex++
		}
	}

	return nil
}

// EncodeJoystick encodes a joystick data message.
// Data format: [X: int8][Y: int8]
func EncodeJoystick(x, y int8) *Message {
	return &Message{
		Type: MsgTypeJoystick,
		Data: []byte{byte(x), byte(y)},
	}
}

// DecodeJoystick decodes a joystick data message.
func DecodeJoystick(msg *Message) (x, y int8, err error) {
	if msg.Type != MsgTypeJoystick {
		return 0, 0, ErrInvalidMessageType
	}

	if len(msg.Data) < 2 {
		return 0, 0, ErrInsufficientData
	}

	x = int8(msg.Data[0])
	y = int8(msg.Data[1])
	return x, y, nil
}

// EncodeLEDSync encodes an LED synchronization message (Phase 2 M3).
// Data format: [count: 1byte][R1][G1][B1][R2][G2][B2]...
// Master -> Slave: Send LED color data to slave
func EncodeLEDSync(colors []byte) *Message {
	if len(colors) > MaxDataSize-1 {
		colors = colors[:MaxDataSize-1]
	}

	count := byte(len(colors) / 3)
	data := make([]byte, 1+len(colors))
	data[0] = count
	copy(data[1:], colors)

	return &Message{
		Type: MsgTypeLEDSync,
		Data: data,
	}
}

// DecodeLEDSync decodes an LED synchronization message (Phase 2 M3).
// Returns RGB color data as byte slice (R, G, B, R, G, B, ...)
func DecodeLEDSync(msg *Message) ([]byte, error) {
	if msg.Type != MsgTypeLEDSync {
		return nil, ErrInvalidMessageType
	}

	if len(msg.Data) < 1 {
		return nil, ErrInsufficientData
	}

	count := int(msg.Data[0])
	expectedLen := 1 + count*3

	if len(msg.Data) < expectedLen {
		return nil, ErrInsufficientData
	}

	colors := make([]byte, count*3)
	copy(colors, msg.Data[1:expectedLen])

	return colors, nil
}

// EncodeOLEDSync encodes an OLED synchronization message (Phase 2 M3).
// Data format: [type: 1byte][payload...]
// Type 0x00: Clear screen
// Type 0x01: Text at position [row][col][textlen][text...]
// Type 0x02: Raw buffer (for full screen updates)
func EncodeOLEDSync(syncType byte, payload []byte) *Message {
	data := make([]byte, 1+len(payload))
	data[0] = syncType
	copy(data[1:], payload)

	return &Message{
		Type: MsgTypeOLEDSync,
		Data: data,
	}
}

// DecodeOLEDSync decodes an OLED synchronization message (Phase 2 M3).
// Returns sync type and payload data.
func DecodeOLEDSync(msg *Message) (syncType byte, payload []byte, err error) {
	if msg.Type != MsgTypeOLEDSync {
		return 0, nil, ErrInvalidMessageType
	}

	if len(msg.Data) < 1 {
		return 0, nil, ErrInsufficientData
	}

	syncType = msg.Data[0]
	if len(msg.Data) > 1 {
		payload = make([]byte, len(msg.Data)-1)
		copy(payload, msg.Data[1:])
	}

	return syncType, payload, nil
}

// OLED sync types (Phase 2 M3)
const (
	OLEDSyncClear  byte = 0x00 // Clear display
	OLEDSyncText   byte = 0x01 // Text at position
	OLEDSyncBuffer byte = 0x02 // Raw buffer update
)

// EncodeKeyEvent はキーイベントメッセージをエンコードする（Phase 3 M2）。
// データ形式: [row:1][col:1][pressed:1]
// Slave -> Master: スレーブ側のキー状態変化を通知
func EncodeKeyEvent(row, col uint8, pressed bool) *Message {
	var p byte
	if pressed {
		p = 1
	}
	return &Message{
		Type: MsgTypeKeyEvent,
		Data: []byte{row, col, p},
	}
}

// DecodeKeyEvent はキーイベントメッセージをデコードする（Phase 3 M2）。
func DecodeKeyEvent(msg *Message) (row, col uint8, pressed bool, err error) {
	if msg.Type != MsgTypeKeyEvent {
		return 0, 0, false, ErrInvalidMessageType
	}

	if len(msg.Data) < 3 {
		return 0, 0, false, ErrInsufficientData
	}

	row = msg.Data[0]
	col = msg.Data[1]
	pressed = msg.Data[2] != 0
	return row, col, pressed, nil
}

// EncodeLEDMode はLEDモードメッセージをエンコードする（Phase 3 M2）。
// データ形式: [mode:1]
// Master -> Slave: LEDの同期モードを通知
func EncodeLEDMode(mode uint8) *Message {
	return &Message{
		Type: MsgTypeLEDMode,
		Data: []byte{mode},
	}
}

// DecodeLEDMode はLEDモードメッセージをデコードする（Phase 3 M2）。
func DecodeLEDMode(msg *Message) (mode uint8, err error) {
	if msg.Type != MsgTypeLEDMode {
		return 0, ErrInvalidMessageType
	}

	if len(msg.Data) < 1 {
		return 0, ErrInsufficientData
	}

	return msg.Data[0], nil
}
