package engine

import (
	"testing"
	"time"

	"github.com/unagiya/keygoard/keycode"
)

func TestComboDetector_RegisterCombo(t *testing.T) {
	cd := NewComboDetector(DefaultComboWindow)

	err := cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)
	if err != nil {
		t.Errorf("RegisterCombo should succeed: %v", err)
	}

	// キー数不足
	err = cd.RegisterCombo([]keycode.Keycode{keycode.KC_A}, keycode.KC_C)
	if err == nil {
		t.Error("RegisterCombo with 1 key should fail")
	}

	// キー数超過
	err = cd.RegisterCombo([]keycode.Keycode{
		keycode.KC_A, keycode.KC_B, keycode.KC_C, keycode.KC_D, keycode.KC_E,
	}, keycode.KC_F)
	if err == nil {
		t.Error("RegisterCombo with 5 keys should fail")
	}
}

func TestComboDetector_RegisterComboMax(t *testing.T) {
	cd := NewComboDetector(DefaultComboWindow)

	// 16コンボ登録
	for i := 0; i < MaxCombos; i++ {
		err := cd.RegisterCombo(
			[]keycode.Keycode{keycode.Keycode(0x04 + i), keycode.Keycode(0x20 + i)},
			keycode.Keycode(0x30+i),
		)
		if err != nil {
			t.Errorf("RegisterCombo %d should succeed: %v", i, err)
		}
	}

	// 17番目は失敗
	err := cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)
	if err == nil {
		t.Error("17th RegisterCombo should fail")
	}
}

func TestComboDetector_TwoKeyCombo(t *testing.T) {
	cd := NewComboDetector(500 * time.Millisecond) // テスト用に長めのウィンドウ

	// A+B → C
	cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)

	// Aを押す → 消費される（保留）
	output, consumed := cd.ProcessKeyPress(keycode.KC_A)
	if output != keycode.KC_NO {
		t.Errorf("first key should not produce output, got 0x%04X", uint16(output))
	}
	if !consumed {
		t.Error("first key should be consumed")
	}

	// Bを押す → コンボ成立、Cが出力
	output, consumed = cd.ProcessKeyPress(keycode.KC_B)
	if output != keycode.KC_C {
		t.Errorf("combo should produce KC_C, got 0x%04X", uint16(output))
	}
	if !consumed {
		t.Error("second key should be consumed")
	}
}

func TestComboDetector_UnrelatedKey(t *testing.T) {
	cd := NewComboDetector(500 * time.Millisecond)

	// A+B → C
	cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)

	// 無関係なキー（D）は消費されない
	output, consumed := cd.ProcessKeyPress(keycode.KC_D)
	if consumed {
		t.Error("unrelated key should not be consumed")
	}
	if output != keycode.KC_NO {
		t.Errorf("unrelated key should not produce output, got 0x%04X", uint16(output))
	}
}

func TestComboDetector_TimeoutReturnsKeys(t *testing.T) {
	cd := NewComboDetector(1 * time.Millisecond) // 超短いウィンドウ

	// A+B → C
	cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)

	// Aを押す
	cd.ProcessKeyPress(keycode.KC_A)

	// タイムアウトを待つ
	time.Sleep(5 * time.Millisecond)

	// Update()で保留キーが返される
	pending, count := cd.Update()
	if count != 1 {
		t.Errorf("should return 1 pending key, got %d", count)
	}
	if pending[0] != keycode.KC_A {
		t.Errorf("pending key should be KC_A, got 0x%04X", uint16(pending[0]))
	}
}

func TestComboDetector_WrongSecondKey(t *testing.T) {
	cd := NewComboDetector(500 * time.Millisecond)

	// A+B → C
	cd.RegisterCombo([]keycode.Keycode{keycode.KC_A, keycode.KC_B}, keycode.KC_C)

	// Aを押す
	cd.ProcessKeyPress(keycode.KC_A)

	// D（コンボに無関係）を押す → 候補がなくなり消費されない
	output, consumed := cd.ProcessKeyPress(keycode.KC_D)
	if output != keycode.KC_NO {
		t.Errorf("should not produce output, got 0x%04X", uint16(output))
	}
	if consumed {
		t.Error("should not consume when no candidates remain")
	}
}
