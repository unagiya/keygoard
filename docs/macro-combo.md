# Macros & Combo Keys

keygoard supports key sequence macros and combo keys triggered by pressing multiple keys simultaneously.

English | [日本語](macro-combo_ja.md)

## Overview

### Macros
- Up to 16 macros
- Each macro supports up to 32 steps
- Step types: key down / key up / key tap / delay
- Declarative definition via `BoardConfig` and dynamic registration via `RegisterMacro()` API

### Combo Keys
- Up to 16 combos
- Each combo uses 2-4 simultaneous keys
- Configurable detection window (default 50ms)
- Declarative definition via `BoardConfig` and dynamic registration via `RegisterCombo()` API

## Macros

### Macro Steps

```go
type MacroStepType uint8

const (
    MacroStepEnd     MacroStepType = 0x00  // End marker
    MacroStepKeyDown MacroStepType = 0x01  // Key press
    MacroStepKeyUp   MacroStepType = 0x02  // Key release
    MacroStepKeyTap  MacroStepType = 0x03  // Press + immediate release
    MacroStepDelay   MacroStepType = 0x04  // Delay (Value = milliseconds)
)

type MacroStep struct {
    Type  MacroStepType
    Value uint16  // Keycode value or delay in milliseconds
}
```

### BoardConfig Definition

```go
cfg := &engine.BoardConfig{
    Macros: []engine.MacroDef{
        {
            ID: 0,  // Macro ID (0-15)
            Steps: []engine.MacroStep{
                // Send Ctrl+C
                {Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
                {Type: engine.MacroStepKeyTap,  Value: uint16(keycode.KC_C)},
                {Type: engine.MacroStepKeyUp,   Value: uint16(keycode.KC_LCTL)},
            },
        },
        {
            ID: 1,
            Steps: []engine.MacroStep{
                // Type "gg" with 100ms delay
                {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_G)},
                {Type: engine.MacroStepDelay,  Value: 100},
                {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_G)},
            },
        },
    },
}
```

### Keymap Assignment

Macros are triggered by keycodes `KC_MACRO0` through `KC_MACRO15` (0x7010-0x701F).

```go
keymap := &engine.Keymap{
    Layers: [16][][]keycode.Keycode{
        { // Layer 0
            {keycode.KC_MACRO0, keycode.KC_MACRO1, keycode.KC_A, keycode.KC_B},
        },
    },
}
```

### Dynamic Registration API

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

err := kb.RegisterMacro(2, []engine.MacroStep{
    {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_H)},
    {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_I)},
})
```

### MacroManager API

```go
mm := engine.NewMacroManager()

// Register macro
mm.RegisterMacro(id, steps)

// Trigger macro
mm.Trigger(id)

// Check execution state
running := mm.IsRunning()

// Called every scan (engine handles this automatically)
kc, stepType, active := mm.Update()
```

## Combo Keys

### BoardConfig Definition

```go
cfg := &engine.BoardConfig{
    Combos: []engine.ComboDef{
        {
            Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_A, keycode.KC_B},
            Count:  2,              // Number of keys
            Output: keycode.KC_ESC, // Keycode emitted on combo activation
        },
        {
            Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_J, keycode.KC_K, keycode.KC_L},
            Count:  3,
            Output: keycode.KC_ENT,
        },
    },
    ComboWindow: 50 * time.Millisecond,  // Detection window (default 50ms)
}
```

### ComboDef Struct

| Field | Type | Description |
|---|---|---|
| `Keys` | `[4]keycode.Keycode` | Component keys (terminated by `KC_NO`) |
| `Count` | `uint8` | Number of component keys (2-4) |
| `Output` | `keycode.Keycode` | Keycode emitted when combo activates |

### Dynamic Registration API

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

err := kb.RegisterCombo(
    []keycode.Keycode{keycode.KC_D, keycode.KC_F},
    keycode.KC_TAB,
)
```

### ComboDetector API

```go
cd := engine.NewComboDetector(50 * time.Millisecond)

// Register combo
cd.RegisterCombo(keys, output)

// Process key press
outputKC, consumed := cd.ProcessKeyPress(kc)
// consumed=true → the original key is held for combo detection

// Process key release
cd.ProcessKeyRelease(kc)

// Timeout processing (called every scan)
pendingKeys, count := cd.Update()
// Returns held keys when the window times out
```

### How Combos Work

1. When a key that is part of a combo is pressed, it is temporarily held
2. If all component keys are pressed within the detection window, the combo activates and the `Output` keycode is emitted
3. If the window times out, held keys are processed as normal key inputs

## Example

```go
package main

import (
    "time"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    cfg := &engine.BoardConfig{
        // ... matrix config, etc. ...
        Macros: []engine.MacroDef{
            {
                ID: 0,
                Steps: []engine.MacroStep{
                    // Ctrl+Z (undo)
                    {Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
                    {Type: engine.MacroStepKeyTap,  Value: uint16(keycode.KC_Z)},
                    {Type: engine.MacroStepKeyUp,   Value: uint16(keycode.KC_LCTL)},
                },
            },
        },
        Combos: []engine.ComboDef{
            {
                Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_A, keycode.KC_S},
                Count:  2,
                Output: keycode.KC_ESC,
            },
        },
        ComboWindow: 50 * time.Millisecond,
    }

    keymap := &engine.Keymap{
        Layers: [16][][]keycode.Keycode{
            {
                {keycode.KC_MACRO0, keycode.KC_A, keycode.KC_S, keycode.KC_D},
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## Constraints

### Macros
- Maximum 16 macros
- Maximum 32 steps per macro
- Only one macro can run at a time
- During execution, one step is processed per scan cycle

### Combo Keys
- Maximum 16 combos
- Each combo consists of 2-4 keys
- Default detection window is 50ms
- Slight input latency while combo keys are held for detection
