# RGB LED & Effects

keygoard supports WS2812 (NeoPixel) compatible LEDs, allowing you to customize your keyboard with various lighting effects.

English | [日本語](led_ja.md)

## Overview

- WS2812 compatible LED control (up to 128 LEDs)
- 3 basic effects (Static / Breathing / Rainbow)
- 3 reactive effects (KeyPress / FadeOut / Ripple)
- Effect registry for managing up to 16 effects
- Both-side LED control for split keyboards (Unified / Independent modes)

## Configuration

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/peripheral/led"
    "machine"
)

cfg := &engine.BoardConfig{
    LED: &led.Config{
        Pin:           machine.GP16,   // WS2812 data pin
        Count:         10,             // Master-side LED count (max 128)
        SlaveLEDCount: 10,             // Slave-side LED count (Unified mode)
        KeyToLED:      keyToLEDMap,    // Key-to-LED mapping (for reactive effects)
        MaxCols:       6,              // Max columns in matrix
    },
}
```

### Config Struct

| Field | Type | Description |
|---|---|---|
| `Pin` | `machine.Pin` | WS2812 data pin |
| `Count` | `int` | Master-side LED count (max 128) |
| `SlaveLEDCount` | `int` | Slave-side LED count (used in Unified mode) |
| `KeyToLED` | `[]int8` | Key coordinate to LED index mapping (-1 = none) |
| `MaxCols` | `int` | Max columns in matrix (used for KeyToLED calculation) |

### KeyToLED Mapping

When using reactive effects, define the mapping between matrix coordinates and LED indices.

```go
// Store LED index at position row*MaxCols+col
// -1 means no mapping
keyToLED := []int8{
    0,  1,  2,  3,  -1, -1,  // Row 0: LEDs for 4 keys
    4,  5,  6,  7,  -1, -1,  // Row 1
    8,  9,  -1, -1, -1, -1,  // Row 2
}
```

## API

### Controller

```go
// Constructor
controller, err := led.New(cfg)

// Set effect
controller.SetEffect(effect)

// Update (call every frame)
controller.Update(tick)

// Individual LED control
controller.SetColor(index, color.RGBA{R: 255, G: 0, B: 0, A: 0})
controller.SetAll(color.RGBA{R: 0, G: 0, B: 255, A: 0})
controller.Clear()

// Get buffer info
buffer := controller.GetBuffer()     // []color.RGBA
count := controller.GetCount()       // Master-side LED count
total := controller.GetTotalCount()  // Master + slave total

// KeyToLED mapping
ledIdx := controller.GetLEDIndex(row, col)  // -1 = no mapping

// Reactive effect notification
controller.NotifyKeyEvent(led.KeyEvent{
    Row: 0, Col: 1, LEDIndex: 1, Pressed: true, Tick: tick,
})
```

### Via Keyboard

```go
// Set effect
kb.SetLEDEffect(led.NewRainbowEffect(5))

// Get controller
ctrl := kb.GetLEDController()
```

## Effects

### Effect Interface

```go
type Effect interface {
    Update(buffer []color.RGBA, tick uint32)
}
```

All effects implement this interface.

### Basic Effects

#### Static

Lights all LEDs with a fixed color.

```go
effect := led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0})
effect.SetColor(color.RGBA{R: 0, G: 255, B: 0})  // Change color
```

#### Breathing

Brightness oscillates using a sine wave.

```go
// speed: lower = slower (0 defaults to 10)
effect := led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255}, 10)
```

#### Rainbow

LEDs cycle through colors using HSV color space.

```go
// speed: lower = slower (0 defaults to 5)
effect := led.NewRainbowEffect(5)
```

### Reactive Effects

Effects that respond to key input. They implement the `ReactiveEffect` interface.

```go
type ReactiveEffect interface {
    Effect
    OnKeyEvent(event KeyEvent)
}

type KeyEvent struct {
    Row      uint8
    Col      uint8
    LEDIndex int8    // -1 = no mapping
    Pressed  bool
    Tick     uint32
}
```

#### KeyPress

Lights the corresponding LED while a key is pressed; turns off on release.

```go
effect := led.NewKeyPressEffect(color.RGBA{R: 255, G: 255, B: 255})
```

#### FadeOut

Lights the corresponding LED on key press, then gradually fades out.

```go
// duration: fade-out ticks (0 defaults to 60)
effect := led.NewFadeOutEffect(color.RGBA{R: 0, G: 255, B: 128}, 60)
```

#### Ripple

A ripple expands outward from the pressed key position.

```go
// posX, posY: physical coordinates for each LED
// duration: ripple duration in ticks (0 defaults to 90)
// Max 8 simultaneous ripples
posX := []uint8{0, 1, 2, 3, 0, 1, 2, 3}
posY := []uint8{0, 0, 0, 0, 1, 1, 1, 1}
effect := led.NewRippleEffect(color.RGBA{R: 128, G: 0, B: 255}, posX, posY, 90)
```

## Effect Registry

Register and switch between multiple effects. Up to 16 effects can be registered.

```go
registry := led.NewRegistry()

// Register effects
id1 := registry.Register(led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0}))
id2 := registry.Register(led.NewRainbowEffect(5))
id3 := registry.Register(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255}, 10))

// Switch effects
registry.SetCurrent(id2)       // By ID
effect := registry.Next()      // Next effect (wraps around)
effect = registry.Previous()   // Previous effect (wraps around)

// Get info
current := registry.Current()      // Current effect
currentID := registry.CurrentID()  // Current effect ID
count := registry.Count()          // Number of registered effects
```

## Split Keyboard Both-Side LED Control

Two modes are available for controlling LEDs in split keyboards.

### Unified Mode (Default)

The master manages LED buffers for both sides.

```go
cfg := &engine.BoardConfig{
    LED: &led.Config{
        Count:         10,  // Master-side LED count
        SlaveLEDCount: 10,  // Slave-side LED count
        // ...
    },
    Split: &engine.SplitConfig{
        LEDSyncMode: engine.LEDSyncUnified,  // Default
        // ...
    },
}
```

- Master computes effects for all LEDs (master + slave)
- Slave-side color data retrieved via `GetSlaveColors()` and sent over UART
- Slave applies received data via `SetColorsFromBytes()`
- Slave key events are forwarded to master, so reactive effects work on both sides

### Independent Mode

Each side runs its own effects independently.

```go
cfg := &engine.BoardConfig{
    Split: &engine.SplitConfig{
        LEDSyncMode: engine.LEDSyncIndependent,
        // ...
    },
}
```

- Each side independently runs local effects
- Slave side runs reactive effects for its own key events
- `MsgTypeLEDMode` protocol message notifies mode changes

## Example

```go
package main

import (
    "image/color"
    "machine"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/peripheral/led"
)

func main() {
    keyToLED := []int8{
        0, 1, 2, 3,
        4, 5, 6, 7,
    }

    cfg := &engine.BoardConfig{
        // ... matrix config, etc. ...
        LED: &led.Config{
            Pin:      machine.GP16,
            Count:    8,
            KeyToLED: keyToLED,
            MaxCols:  4,
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // Set effect
    kb.SetLEDEffect(led.NewRainbowEffect(5))

    kb.Run()
}
```

## Constraints

- Maximum 128 LEDs
- Effect registry limited to 16 effects
- Ripple effect supports up to 8 simultaneous ripples
- WS2812 GRB byte order supported
- In Unified mode, color data for all LEDs is transferred over UART every frame; be mindful of bandwidth with large LED counts
