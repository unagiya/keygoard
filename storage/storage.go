// Package storage はFlash直接操作によるキーマップ・設定の永続化を提供する。
// tinyfs不使用、依存関係追加なし。
package storage

import "errors"

// ストレージエラー
var (
	ErrInvalidMagic    = errors.New("storage: invalid magic")
	ErrInvalidChecksum = errors.New("storage: invalid checksum")
	ErrDataTooLarge    = errors.New("storage: data too large")
	ErrNotInitialized  = errors.New("storage: not initialized")
)

// ストレージレイアウト定数
const (
	HeaderSize       = 16         // StorageHeaderサイズ
	MagicSize        = 4          // マジックバイト
	KeymapOffset     = HeaderSize // キーマップデータのオフセット
	KeymapMaxSize    = 8192       // キーマップ最大サイズ（8KB）
	ConfigOffset     = 8208       // 設定データのオフセット（16+8192）
	ConfigMaxSize    = 64         // 設定最大サイズ
	MacroOffset      = 8272       // マクロデータのオフセット（8208+64）
	MacroMaxSize     = 1536       // マクロ最大サイズ（16×32×3）
	TotalStorageSize = 9808       // 合計サイズ
	SectorSize       = 4096       // Flashセクターサイズ
	RequiredSectors  = 3          // 必要セクター数（12KB）
	StorageVersion   = 1          // 現在のバージョン
)

// Magic はストレージのマジックバイト。
var Magic = [MagicSize]byte{'K', 'B', 'D', 0x00}

// StorageHeader はFlashデータのヘッダ。
type StorageHeader struct {
	Magic    [MagicSize]byte // "KBD\x00"
	Version  uint8           // ストレージバージョン
	Reserved [3]byte         // 予約
	DataSize uint32          // データサイズ（ヘッダ除く）
	Checksum uint32          // XOR回転チェックサム
}

// BlockDevice はFlashデバイスのインターフェース。
// machine.BlockDeviceと同じシグネチャ。
type BlockDevice interface {
	ReadAt(buf []byte, off int64) (int, error)
	WriteAt(p []byte, off int64) (int, error)
	EraseBlocks(start, length int64) error
}

// Storage はFlash永続化ストレージ。
type Storage struct {
	flash  BlockDevice
	offset int64
	size   int64
}

// New は新しいStorageを作成する。
// flashはBlockDeviceインターフェースを実装するデバイス。
// offsetはFlash内のストレージ開始位置。
func New(flash BlockDevice, offset int64) (*Storage, error) {
	if flash == nil {
		return nil, ErrNotInitialized
	}
	return &Storage{
		flash:  flash,
		offset: offset,
		size:   int64(RequiredSectors * SectorSize),
	}, nil
}

// Read はFlashからデータを読み出す。
func (s *Storage) Read(buf []byte, offset int64) error {
	_, err := s.flash.ReadAt(buf, s.offset+offset)
	return err
}

// Write はFlashにデータを書き込む（消去→書き込み自動）。
// 影響するセクターを消去してから書き込む。
func (s *Storage) Write(data []byte, offset int64) error {
	absOffset := s.offset + offset

	// 影響するセクター範囲を計算
	startSector := absOffset / int64(SectorSize)
	endSector := (absOffset + int64(len(data)) - 1) / int64(SectorSize)
	numSectors := endSector - startSector + 1

	// セクターを消去
	if err := s.flash.EraseBlocks(startSector, numSectors); err != nil {
		return err
	}

	// データを書き込み
	_, err := s.flash.WriteAt(data, absOffset)
	return err
}

// ReadHeader はストレージヘッダを読み出す。
func (s *Storage) ReadHeader() (*StorageHeader, error) {
	buf := make([]byte, HeaderSize)
	if err := s.Read(buf, 0); err != nil {
		return nil, err
	}

	header := &StorageHeader{}
	copy(header.Magic[:], buf[0:4])
	header.Version = buf[4]
	header.DataSize = uint32(buf[8]) | uint32(buf[9])<<8 | uint32(buf[10])<<16 | uint32(buf[11])<<24
	header.Checksum = uint32(buf[12]) | uint32(buf[13])<<8 | uint32(buf[14])<<16 | uint32(buf[15])<<24
	return header, nil
}

// WriteHeader はストレージヘッダを書き込む。
func (s *Storage) WriteHeader(header *StorageHeader) error {
	buf := make([]byte, HeaderSize)
	copy(buf[0:4], header.Magic[:])
	buf[4] = header.Version
	buf[8] = byte(header.DataSize)
	buf[9] = byte(header.DataSize >> 8)
	buf[10] = byte(header.DataSize >> 16)
	buf[11] = byte(header.DataSize >> 24)
	buf[12] = byte(header.Checksum)
	buf[13] = byte(header.Checksum >> 8)
	buf[14] = byte(header.Checksum >> 16)
	buf[15] = byte(header.Checksum >> 24)
	return s.Write(buf, 0)
}

// IsValid はストレージが有効か（マジック+チェックサム）を検証する。
func (s *Storage) IsValid() bool {
	header, err := s.ReadHeader()
	if err != nil {
		return false
	}

	// マジック検証
	if header.Magic != Magic {
		return false
	}

	// バージョン検証
	if header.Version != StorageVersion {
		return false
	}

	// データ読み出しとチェックサム検証
	if header.DataSize == 0 || header.DataSize > TotalStorageSize-HeaderSize {
		return false
	}

	data := make([]byte, header.DataSize)
	if err := s.Read(data, HeaderSize); err != nil {
		return false
	}

	checksum := CalculateChecksum(data)
	return checksum == header.Checksum
}

// CalculateChecksum はXOR回転チェックサムを計算する。
func CalculateChecksum(data []byte) uint32 {
	var checksum uint32
	for _, b := range data {
		checksum ^= uint32(b)
		// 左回転
		checksum = (checksum << 1) | (checksum >> 31)
	}
	return checksum
}

// SaveAll は全データ（キーマップ+設定+マクロ）をFlashに保存する。
func (s *Storage) SaveAll(keymapData, configData, macroData []byte) error {
	if len(keymapData) > KeymapMaxSize {
		return ErrDataTooLarge
	}
	if len(configData) > ConfigMaxSize {
		return ErrDataTooLarge
	}
	if len(macroData) > MacroMaxSize {
		return ErrDataTooLarge
	}

	// 全データを1バッファに構築
	totalData := make([]byte, TotalStorageSize-HeaderSize)
	copy(totalData[0:], keymapData)
	copy(totalData[KeymapMaxSize:], configData)
	copy(totalData[KeymapMaxSize+ConfigMaxSize:], macroData)

	// チェックサム計算
	checksum := CalculateChecksum(totalData)

	// ヘッダ書き込み（セクター消去含む）
	header := &StorageHeader{
		Magic:    Magic,
		Version:  StorageVersion,
		DataSize: uint32(len(totalData)),
		Checksum: checksum,
	}

	// 全セクターを消去
	startSector := s.offset / int64(SectorSize)
	if err := s.flash.EraseBlocks(startSector, RequiredSectors); err != nil {
		return err
	}

	// ヘッダを書き込み
	headerBuf := make([]byte, HeaderSize)
	copy(headerBuf[0:4], header.Magic[:])
	headerBuf[4] = header.Version
	headerBuf[8] = byte(header.DataSize)
	headerBuf[9] = byte(header.DataSize >> 8)
	headerBuf[10] = byte(header.DataSize >> 16)
	headerBuf[11] = byte(header.DataSize >> 24)
	headerBuf[12] = byte(header.Checksum)
	headerBuf[13] = byte(header.Checksum >> 8)
	headerBuf[14] = byte(header.Checksum >> 16)
	headerBuf[15] = byte(header.Checksum >> 24)
	if _, err := s.flash.WriteAt(headerBuf, s.offset); err != nil {
		return err
	}

	// データを書き込み
	_, err := s.flash.WriteAt(totalData, s.offset+HeaderSize)
	return err
}

// LoadAll はFlashから全データ（キーマップ+設定+マクロ）を読み出す。
func (s *Storage) LoadAll() (keymapData, configData, macroData []byte, err error) {
	if !s.IsValid() {
		return nil, nil, nil, ErrInvalidMagic
	}

	keymapData = make([]byte, KeymapMaxSize)
	if err := s.Read(keymapData, KeymapOffset); err != nil {
		return nil, nil, nil, err
	}

	configData = make([]byte, ConfigMaxSize)
	if err := s.Read(configData, ConfigOffset); err != nil {
		return nil, nil, nil, err
	}

	macroData = make([]byte, MacroMaxSize)
	if err := s.Read(macroData, MacroOffset); err != nil {
		return nil, nil, nil, err
	}

	return keymapData, configData, macroData, nil
}
