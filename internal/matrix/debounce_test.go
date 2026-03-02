package matrix

import "testing"

func TestDebouncerNoChange(t *testing.T) {
	var d Debouncer
	var raw [RowCount][ColCount]bool

	state, changed := d.Update(raw)
	if changed {
		t.Error("変化なしのとき changed は false であるべき")
	}
	if state != raw {
		t.Error("変化なしのとき state は raw と一致するべき")
	}
}

func TestDebouncerRejectsNoise(t *testing.T) {
	var d Debouncer
	var raw [RowCount][ColCount]bool

	// 1 回だけ押す → 閾値未満なので確定しない
	raw[0][0] = true
	_, changed := d.Update(raw)
	if changed {
		t.Error("1 回の読み取りで確定すべきでない")
	}

	// すぐに離す（ノイズ）
	raw[0][0] = false
	state, changed := d.Update(raw)
	if changed {
		t.Error("ノイズは状態変化とみなすべきでない")
	}
	if state[0][0] {
		t.Error("ノイズ後もキーは解放状態であるべき")
	}
}

func TestDebouncerAcceptsStablePress(t *testing.T) {
	var d Debouncer
	var raw [RowCount][ColCount]bool
	raw[0][0] = true

	// 閾値 - 1 回は確定しない
	for i := uint8(0); i < debounceThr-1; i++ {
		_, changed := d.Update(raw)
		if changed {
			t.Errorf("閾値未満（i=%d）で確定すべきでない", i)
		}
	}

	// 閾値に達したら確定
	state, changed := d.Update(raw)
	if !changed {
		t.Error("閾値到達後に確定すべき")
	}
	if !state[0][0] {
		t.Error("デバウンス後のキー状態は押下であるべき")
	}
}

func TestDebouncerAcceptsStableRelease(t *testing.T) {
	var d Debouncer
	var raw [RowCount][ColCount]bool

	// まずキーを確定させる
	raw[0][0] = true
	for i := uint8(0); i < debounceThr; i++ {
		d.Update(raw)
	}

	// 離す
	raw[0][0] = false
	for i := uint8(0); i < debounceThr-1; i++ {
		_, changed := d.Update(raw)
		if changed {
			t.Errorf("閾値未満（i=%d）で確定すべきでない", i)
		}
	}
	state, changed := d.Update(raw)
	if !changed {
		t.Error("閾値到達後に確定すべき")
	}
	if state[0][0] {
		t.Error("デバウンス後のキー状態は解放であるべき")
	}
}
