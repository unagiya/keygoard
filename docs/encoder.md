# Rotary Encoder

keygoard supports rotary encoders, allowing you to map rotation and click actions to keycodes.

English | [日本語](encoder_ja.md)

## Overview

- Up to 2 rotary encoders
- Clockwise / counter-clockwise event detection
- Click (push switch) support
- Mapping to any keycode
- High-precision rotation detection using Gray code lookup table

## Configuration

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/peripheral/encoder"
    "machine"
)

cfg := &engine.BoardConfig{
    Encoders: []*encoder.Config{
        {
            PinA:     machine.GP10,         // Encoder A pin (pull-up input)
            PinB:     machine.GP11,         // Encoder B pin (pull-up input)
            PinClick: machine.GP12,         // Click button pin
            CW:       keycode.KC_VOLU,      // Clockwise → Volume Up
            CCW:      keycode.KC_VOLD,      // Counter-clockwise → Volume Down
            Click:    keycode.KC_MUTE,      // Click → Mute
        },
        {
            PinA:     machine.GP13,
            PinB:     machine.GP14,
            PinClick: machine.NoPin,        // No click button
            CW:       keycode.KC_PGDN,
            CCW:      keycode.KC_PGUP,
        },
    },
}
```

### Config Struct

| Field | Type | Description |
|---|---|---|
| `PinA` | `machine.Pin` | Encoder A pin (pull-up input) |
| `PinB` | `machine.Pin` | Encoder B pin (pull-up input) |
| `PinClick` | `machine.Pin` | Click button pin (`machine.NoPin` if not needed) |
| `CW` | `keycode.Keycode` | Keycode for clockwise rotation |
| `CCW` | `keycode.Keycode` | Keycode for counter-clockwise rotation |
| `Click` | `keycode.Keycode` | Keycode for click (optional) |

## Events

```go
type Event uint8

const (
    EventNone  Event = iota  // No change
    EventCW                  // Clockwise
    EventCCW                 // Counter-clockwise
    EventClick               // Button click
)
```

## API

### Encoder

```go
// Constructor
enc := encoder.New(cfg)

// Update (called every 1ms, engine handles this automatically)
event := enc.Update()

// Convert event to keycode
kc := enc.GetKeycode(event)

// Cumulative rotation count
pos := enc.GetPosition()

// Reset count
enc.ResetPosition()
```

### Engine Integration

Simply pass the configuration to `BoardConfig.Encoders`, and the engine automatically:

1. Calls `Update()` every scan cycle
2. Processes the corresponding keycode as a normal key input when an event occurs
3. Includes it in the HID report

## Example

```go
package main

import (
    "machine"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/peripheral/encoder"
)

func main() {
    cfg := &engine.BoardConfig{
        // ... matrix config, etc. ...
        Encoders: []*encoder.Config{
            {
                PinA:     machine.GP10,
                PinB:     machine.GP11,
                PinClick: machine.GP12,
                CW:       keycode.KC_VOLU,
                CCW:      keycode.KC_VOLD,
                Click:    keycode.KC_MUTE,
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## Constraints

- Up to 2 encoders supported
- Pins are configured as pull-up inputs
- Rotation detected by polling at 1ms intervals (not interrupt-driven)
- Steps may be missed during very fast rotation
