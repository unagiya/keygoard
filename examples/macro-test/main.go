package main

import "github.com/unagiya/keygoard/engine"

// マクロデモ
// KC_MACRO0キーを押すとCtrl+C → 100ms遅延 → Ctrl+Vの
// コピー＆ペーストマクロが実行される。
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
