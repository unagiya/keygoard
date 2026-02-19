package main

import (
	"image/color"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/peripheral/led"
)

func main() {
	cfg := newBoardConfig()
	km := newKeymap()

	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// リアクティブLEDエフェクト（キー押下で白く光る）
	kb.SetLEDEffect(led.NewKeyPressEffect(color.RGBA{R: 255, G: 255, B: 255, A: 255}))

	if err := kb.Run(); err != nil {
		panic(err)
	}
}
