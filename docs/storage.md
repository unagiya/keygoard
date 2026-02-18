# Persistent Storage

keygoard supports saving settings to Flash memory, allowing keymaps, configuration, and macro data to persist across power cycles.

English | [日本語](storage_ja.md)

## Overview

- Save, load, and reset data in Flash memory
- Keymap storage (up to 8KB)
- Settings storage (HID mode, LED effect, brightness)
- Macro data storage (up to 1.5KB)
- Data integrity verification via XOR rotation checksum

## Configuration

### BoardConfig

```go
cfg := &engine.BoardConfig{
    Storage: true,  // Enable Flash persistence
    // ...
}
```

### Initializing Storage

```go
import (
    "github.com/unagiya/keygoard/storage"
    "machine"
)

// Initialize with Flash device and offset
store, err := storage.New(machine.Flash, 0x100000)  // Adjust offset for your environment

// Set storage on keyboard
kb.SetStorage(store)

// Load settings at startup
kb.LoadSettings()
```

## Storage Layout

| Region | Offset | Size | Description |
|---|---|---|---|
| Header | 0 | 16B | Magic, version, checksum |
| Keymap | 16 | Up to 8,192B | Keycode data for layers x rows x cols |
| Config | 8,208 | Up to 64B | HID mode, LED settings, etc. |
| Macros | 8,272 | Up to 1,536B | 16 macros x 32 steps x 3 bytes |
| **Total** | | **9,808B** | **Required sectors: 3 (12KB)** |

### StorageHeader

```go
type StorageHeader struct {
    Magic    [4]byte  // "KBD\x00"
    Version  uint8    // 1
    Reserved [3]byte
    DataSize uint32   // Data size excluding header
    Checksum uint32   // XOR rotation checksum
}
```

## API

### High-Level API (via Keyboard)

```go
// Save settings (keymap + HID mode)
err := kb.SaveSettings()

// Load settings (restore HID mode)
err := kb.LoadSettings()

// Reset settings (overwrite with empty data)
err := kb.ResetSettings()
```

### Low-Level API (Storage Package)

```go
// Initialize
store, err := storage.New(flash, offset)

// Read/write individual data
store.Read(buf, offset)
store.Write(data, offset)

// Header operations
header, err := store.ReadHeader()
store.WriteHeader(header)

// Validity check (magic + checksum verification)
valid := store.IsValid()

// Bulk save/load
store.SaveAll(keymapData, configData, macroData)
keymapData, configData, macroData, err := store.LoadAll()

// Checksum calculation
checksum := storage.CalculateChecksum(data)
```

### ConfigData

```go
type ConfigData struct {
    HIDMode       uint8    // 0=6KRO, 1=NKRO
    LEDEffectID   uint8    // Current LED effect ID
    LEDBrightness uint8    // LED brightness 0-255
    Reserved      [61]byte
}

// Serialize / Deserialize
data := storage.SerializeConfig(cfg)
cfg, err := storage.DeserializeConfig(data)
```

### Keymap Serialization

```go
// Format: [layers(1)][rows(1)][cols(1)][reserved(1)][keycode(2B) x L x R x C]
data := storage.SerializeKeymap(layers, numLayers, rows, cols)
layers, numLayers, rows, cols, err := storage.DeserializeKeymap(data)
```

### Macro Serialization

```go
// Each step: Type(1B) + Value(2B) = 3 bytes
data := storage.SerializeMacros(macros)
macros, err := storage.DeserializeMacros(data)
```

### BlockDevice Interface

```go
type BlockDevice interface {
    ReadAt(buf []byte, off int64) (int, error)
    WriteAt(p []byte, off int64) (int, error)
    EraseBlocks(start, length int64) error
}
```

`machine.Flash` satisfies this interface.

## Example

```go
package main

import (
    "machine"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/storage"
)

func main() {
    cfg := &engine.BoardConfig{
        Storage: true,
        // ... other settings ...
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // Initialize and set storage
    store, _ := storage.New(machine.Flash, 0x100000)
    kb.SetStorage(store)

    // Load saved settings
    if err := kb.LoadSettings(); err != nil {
        // Error on first boot or data corruption (normal behavior)
    }

    // Run keyboard
    kb.Run()

    // Save after settings change (e.g., after NKRO mode toggle)
    // kb.SaveSettings()
}
```

## Errors

| Error | Description |
|---|---|
| `ErrInvalidMagic` | Invalid magic bytes in header |
| `ErrInvalidChecksum` | Checksum mismatch |
| `ErrDataTooLarge` | Data exceeds region size |
| `ErrNotInitialized` | Storage not initialized |

## Constraints

- Flash writes occur at sector granularity (4KB), requiring sector erasure before writing
- 3 sectors required (12KB)
- Frequent writes affect Flash lifespan; save only when settings change
- Checksum uses XOR rotation (not CRC)
