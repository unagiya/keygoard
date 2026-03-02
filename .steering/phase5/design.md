# Phase 5: パッケージ構造リファクタリング — 設計書

## ディレクトリ構造

### Before（Phase 4）

```
keygoard/
├── keycode/
│   ├── keycode.go
│   └── keycode_test.go
├── matrix/
│   ├── const.go              ← RowCount/ColCount 定数
│   ├── debounce.go
│   ├── debounce_test.go
│   └── matrix.go             ← tinygo
├── engine/
│   ├── keyboard.go           ← tinygo（オーケストレーション + HID）
│   ├── config.go             ← tinygo（matrix.Scanner 依存）
│   ├── keymap.go             ← matrix.RowCount/ColCount 依存
│   ├── layer.go
│   ├── layer_test.go
│   ├── tap.go                ← matrix.RowCount/ColCount 依存
│   ├── tap_test.go
│   ├── peripheral.go
│   ├── errors.go
│   └── peripheral/
│       ├── encoder/
│       ├── led/
│       └── oled/
└── examples/zero-kb02/
    ├── main.go               ← 6 import + machine
    ├── config.go
    └── keymap.go
```

### After（Phase 5）

```
keygoard/
├── keycode/                   ← 公開（変更なし）
│   ├── keycode.go
│   └── keycode_test.go
├── engine/                    ← 公開（Facade）
│   ├── keyboard.go           ← tinygo（オーケストレーション + HID）
│   ├── config.go             ← tinygo（Config + ペリフェラル設定）
│   ├── types.go              ← Pin / I2CBus / Rotation 型定義
│   ├── keymap.go             ← Keymap + RowCount/ColCount 再エクスポート
│   ├── peripheral.go         ← Peripheral インターフェース（変更なし）
│   └── errors.go             ← 変更なし
├── internal/
│   ├── matrix/
│   │   ├── const.go          ← RowCount/ColCount 定数（定義元）
│   │   ├── debounce.go
│   │   ├── debounce_test.go
│   │   └── matrix.go         ← tinygo
│   ├── layer/
│   │   ├── resolver.go
│   │   └── resolver_test.go
│   ├── tap/
│   │   ├── detector.go
│   │   └── detector_test.go
│   └── peripheral/
│       ├── encoder/
│       │   ├── encoder.go    ← tinygo
│       │   ├── quadrature.go
│       │   └── quadrature_test.go
│       ├── led/
│       │   └── led.go        ← tinygo
│       └── oled/
│           ├── oled.go       ← tinygo
│           ├── font.go
│           └── font_test.go
└── examples/zero-kb02/
    ├── main.go               ← 2 import（engine + keycode）+ image/color
    └── keymap.go
```

## パッケージ依存関係

```
利用者の main.go
    │
    ├── engine（公開 Facade）
    │     ├── internal/matrix
    │     ├── internal/layer
    │     ├── internal/tap
    │     ├── internal/peripheral/encoder
    │     ├── internal/peripheral/led
    │     └── internal/peripheral/oled
    │
    └── keycode（公開 定数）
```

- `engine` → `internal/*`: OK（モジュール内）
- `internal/*` → `keycode`: OK（公開パッケージ）
- `internal/*` 同士: `tap` → `internal/matrix`（RowCount/ColCount のみ）
- 外部 → `internal/*`: コンパイルエラー

## 新規・変更ファイルの設計

### engine/types.go（新規）

ハードウェア抽象型を定義する。build tag なし。

```go
package engine

// Pin は GPIO ピン番号です。
// GPIO 番号を整数で指定します（例: GP5 → Pin(5)）。
type Pin uint8

// I2CBus は I2C バスの識別子です。
type I2CBus uint8

const (
    I2C0 I2CBus = 0
    I2C1 I2CBus = 1
)

// Rotation はディスプレイの回転角度（90° 単位）です。
type Rotation uint8

const (
    Rotation0   Rotation = 0
    Rotation90  Rotation = 1
    Rotation180 Rotation = 2
    Rotation270 Rotation = 3
)
```

### engine/config.go（変更）

マトリクスとペリフェラルの設定を統合する。

```go
//go:build tinygo

package engine

import "image/color"

type Config struct {
    ProductName string

    // マトリクス設定（旧 matrix.New() の引数を吸収）
    ColPins []Pin
    RowPins []Pin

    // キーマップ
    Keymap *Keymap

    // 標準ペリフェラル（nil の場合はスキップ）
    Encoder *EncoderConfig
    LED     *LEDConfig
    OLED    *OLEDConfig

    // カスタムペリフェラル（上級者向け拡張ポイント）
    Peripherals []Peripheral
}

type EncoderConfig struct {
    PinA   Pin
    PinB   Pin
    KeyCW  keycode.Keycode
    KeyCCW keycode.Keycode
}

// MaxLayerColors はレイヤーごとの色設定の最大数です。
const MaxLayerColors = 4

type LEDConfig struct {
    Pin         Pin
    Count       uint8
    LayerColors [MaxLayerColors]color.RGBA
}

type OLEDConfig struct {
    Bus      I2CBus
    SDA      Pin
    SCL      Pin
    Address  uint16
    Width    int16
    Height   int16
    Rotation Rotation
}
```

### engine/keymap.go（変更）

`matrix` パッケージの定数を再エクスポートする。build tag なし。

```go
package engine

import (
    "github.com/unagiya/keygoard/internal/matrix"
    "github.com/unagiya/keygoard/keycode"
)

const (
    RowCount = matrix.RowCount
    ColCount = matrix.ColCount
)

type Keymap struct {
    Layers [MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

利用者は `engine.RowCount` / `engine.ColCount` でアクセスする。

### engine/keyboard.go（変更）

`New()` 内でマトリクスとペリフェラルを内部生成する。

```go
func New(cfg *Config) *Keyboard {
    // Pin → machine.Pin 変換してマトリクス生成
    cols := pinsToCols(cfg.ColPins)
    rows := pinsToRows(cfg.RowPins)
    scanner := matrix.New(cols, rows)

    // 標準ペリフェラルを内部生成
    kb := &Keyboard{ ... }
    if cfg.Encoder != nil {
        kb.addPeripheral(encoder.New(...))
    }
    if cfg.LED != nil {
        kb.addPeripheral(led.New(...))
    }
    if cfg.OLED != nil {
        kb.addPeripheral(oled.New(...))
    }

    // カスタムペリフェラルを追加
    for _, p := range cfg.Peripherals {
        kb.addPeripheral(p)
    }

    return kb
}
```

Pin 変換ヘルパー（keyboard.go 内のプライベート関数）:

```go
func pinsToCols(pins []Pin) [matrix.ColCount]machine.Pin {
    var cols [matrix.ColCount]machine.Pin
    for i := 0; i < len(pins) && i < matrix.ColCount; i++ {
        cols[i] = machine.Pin(pins[i])
    }
    return cols
}
```

### internal/layer/resolver.go（移動）

`engine/layer.go` を移動。パッケージ名を `layer` に変更。

変更点:
- `package engine` → `package layer`
- `MaxLayers` 定数はこのパッケージで定義（engine から再エクスポート）
- `Keymap` 型への依存 → engine から引数で受け取る形に変更

```go
package layer

const MaxLayers = 4

type Resolver struct {
    keymap interface{ Layer(l, row, col int) keycode.Keycode }
    active [MaxLayers]bool
}
```

ただし、Keymap を interface にすると複雑になるため、
Resolver が直接 `[MaxLayers][RowCount][ColCount]keycode.Keycode` の配列ポインタを受け取る方がシンプル。

```go
package layer

import (
    "github.com/unagiya/keygoard/internal/matrix"
    "github.com/unagiya/keygoard/keycode"
)

const MaxLayers = 4

type Resolver struct {
    layers *[MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
    active [MaxLayers]bool
}

func NewResolver(layers *[MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode) *Resolver {
    r := &Resolver{layers: layers}
    r.active[0] = true
    return r
}
```

engine 側: `layer.NewResolver(&cfg.Keymap.Layers)`

### internal/tap/detector.go（移動）

`engine/tap.go` を移動。パッケージ名を `tap` に変更。

変更点:
- `package engine` → `package tap`
- `TapDetector` → `Detector`（パッケージ修飾で `tap.Detector` となり冗長が解消）
- `matrix.RowCount/ColCount` → `internal/matrix` から import

### internal/matrix/（移動）

現在の `matrix/` をそのまま移動。import パスのみ変更。

### internal/peripheral/*（移動）

現在の `engine/peripheral/*` を移動。変更点:
- `machine.Pin` 引数は engine が変換して渡す
- Config 型は engine 側で定義（internal 側は直接フィールドを受け取る）

## examples/zero-kb02/ の変更

### main.go（書き直し）

```go
package main

import (
    "image/color"

    _ "machine/usb/hid/keyboard"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    kb := engine.New(&engine.Config{
        ProductName: "zero-kb02",
        ColPins:     []engine.Pin{5, 6, 7, 8},
        RowPins:     []engine.Pin{9, 10, 11},
        Keymap:      defaultKeymap,
        Encoder: &engine.EncoderConfig{
            PinA: 3, PinB: 4,
            KeyCW: keycode.UpArrow, KeyCCW: keycode.DownArrow,
        },
        LED: &engine.LEDConfig{
            Pin:   1,
            Count: 12,
            LayerColors: [engine.MaxLayerColors]color.RGBA{
                {R: 0, G: 8, B: 0},
                {R: 0, G: 0, B: 8},
                {R: 8, G: 0, B: 0},
                {R: 4, G: 4, B: 4},
            },
        },
        OLED: &engine.OLEDConfig{
            Bus:      engine.I2C0,
            SDA:      12,
            SCL:      13,
            Address:  0x3C,
            Width:    128,
            Height:   64,
            Rotation: engine.Rotation180,
        },
    })
    kb.Run()
}
```

### config.go（削除）

ピン定義は main.go の Config に直接記述するため不要。

### keymap.go（変更）

```go
package main

import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

var defaultKeymap = &engine.Keymap{
    Layers: [engine.MaxLayers][engine.RowCount][engine.ColCount]keycode.Keycode{
        // Layer 0, Layer 1 ... (内容は同一)
    },
}
```

## 定数の所在と再エクスポート

定義元と公開先の対応表:

| 定数 | 定義元 | 公開先 |
|------|--------|--------|
| `RowCount` / `ColCount` | `internal/matrix` | `engine.RowCount` / `engine.ColCount` |
| `MaxLayers` | `internal/layer` | `engine.MaxLayers` |
| `MaxPeripherals` | `engine` | `engine.MaxPeripherals`（変更なし） |
| `MaxLayerColors` | `engine` | `engine.MaxLayerColors`（LED から移動） |

## Makefile の変更

テスト対象パスを更新:

```makefile
test: build
	go test ./keycode/... ./internal/...
```

## 影響を受けないもの

- `keycode/` パッケージ: 変更なし
- `Peripheral` インターフェース: 定義場所・シグネチャともに変更なし
- キーボードとしての動作: 同一
- ホットパスのメモリ特性: 固定サイズ配列のまま
