# waveshare-rp2040-zero Example

Example firmware for the **zero-kb02** board by [sago35/keyboards](https://github.com/sago35/keyboards).

[日本語 README](README_ja.md)

## Target Hardware

- Waveshare RP2040-Zero (MCU)
- 3-row × 4-column key matrix (12 keys)
- WS2812 RGB LEDs × 12 (one per key)
- Rotary encoder × 1 (with click button)
- Analog joystick × 1
- SSD1306 OLED (128×32)

## Pin Assignment

| Function | Pin |
|----------|-----|
| Matrix rows | GPIO9, GPIO10, GPIO11 |
| Matrix cols | GPIO5, GPIO6, GPIO7, GPIO8 |
| WS2812 LED | GPIO1 |
| Encoder A | GPIO3 |
| Encoder B | GPIO4 |
| Encoder button | GPIO2 |
| Joystick X | GPIO29 (ADC3) |
| Joystick Y | GPIO28 (ADC2) |
| OLED SDA | GPIO12 (I2C0) |
| OLED SCL | GPIO13 (I2C0) |

## Key Layout

### Layer 0 (default)

```
┌───┬───┬───┬───┐
│ Q │ W │ E │ R │
├───┼───┼───┼───┤
│ A │ S │ D │ F │
├───┼───┼───┼───┤
│ Z │ X │ C │MO1│
└───┴───┴───┴───┘
```

### Layer 1 (hold MO1)

```
┌─────┬────┬────┬────┐
│ ESC │ F1 │ F2 │ F3 │
├─────┼────┼────┼────┤
│ TAB │ F4 │ F5 │ F6 │
├─────┼────┼────┼────┤
│     │ F7 │ F8 │    │
└─────┴────┴────┴────┘
```

### Encoder

| Action | Keycode |
|--------|---------|
| Clockwise | KC_PGUP (page up) |
| Counter-clockwise | KC_PGDN (page down) |
| Click | KC_ENT (enter) |

## LED Effects

The default effect is **KeyPressEffect** (white glow on key press).
Edit `SetLEDEffect` in `main.go` to switch effects:

```go
// Breathing effect (blue)
kb.SetLEDEffect(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255, A: 255}, 10))

// Rainbow effect
kb.SetLEDEffect(led.NewRainbowEffect(5))
```

## Build & Flash

```bash
# Build firmware
tinygo build -target=waveshare-rp2040-zero -o firmware-waveshare-zero.uf2 .

# Flash to RP2040
# 1. Hold BOOTSEL button while connecting to USB
# 2. Copy firmware to the device
cp firmware-waveshare-zero.uf2 /Volumes/RPI-RP2/
```

## Customization

Edit `config.go` to change pin assignments or peripheral settings.
Edit `keymap.go` to modify the key layout.
