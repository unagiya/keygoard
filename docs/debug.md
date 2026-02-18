# Debug & Performance

keygoard provides debug logging and error statistics to assist with development and troubleshooting.

English | [日本語](debug_ja.md)

## Overview

- Lightweight debug logging via DebugLogger (USB CDC serial)
- Error counters and statistics via Stats struct
- Performance optimization guidelines

## Debug Mode

### Enabling

```go
cfg := &engine.BoardConfig{
    Debug: true,  // Enable debug output
    // ...
}
```

When `Debug: true` is set, the following logs are output via USB CDC serial (`println()`).

### Log Output

#### Startup Logs
- Matrix size (rows x cols)
- Enabled peripherals (joystick, OLED, encoders, LED)
- Split configuration (master/slave, baud rate)

#### Runtime Logs
- Key events (row, col, pressed/released)
- Layer changes (from → to)
- Split connection status (connected/disconnected)

### DebugLogger

A lightweight logger used internally by the engine. When `enabled=false`, all methods are no-ops with no performance impact.

```go
// Used internally by the engine (no direct instantiation needed)
logger.Log("engine", "Keyboard initialization complete")
logger.LogError("hid", err)
logger.LogKeyEvent(row, col, pressed)
logger.LogLayerChange(fromLayer, toLayer)
logger.LogSplitStatus(connected)
```

## Error Statistics

### Stats Struct

```go
type Stats struct {
    SplitRxErrors    uint16  // Split receive error count
    SplitTxErrors    uint16  // Split transmit error count
    HIDSendErrors    uint16  // HID report send error count
    LEDWriteErrors   uint16  // LED write error count
    ScanCount        uint32  // Scan count (cumulative)
    SplitDisconnects uint16  // Slave disconnect count
}
```

### Retrieving Statistics

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

// Retrieve during keyboard operation
stats := kb.GetStats()

// Check individual counters
println("Scan count:", stats.ScanCount)
println("Split Rx errors:", stats.SplitRxErrors)
println("HID send errors:", stats.HIDSendErrors)
```

## Performance Optimization Tips

The following optimizations are implemented in keygoard. Use these as guidelines for your development.

### Zero-Allocation Design

Hot paths (code executed every scan cycle) have 95% reduced heap allocations.

- **Pre-allocated buffers**: Matrix scan, debounce, LED, and split protocol buffers are allocated at startup
- **Zero-allocation APIs**: `EncodeInto()` / `DecodeMatrixStateInto()` write to existing buffers
- **Fixed-size array buffers**: TapDetector uses fixed-length arrays

### OLED Optimization

- Change detection ensures I2C transfers only occur when content changes
- Fixed byte serialization avoids heap allocation during serialization

### TinyGo Performance Guidelines

| Avoid | Recommended Alternative |
|---|---|
| `fmt.Sprintf()` | `println()` / fixed string concatenation |
| Dynamic growth via `append()` | Pre-allocated slices |
| `map` in hot paths | Arrays / fixed-size structs |
| Excessive goroutine usage | Sequential processing in main loop |
| Excessive interface boxing | Direct use of concrete types |

### Memory Targets

- Total memory usage: under 128KB
- Hot path heap allocations: target zero

## Example

```go
package main

import (
    "github.com/unagiya/keygoard/engine"
)

func main() {
    cfg := &engine.BoardConfig{
        Debug: true,  // Enable debug output
        // ... other settings ...
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // Use in OLED display callback to show stats, etc.
    // stats := kb.GetStats()

    kb.Run()
}
```

## Constraints

- Debug logs are output via USB CDC serial, requiring a serial monitor
- Slight performance impact when debug mode is enabled (`println()` overhead)
- Stats counters are `uint16` and will overflow beyond 65535
- ScanCount is `uint32` and will overflow after approximately 49 days of continuous operation (at 1000Hz)
