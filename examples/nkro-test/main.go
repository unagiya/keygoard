package main

import "github.com/unagiya/keygoard/engine"

// NKRO切り替えデモ
// キーマップにKC_NKRO_TOGGLE / KC_NKRO_ON / KC_NKRO_OFFを配置し、
// 6KROとNKROの切り替えを試すサンプル。
func main() {
	cfg := newBoardConfig()
	km := newKeymap()

	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	if err := kb.Run(); err != nil {
		panic(err)
	}
}
