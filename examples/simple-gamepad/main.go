package main

import (
	"github.com/unagiya/keygoard/engine"
)

func main() {
	// Create keyboard instance
	keymap := engine.NewKeymap(keymapLayers)
	kb, err := engine.NewKeyboard(&boardConfig, keymap)
	if err != nil {
		panic(err)
	}

	// Start keyboard
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
