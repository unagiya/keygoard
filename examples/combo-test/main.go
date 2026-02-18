package main

import "github.com/unagiya/keygoard/engine"

// コンボキーデモ
// AキーとBキーを同時押し（50ms以内）するとEscが出力される。
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
