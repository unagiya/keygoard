# NKRO (N-Key Rollover)

keygoard supports dynamic switching between 6KRO (6-Key Rollover) and NKRO (N-Key Rollover) modes.

English | [日本語](nkro_ja.md)

## Overview

- Dynamic switching between 6KRO and NKRO
- Bitmap-based NKRO report (up to 120 simultaneous keys)
- Mode switching via keycodes (toggle / on / off)
- Unified keyboard + gamepad management via CompositeHID

## Configuration

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/hid"
)

cfg := &engine.BoardConfig{
    DefaultHIDMode: hid.HIDModeNKRO,  // Start in NKRO mode
    // ...
}
```

### HIDMode

```go
type HIDMode uint8

const (
    HIDMode6KRO HIDMode = 0  // 6KRO (default)
    HIDModeNKRO HIDMode = 1  // NKRO
)
```

## Switching via Keycodes

Place these keycodes in your keymap to allow users to switch modes.

| Keycode | Value | Description |
|---|---|---|
| `KC_NKRO_TOGGLE` | `0x7000` | Toggle NKRO (6KRO ↔ NKRO) |
| `KC_NKRO_ON` | `0x7001` | Enable NKRO |
| `KC_NKRO_OFF` | `0x7002` | Disable NKRO |

```go
keymap := &engine.Keymap{
    Layers: [16][][]keycode.Keycode{
        { // Layer 0
            {keycode.KC_A, keycode.KC_B, keycode.KC_C, keycode.KC_NKRO_TOGGLE},
        },
    },
}
```

## API

### UnifiedKeyboardHID

A wrapper that transparently switches between 6KRO and NKRO.

```go
uhid := hid.NewUnifiedKeyboardHID(hid.HIDMode6KRO)

// Key operations (uses appropriate report based on current mode)
uhid.AddKey(kc)
uhid.SetModifier(mask)
uhid.SendReport()
uhid.Clear()

// Mode switching
uhid.SetMode(hid.HIDModeNKRO)
uhid.ToggleMode()

// State queries
mode := uhid.GetMode()
report := uhid.GetReport()         // 6KRO report
nkroReport := uhid.GetNKROReport() // NKRO report
```

### NKROReport

```go
type NKROReport struct {
    Modifier uint8     // Modifier byte
    Reserved uint8
    Keys     [15]byte  // Bitmap (120 keys, 0x04-0x7B)
}

report := &hid.NKROReport{}
report.Clear()
report.AddKey(0x04)     // Add key (0x04-0x7B only)
report.RemoveKey(0x04)  // Remove key
report.HasKey(0x04)     // Check if key is present
data := report.ToBytes()  // 17-byte report data
```

### CompositeHID

Manages keyboard and gamepad HID reports together.

```go
chid := hid.NewCompositeHID()                          // 6KRO default
chid := hid.NewCompositeHIDWithMode(hid.HIDModeNKRO)   // Start with NKRO

// Access each device
keyboard := chid.Keyboard()  // *UnifiedKeyboardHID
gamepad := chid.Gamepad()    // *GamepadHID

// Send reports
chid.SendReports()
chid.Clear()
```

## 6KRO vs NKRO

| Feature | 6KRO | NKRO |
|---|---|---|
| Simultaneous keys | Up to 6 + modifiers | Up to 120 + modifiers |
| Report size | 8 bytes | 17 bytes |
| Compatibility | Works with all OS | May not work in some BIOS |
| Default | Yes | No |

## Example

```go
package main

import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/hid"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    cfg := &engine.BoardConfig{
        DefaultHIDMode: hid.HIDModeNKRO,
        // ... other settings ...
    }

    keymap := &engine.Keymap{
        Layers: [16][][]keycode.Keycode{
            { // Layer 0: normal keys + NKRO toggle
                {keycode.KC_A, keycode.KC_B, keycode.KC_C, keycode.KC_NKRO_TOGGLE},
            },
            { // Layer 1: NKRO ON/OFF
                {keycode.KC_NKRO_ON, keycode.KC_NKRO_OFF, keycode.KC_TRNS, keycode.KC_TRNS},
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## Constraints

- Due to TinyGo USB stack limitations, even in NKRO mode, USB reports are sent in 6KRO-compatible format (up to 6 keys extracted from bitmap). Full NKRO will be available after TinyGo supports custom HID descriptors
- NKRO report covers key range 0x04-0x7B (120 keys)
- Some BIOS and bootloaders may not recognize NKRO reports
- When Storage is enabled, HID mode state is persisted
