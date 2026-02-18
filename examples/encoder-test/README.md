# Encoder Test Example

This example demonstrates rotary encoder functionality (Phase 2).

## Features

- 4x4 matrix keyboard
- 2 rotary encoders with click buttons
- OLED display
- Gaming-optimized timing

## Hardware Requirements

- RP2040 board (e.g., Raspberry Pi Pico)
- 4x4 key switch matrix
- 2x rotary encoders (EC11 or similar)
- 128x32 OLED display (I2C)

## Pin Configuration

### Matrix
- Rows: GP0, GP1, GP2, GP3
- Cols: GP4, GP5, GP6, GP7

### Encoder 1 (Volume Control)
- Pin A: GP8
- Pin B: GP9
- Click: GP10
- CW: Volume Up (Keypad +)
- CCW: Volume Down (Keypad -)
- Click: Mute (Keypad Enter)

### Encoder 2 (Navigation)
- Pin A: GP11
- Pin B: GP12
- Click: GP13
- CW: Down Arrow
- CCW: Up Arrow
- Click: Enter

### OLED
- I2C0 (default pins: GP4=SDA, GP5=SCL)
- Address: 0x3C (default)

## Wiring Diagram

```
Rotary Encoder (EC11)
┌─────────┐
│   A  ───┼── GP8/GP11
│   B  ───┼── GP9/GP12
│  GND ───┼── GND
│   C  ───┼── GP10/GP13 (click)
│  C+ ───┼── 3.3V
└─────────┘
```

## Building

```bash
tinygo build -target=pico -o firmware.uf2 .
```

## Flashing

1. Hold BOOTSEL button while connecting RP2040 to USB
2. Copy firmware.uf2 to the RPI-RP2 drive:
   ```bash
   cp firmware.uf2 /Volumes/RPI-RP2/
   ```

## Usage

### Encoder 1 (Volume Control)
- Rotate CW: Volume Up
- Rotate CCW: Volume Down
- Click: Mute

### Encoder 2 (Navigation)
- Rotate CW: Down Arrow
- Rotate CCW: Up Arrow
- Click: Enter

## Customization

Edit `config.go` to change encoder pin assignments and key mappings:

```go
Encoders: []*encoder.Config{
    {
        PinA:     machine.GP8,
        PinB:     machine.GP9,
        PinClick: machine.GP10,
        CW:       keycode.KC_PPLS, // Change to your desired key
        CCW:      keycode.KC_PMNS,
        Click:    keycode.KC_PENT,
    },
}
```

## Troubleshooting

### Encoder doesn't respond
- Check wiring (A, B, GND)
- Verify pin configuration in config.go
- Test encoder with multimeter (should show changing resistance)

### Encoder skips steps
- Try adding a 0.1μF capacitor between A-GND and B-GND
- Check for loose connections

### Wrong direction
- Swap PinA and PinB in configuration

### Click button doesn't work
- Check if PinClick is connected
- Verify pull-up configuration (built-in pull-up is enabled)

## License

MIT License
