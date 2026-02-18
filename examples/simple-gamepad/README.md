# Simple Gamepad Keyboard Example

This example demonstrates a simple keyboard with gamepad functionality.

## Features

- 4x6 matrix keyboard
- Analog joystick (X/Y axes)
- OLED display showing layer and joystick status
- Gaming-optimized timing (1ms scan, 3ms debounce, 1000Hz USB)

## Hardware Requirements

- RP2040 board (e.g., Raspberry Pi Pico)
- 4x6 key switch matrix
- Analog joystick module (connected to ADC0/ADC1)
- 128x32 OLED display (I2C)

## Pin Configuration

### Matrix
- Rows: GP0, GP1, GP2, GP3
- Cols: GP4, GP5, GP6, GP7, GP8, GP9

### Joystick
- X axis: GP26 (ADC0)
- Y axis: GP27 (ADC1)

### OLED
- I2C0 (default pins: GP4=SDA, GP5=SCL)
- Address: 0x3C (default)

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

## Layout

### Layer 0 (Default)
```
ESC  1    2    3    4    5
TAB  Q    W    E    R    T
CTRL A    S    D    F    GP1
SHFT Z    X    C    MO1  SPC
```

### Layer 1 (Function)
```
`    F1   F2   F3   F4   F5
---  ---  UP   ---  ---  ---
---  LEFT DOWN RGHT ---  GP2
---  ---  ---  ---  ---  ENT
```

## Customization

Edit `config.go` to change pin assignments and timing.
Edit `keymap.go` to customize the key layout.
