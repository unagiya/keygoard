# keygoard

**TinyGo-based keyboard firmware framework for custom gaming keyboards**

A modern, high-performance keyboard firmware framework built with TinyGo for RP2040/RP2350 microcontrollers. Designed specifically for gaming keyboards with support for joysticks, OLED displays, and split keyboard configurations.

## Features

- **Matrix Scanning**: COL2ROW/ROW2COL support with configurable debouncing
- **High Performance**: 1ms scan interval, 1000Hz USB report rate (gaming-optimized)
- **Layer System**: 16 layers with MO, TG, TT, LT operations
- **USB HID Composite**: 6KRO/NKRO Keyboard + 32-button Gamepad
- **Analog Joystick**: Dual-axis with oversampling, calibration, and deadzone
- **OLED Display**: SSD1306 128x32/64 support on both sides
- **Split Keyboard**: UART bidirectional communication at 460800 bps
- **Rotary Encoder**: Up to 2 encoders with click support
- **RGB LED**: WS2812 with reactive effects (KeyPress, FadeOut, Ripple, Rainbow, Breathing)
- **Macro**: Key sequence macros with delay support (up to 16 macros)
- **Combo Keys**: Simultaneous key press detection (up to 16 combos)
- **Persistent Storage**: Flash-based keymap and settings persistence
- **Easy Setup**: CLI tool for project generation

## Quick Start

### Installation

```bash
# Install the CLI tool
go install github.com/unagiya/keygoard/cmd/keygoard@latest

# Create a new keyboard project
keygoard init my-keyboard
cd my-keyboard
```

### Configuration

Edit `config.go` to match your hardware:

```go
var boardConfig = engine.BoardConfig{
    Rows: []machine.Pin{
        machine.GP0, machine.GP1, machine.GP2, machine.GP3,
    },
    Cols: []machine.Pin{
        machine.GP4, machine.GP5, machine.GP6, machine.GP7,
    },
    MatrixType: engine.COL2ROW,

    Joystick: &joystick.Config{
        PinX: machine.ADC0,
        PinY: machine.ADC1,
    },

    OLED: &oled.Config{
        I2C:    machine.I2C0,
        Width:  128,
        Height: 32,
    },
}
```

Edit `keymap.go` to define your layout:

```go
var keymapLayers = [16][][]keycode.Keycode{
    // Layer 0
    {
        {keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
        {keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.MO(1)},
    },
    // Layer 1
    {
        {keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4},
        {keycode.KC_F1, keycode.KC_F2, keycode.KC_F3, keycode.KC_TRNS},
    },
}
```

### Building & Flashing

```bash
# Build firmware
tinygo build -target=pico -o firmware.uf2 .

# Flash to RP2040
# 1. Hold BOOTSEL button while connecting to USB
# 2. Copy firmware to the device
cp firmware.uf2 /Volumes/RPI-RP2/
```

## Architecture

```
keygoard/
├── engine/              # Core keyboard engine
├── matrix/              # Matrix scanning & debouncing
├── keycode/             # Keycode definitions
├── hid/                 # USB HID reports
├── peripheral/          # Joystick, OLED, LED, Encoder support
├── split/               # Split keyboard communication
├── storage/             # Flash-based persistent storage
└── cmd/keygoard/        # CLI tool
```

## Hardware Support

### Microcontrollers
- RP2040 (Raspberry Pi Pico)
- RP2350
- Waveshare RP2040-Zero

### Matrix
- Up to 16x8 per side (128 keys)
- COL2ROW and ROW2COL support
- 3ms debounce (configurable)

### Peripherals
- **Joystick**: 2 axes (4 axes for split keyboards)
- **OLED**: SSD1306 (I2C), 128x32 or 128x64
- **Split**: UART at 460800 bps (bidirectional)
- **Rotary Encoder**: Up to 2 encoders
- **RGB LED**: WS2812, up to 128 LEDs

## Examples

See the `examples/` directory for complete working examples:

- `simple-gamepad` - Basic 4x6 keyboard with joystick and OLED
- `encoder-test` - Rotary encoder demonstration
- `led-test` - RGB LED effects showcase
- `split-bidirectional` - Full bidirectional split keyboard
- `reactive-led` - Reactive LED effects
- `split-led` - Both-side LED control
- `nkro-test` - NKRO functionality
- `macro-test` - Macro recording and playback
- `combo-test` - Combo key combinations
- `persistent-keymap` - Flash-based keymap persistence
- `waveshare-rp2040-zero` - zero-kb02 board (3×4 matrix, encoder, joystick, OLED)

For detailed keycode reference, split keyboard configuration, and performance tuning, see the [documentation](docs/README_ja.md).

## Development

### Requirements
- Go 1.21 or later
- TinyGo 0.30.0 or later
- RP2040 development board

### Building from Source
```bash
git clone https://github.com/unagiya/keygoard.git
cd keygoard
go mod download
```

### Running Examples
```bash
cd examples/simple-gamepad
tinygo build -target=pico -o firmware.uf2 .
```

## Roadmap

See [ROADMAP.md](ROADMAP.md) for detailed development history and future possibilities.

## Design Philosophy

1. **Gaming First**: Optimized for low latency and high performance
2. **Type Safe**: Leverage Go's type system for reliable firmware
3. **Memory Efficient**: Target <128KB RAM usage
4. **Extensible**: Easy to add new features and peripherals

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Acknowledgments

- Built with [TinyGo](https://tinygo.org/)
- Inspired by QMK, KMK, and ZMK firmware projects
- Uses [tinygo-drivers](https://github.com/tinygo-org/drivers) for peripherals
