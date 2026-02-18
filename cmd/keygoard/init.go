package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func initProject(name string) error {
	// Create project directory
	if err := os.Mkdir(name, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create files
	files := map[string]string{
		"main.go":   mainGoTemplate,
		"config.go": configGoTemplate,
		"keymap.go": keymapGoTemplate,
		"go.mod":    goModTemplate(name),
	}

	for filename, content := range files {
		path := filepath.Join(name, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}
	}

	fmt.Printf("Created project: %s/\n", name)
	fmt.Println("  - main.go")
	fmt.Println("  - config.go")
	fmt.Println("  - keymap.go")
	fmt.Println("  - go.mod")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", name)
	fmt.Println("  # Edit config.go and keymap.go for your keyboard")
	fmt.Println("  tinygo build -target=pico -o firmware.uf2 .")
	fmt.Println("  # Flash firmware.uf2 to your RP2040")

	return nil
}

func goModTemplate(name string) string {
	return fmt.Sprintf(`module %s

go 1.21

require github.com/unagiya/keygoard v%s
`, name, Version)
}

const mainGoTemplate = `package main

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
`

const configGoTemplate = `package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/oled"
)

var boardConfig = engine.BoardConfig{
	// Matrix configuration
	// Adjust these pins for your keyboard
	Rows: []machine.Pin{
		machine.GP0, // Row 0
		machine.GP1, // Row 1
		machine.GP2, // Row 2
		machine.GP3, // Row 3
	},
	Cols: []machine.Pin{
		machine.GP4,  // Col 0
		machine.GP5,  // Col 1
		machine.GP6,  // Col 2
		machine.GP7,  // Col 3
		machine.GP8,  // Col 4
		machine.GP9,  // Col 5
	},
	MatrixType: engine.COL2ROW,

	// Timing (optional, uses defaults if omitted)
	ScanInterval:  1 * time.Millisecond,  // 1ms = 1000Hz
	DebounceTime:  3 * time.Millisecond,  // 3ms debounce
	USBReportRate: 1 * time.Millisecond,  // 1ms = 1000Hz

	// Joystick configuration (optional)
	Joystick: &joystick.Config{
		PinX:    machine.ADC0, // GP26
		PinY:    machine.ADC1, // GP27
		InvertX: false,
		InvertY: false,
	},

	// OLED configuration (optional)
	OLED: &oled.Config{
		I2C:    machine.I2C0,
		Width:  128,
		Height: 32,
	},

	// Split keyboard configuration (optional)
	// Uncomment for split keyboard
	// Split: &engine.SplitConfig{
	// 	IsMaster:  true,  // Set to false for slave side
	// 	UART:      machine.UART1,
	// 	BaudRate:  460800,
	// 	SlaveRows: 4,  // Number of rows on slave side
	// 	SlaveCols: 6,  // Number of columns on slave side
	// },
}
`

const keymapGoTemplate = `package main

import "github.com/unagiya/keygoard/keycode"

// Define your keymap here
// Structure: [16 layers][rows][cols]
var keymapLayers = [16][][]keycode.Keycode{
	// Layer 0: Default layer
	{
		// Row 0
		{keycode.KC_ESC, keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4, keycode.KC_5},
		// Row 1
		{keycode.KC_TAB, keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R, keycode.KC_T},
		// Row 2
		{keycode.KC_LCTL, keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_F, keycode.KC_G},
		// Row 3
		{keycode.KC_LSFT, keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.KC_V, keycode.KC_B},
	},

	// Layer 1: Function layer
	{
		// Row 0
		{keycode.KC_GRV, keycode.KC_F1, keycode.KC_F2, keycode.KC_F3, keycode.KC_F4, keycode.KC_F5},
		// Row 1
		{keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_UP, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS},
		// Row 2
		{keycode.KC_TRNS, keycode.KC_LEFT, keycode.KC_DOWN, keycode.KC_RGHT, keycode.KC_TRNS, keycode.KC_TRNS},
		// Row 3
		{keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS},
	},

	// Layers 2-15: Empty (nil)
	nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
}
`
