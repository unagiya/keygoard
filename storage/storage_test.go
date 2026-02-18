package storage

import (
	"testing"
)

// mockFlash はテスト用のFlashデバイスモック。
type mockFlash struct {
	data [(RequiredSectors + 1) * SectorSize]byte // オフセットテスト用に余裕を持たせる
}

func (m *mockFlash) ReadAt(buf []byte, off int64) (int, error) {
	copy(buf, m.data[off:off+int64(len(buf))])
	return len(buf), nil
}

func (m *mockFlash) WriteAt(p []byte, off int64) (int, error) {
	copy(m.data[off:], p)
	return len(p), nil
}

func (m *mockFlash) EraseBlocks(start, length int64) error {
	// セクター消去（0xFFで埋める）
	for i := start * int64(SectorSize); i < (start+length)*int64(SectorSize); i++ {
		if i < int64(len(m.data)) {
			m.data[i] = 0xFF
		}
	}
	return nil
}

func TestStorage_New(t *testing.T) {
	flash := &mockFlash{}
	s, err := New(flash, 0)
	if err != nil {
		t.Fatalf("New should succeed: %v", err)
	}
	if s == nil {
		t.Fatal("Storage should not be nil")
	}
}

func TestStorage_NewNilFlash(t *testing.T) {
	_, err := New(nil, 0)
	if err != ErrNotInitialized {
		t.Errorf("New with nil flash should return ErrNotInitialized, got %v", err)
	}
}

func TestStorage_IsValid_Empty(t *testing.T) {
	flash := &mockFlash{}
	s, _ := New(flash, 0)
	if s.IsValid() {
		t.Error("empty storage should not be valid")
	}
}

func TestStorage_SaveAllAndIsValid(t *testing.T) {
	flash := &mockFlash{}
	s, _ := New(flash, 0)

	keymapData := make([]byte, 100)
	keymapData[0] = 0x42
	configData := make([]byte, 10)
	configData[0] = 0x01
	macroData := make([]byte, 96)

	err := s.SaveAll(keymapData, configData, macroData)
	if err != nil {
		t.Fatalf("SaveAll should succeed: %v", err)
	}

	if !s.IsValid() {
		t.Error("storage should be valid after SaveAll")
	}
}

func TestStorage_SaveAllAndLoadAll(t *testing.T) {
	flash := &mockFlash{}
	s, _ := New(flash, 0)

	// テストデータ
	keymapData := make([]byte, 200)
	for i := range keymapData {
		keymapData[i] = byte(i)
	}
	configData := make([]byte, 10)
	configData[0] = 0xAB
	macroData := make([]byte, 96)
	macroData[0] = 0xCD

	// 保存
	err := s.SaveAll(keymapData, configData, macroData)
	if err != nil {
		t.Fatalf("SaveAll: %v", err)
	}

	// 読み出し
	loadedKeymap, loadedConfig, loadedMacro, err := s.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}

	// キーマップ検証（パディングあり）
	for i := 0; i < len(keymapData); i++ {
		if loadedKeymap[i] != keymapData[i] {
			t.Errorf("keymap[%d]: got 0x%02X, want 0x%02X", i, loadedKeymap[i], keymapData[i])
			break
		}
	}

	// 設定検証
	if loadedConfig[0] != 0xAB {
		t.Errorf("config[0]: got 0x%02X, want 0xAB", loadedConfig[0])
	}

	// マクロ検証
	if loadedMacro[0] != 0xCD {
		t.Errorf("macro[0]: got 0x%02X, want 0xCD", loadedMacro[0])
	}
}

func TestStorage_SaveAll_DataTooLarge(t *testing.T) {
	flash := &mockFlash{}
	s, _ := New(flash, 0)

	// キーマップが大きすぎる
	largeData := make([]byte, KeymapMaxSize+1)
	err := s.SaveAll(largeData, nil, nil)
	if err != ErrDataTooLarge {
		t.Errorf("should return ErrDataTooLarge, got %v", err)
	}
}

func TestCalculateChecksum(t *testing.T) {
	data1 := []byte{0x01, 0x02, 0x03}
	data2 := []byte{0x01, 0x02, 0x03}

	c1 := CalculateChecksum(data1)
	c2 := CalculateChecksum(data2)
	if c1 != c2 {
		t.Error("same data should produce same checksum")
	}

	data3 := []byte{0x01, 0x02, 0x04}
	c3 := CalculateChecksum(data3)
	if c1 == c3 {
		t.Error("different data should produce different checksum")
	}
}

func TestStorage_WithOffset(t *testing.T) {
	flash := &mockFlash{}
	// オフセット4096（セクター1から開始）
	s, _ := New(flash, 4096)

	keymapData := []byte{0x42}
	configData := []byte{0x01}
	macroData := []byte{0x02}

	err := s.SaveAll(keymapData, configData, macroData)
	if err != nil {
		t.Fatalf("SaveAll with offset: %v", err)
	}

	if !s.IsValid() {
		t.Error("storage with offset should be valid")
	}
}
