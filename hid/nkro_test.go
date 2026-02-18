package hid

import "testing"

func TestNKROReport_AddKey(t *testing.T) {
	r := &NKROReport{}

	// 正常なキー追加
	if !r.AddKey(0x04) { // KC_A
		t.Error("AddKey(0x04) should return true")
	}
	if !r.HasKey(0x04) {
		t.Error("HasKey(0x04) should return true after AddKey")
	}

	// 複数キー追加
	if !r.AddKey(0x05) { // KC_B
		t.Error("AddKey(0x05) should return true")
	}
	if !r.AddKey(0x06) { // KC_C
		t.Error("AddKey(0x06) should return true")
	}

	// 範囲外のキーは失敗
	if r.AddKey(0x03) {
		t.Error("AddKey(0x03) should return false (below range)")
	}
	if r.AddKey(0x7C) {
		t.Error("AddKey(0x7C) should return false (above range)")
	}
}

func TestNKROReport_RemoveKey(t *testing.T) {
	r := &NKROReport{}
	r.AddKey(0x04)
	r.AddKey(0x05)

	r.RemoveKey(0x04)
	if r.HasKey(0x04) {
		t.Error("HasKey(0x04) should return false after RemoveKey")
	}
	if !r.HasKey(0x05) {
		t.Error("HasKey(0x05) should still return true")
	}
}

func TestNKROReport_ManyKeys(t *testing.T) {
	r := &NKROReport{}

	// 6キー以上の同時押し（NKROの利点）
	keys := []uint8{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}
	for _, k := range keys {
		if !r.AddKey(k) {
			t.Errorf("AddKey(0x%02X) should return true", k)
		}
	}

	// 全キーが存在することを確認
	for _, k := range keys {
		if !r.HasKey(k) {
			t.Errorf("HasKey(0x%02X) should return true", k)
		}
	}
}

func TestNKROReport_ToBytes(t *testing.T) {
	r := &NKROReport{}
	r.Modifier = 0x03 // LCtrl + LShift
	r.AddKey(0x04)    // KC_A → bit 0 of byte 0

	bytes := r.ToBytes()
	if len(bytes) != 17 {
		t.Errorf("ToBytes() length should be 17, got %d", len(bytes))
	}
	if bytes[0] != 0x03 {
		t.Errorf("Modifier byte should be 0x03, got 0x%02X", bytes[0])
	}
	if bytes[1] != 0x00 {
		t.Errorf("Reserved byte should be 0x00, got 0x%02X", bytes[1])
	}
	// 0x04 → bit 0 of Keys[0] → bytes[2]
	if bytes[2] != 0x01 {
		t.Errorf("Keys byte 0 should be 0x01, got 0x%02X", bytes[2])
	}
}

func TestNKROReport_Clear(t *testing.T) {
	r := &NKROReport{}
	r.Modifier = 0xFF
	r.AddKey(0x04)
	r.AddKey(0x50)

	r.Clear()
	if r.Modifier != 0 {
		t.Error("Modifier should be 0 after Clear")
	}
	if r.HasKey(0x04) {
		t.Error("HasKey(0x04) should return false after Clear")
	}
	if r.HasKey(0x50) {
		t.Error("HasKey(0x50) should return false after Clear")
	}
}

func TestNKROReport_BitmapAccuracy(t *testing.T) {
	r := &NKROReport{}

	// 各バイト境界のキーをテスト
	testCases := []struct {
		keycode  uint8
		byteIdx  int
		expected byte
	}{
		{0x04, 0, 0x01}, // bit 0 of byte 0
		{0x0B, 0, 0x80}, // bit 7 of byte 0
		{0x0C, 1, 0x01}, // bit 0 of byte 1
		{0x13, 1, 0x80}, // bit 7 of byte 1
	}

	for _, tc := range testCases {
		r.Clear()
		r.AddKey(tc.keycode)
		if r.Keys[tc.byteIdx] != tc.expected {
			t.Errorf("keycode 0x%02X: Keys[%d] should be 0x%02X, got 0x%02X",
				tc.keycode, tc.byteIdx, tc.expected, r.Keys[tc.byteIdx])
		}
	}
}
