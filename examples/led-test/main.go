package main

import (
	"image/color"
	"machine"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/led"
)

func main() {
	// Configure board with LED
	cfg := &engine.BoardConfig{
		Rows: []machine.Pin{
			machine.GP0, machine.GP1, machine.GP2, machine.GP3,
		},
		Cols: []machine.Pin{
			machine.GP4, machine.GP5, machine.GP6,
		},
		MatrixType: engine.COL2ROW,

		// LED configuration
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: 12, // 12 LEDs
		},
	}

	// Define keymap with layer switching
	km := engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0: Static red
		{
			{keycode.KC_1, keycode.KC_2, keycode.KC_3},
			{keycode.KC_4, keycode.KC_5, keycode.KC_6},
			{keycode.KC_7, keycode.KC_8, keycode.KC_9},
			{keycode.KC_0, keycode.MO(1), keycode.MO(2)},
		},
		// Layer 1: Breathing blue
		{
			{keycode.KC_Q, keycode.KC_W, keycode.KC_E},
			{keycode.KC_R, keycode.KC_T, keycode.KC_Y},
			{keycode.KC_U, keycode.KC_I, keycode.KC_O},
			{keycode.KC_P, keycode.KC_TRNS, keycode.KC_NO},
		},
		// Layer 2: Rainbow
		{
			{keycode.KC_A, keycode.KC_S, keycode.KC_D},
			{keycode.KC_F, keycode.KC_G, keycode.KC_H},
			{keycode.KC_J, keycode.KC_K, keycode.KC_L},
			{keycode.KC_ENT, keycode.KC_NO, keycode.KC_TRNS},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// Set initial LED effect (Static red for Layer 0)
	kb.SetLEDEffect(led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	// Note: Layer-based LED effect switching will be implemented in Phase 2 M3
	// For now, you can manually switch effects by calling:
	// kb.SetLEDEffect(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255, A: 255}, 10))
	// kb.SetLEDEffect(led.NewRainbowEffect(5))

	// Run keyboard
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
