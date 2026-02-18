package main

import (
	"image/color"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/peripheral/led"
)

func main() {
	cfg := newBoardConfig()
	km := newKeymap()

	// キーボードを作成
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// エフェクトレジストリを作成し、各種エフェクトを登録
	registry := led.NewRegistry()

	// 0: 静的エフェクト（赤）
	registry.Register(led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	// 1: キー押下エフェクト（シアン）
	registry.Register(led.NewKeyPressEffect(color.RGBA{R: 0, G: 255, B: 255, A: 255}))

	// 2: フェードアウトエフェクト（マゼンタ、90tick持続）
	registry.Register(led.NewFadeOutEffect(color.RGBA{R: 255, G: 0, B: 255, A: 255}, 90))

	// 3: 波紋エフェクト（白、90tick持続）
	registry.Register(led.NewRippleEffect(
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		ledPositionsX,
		ledPositionsY,
		90,
	))

	// 4: レインボーエフェクト
	registry.Register(led.NewRainbowEffect(5))

	// 5: ブリージングエフェクト（青）
	registry.Register(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255, A: 255}, 10))

	// 初期エフェクトを設定（キー押下エフェクト）
	registry.SetCurrent(1)
	kb.SetLEDEffect(registry.Current())

	// キーボードを実行
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
