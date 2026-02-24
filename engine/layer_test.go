package engine

import (
	"testing"

	"github.com/unagiya/keygoard/keycode"
)

func TestResolveSingleLayer(t *testing.T) {
	km := &Keymap{}
	km.Layers[0] = testLayer(keycode.A, keycode.B, keycode.C, keycode.D)
	r := NewResolver(km)

	if got := r.Resolve(0, 0); got != keycode.A {
		t.Errorf("Resolve(0,0) = %#x, want %#x", got, keycode.A)
	}
	if got := r.Resolve(0, 1); got != keycode.B {
		t.Errorf("Resolve(0,1) = %#x, want %#x", got, keycode.B)
	}
}

func TestResolveHigherLayerPriority(t *testing.T) {
	km := &Keymap{}
	km.Layers[0] = testLayer(keycode.A, keycode.B, keycode.C, keycode.D)
	km.Layers[1] = testLayer(keycode.Num1, keycode.Num2, keycode.Num3, keycode.Num4)
	r := NewResolver(km)
	r.Activate(1)

	// レイヤー 1 が優先される
	if got := r.Resolve(0, 0); got != keycode.Num1 {
		t.Errorf("Resolve(0,0) = %#x, want %#x", got, keycode.Num1)
	}
}

func TestResolveTRNSFallthrough(t *testing.T) {
	km := &Keymap{}
	km.Layers[0] = testLayer(keycode.A, keycode.B, keycode.C, keycode.D)
	// レイヤー 1: (0,0) は TRNS、(0,1) は Num2
	km.Layers[1][0][0] = keycode.TRNS
	km.Layers[1][0][1] = keycode.Num2
	r := NewResolver(km)
	r.Activate(1)

	// TRNS → レイヤー 0 にフォールバック
	if got := r.Resolve(0, 0); got != keycode.A {
		t.Errorf("Resolve(0,0) = %#x, want %#x (TRNS fallthrough)", got, keycode.A)
	}
	// 非 TRNS → レイヤー 1 のまま
	if got := r.Resolve(0, 1); got != keycode.Num2 {
		t.Errorf("Resolve(0,1) = %#x, want %#x", got, keycode.Num2)
	}
}

func TestResolveAllTRNSReturnsNone(t *testing.T) {
	km := &Keymap{}
	// レイヤー 0 も TRNS
	km.Layers[0][0][0] = keycode.TRNS
	km.Layers[1][0][0] = keycode.TRNS
	r := NewResolver(km)
	r.Activate(1)

	if got := r.Resolve(0, 0); got != keycode.None {
		t.Errorf("Resolve(0,0) = %#x, want None (%#x)", got, keycode.None)
	}
}

func TestResolveNoneOnLayer0(t *testing.T) {
	km := &Keymap{}
	// レイヤー 0 が None（ゼロ値）の場合はそのまま None を返す
	r := NewResolver(km)

	if got := r.Resolve(0, 0); got != keycode.None {
		t.Errorf("Resolve(0,0) = %#x, want None", got)
	}
}

func TestActivateDeactivate(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	if r.IsActive(1) {
		t.Error("layer 1 should be inactive initially")
	}

	r.Activate(1)
	if !r.IsActive(1) {
		t.Error("layer 1 should be active after Activate")
	}

	r.Deactivate(1)
	if r.IsActive(1) {
		t.Error("layer 1 should be inactive after Deactivate")
	}
}

func TestToggle(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	r.Toggle(2)
	if !r.IsActive(2) {
		t.Error("layer 2 should be active after first Toggle")
	}

	r.Toggle(2)
	if r.IsActive(2) {
		t.Error("layer 2 should be inactive after second Toggle")
	}
}

func TestLayer0CannotBeDeactivated(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	r.Deactivate(0)
	if !r.IsActive(0) {
		t.Error("layer 0 should always be active")
	}

	r.Toggle(0)
	if !r.IsActive(0) {
		t.Error("layer 0 should always be active even after Toggle")
	}
}

func TestLayer0AlwaysActive(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	if !r.IsActive(0) {
		t.Error("layer 0 should be active on init")
	}
}

func TestIsActiveOutOfRange(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	if r.IsActive(-1) {
		t.Error("IsActive(-1) should return false")
	}
	if r.IsActive(MaxLayers) {
		t.Errorf("IsActive(%d) should return false", MaxLayers)
	}
}

func TestActivateOutOfRange(t *testing.T) {
	km := &Keymap{}
	r := NewResolver(km)

	// パニックしないことを確認
	r.Activate(-1)
	r.Activate(MaxLayers)
	r.Deactivate(-1)
	r.Deactivate(MaxLayers)
	r.Toggle(-1)
	r.Toggle(MaxLayers)
}

// testLayer はテスト用に最初の 4 キー（row=0, col=0..3）を設定したレイヤーを返します。
func testLayer(k0, k1, k2, k3 keycode.Keycode) [3][4]keycode.Keycode {
	var layer [3][4]keycode.Keycode
	layer[0][0] = k0
	layer[0][1] = k1
	layer[0][2] = k2
	layer[0][3] = k3
	return layer
}
