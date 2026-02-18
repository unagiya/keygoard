# RGB LED・エフェクト

keygoardはWS2812（NeoPixel）互換LEDをサポートし、多彩なエフェクトでキーボードをカスタマイズできます。

[English](led.md) | 日本語

## 概要

- WS2812互換LEDの制御（最大128個）
- 基本エフェクト3種（Static / Breathing / Rainbow）
- リアクティブエフェクト3種（KeyPress / FadeOut / Ripple）
- エフェクトレジストリによる切り替え管理（最大16エフェクト）
- 分割キーボードでの両側LED制御（Unified / Independent モード）

## 設定

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/peripheral/led"
    "machine"
)

cfg := &engine.BoardConfig{
    LED: &led.Config{
        Pin:           machine.GP16,   // WS2812 データピン
        Count:         10,             // マスター側LED数（最大128）
        SlaveLEDCount: 10,             // スレーブ側LED数（Unifiedモード時）
        KeyToLED:      keyToLEDMap,    // キー→LEDマッピング（リアクティブ用）
        MaxCols:       6,              // マトリクスの最大列数
    },
}
```

### Config構造体

| フィールド | 型 | 説明 |
|---|---|---|
| `Pin` | `machine.Pin` | WS2812データピン |
| `Count` | `int` | マスター側LED数（最大128） |
| `SlaveLEDCount` | `int` | スレーブ側LED数（Unifiedモード時に使用） |
| `KeyToLED` | `[]int8` | キー座標→LEDインデックスのマッピング（-1=なし） |
| `MaxCols` | `int` | マトリクスの最大列数（KeyToLEDの計算に使用） |

### KeyToLEDマッピング

リアクティブエフェクトを使用する場合、キーのマトリクス座標とLEDインデックスの対応を定義します。

```go
// row*MaxCols+col の位置にLEDインデックスを格納
// -1 はマッピングなし
keyToLED := []int8{
    0,  1,  2,  3,  -1, -1,  // Row 0: 4キー分のLED
    4,  5,  6,  7,  -1, -1,  // Row 1
    8,  9,  -1, -1, -1, -1,  // Row 2
}
```

## API

### Controller

```go
// コンストラクタ
controller, err := led.New(cfg)

// エフェクト設定
controller.SetEffect(effect)

// 更新（毎フレーム呼び出し）
controller.Update(tick)

// 個別LED操作
controller.SetColor(index, color.RGBA{R: 255, G: 0, B: 0, A: 0})
controller.SetAll(color.RGBA{R: 0, G: 0, B: 255, A: 0})
controller.Clear()

// バッファ取得
buffer := controller.GetBuffer()     // []color.RGBA
count := controller.GetCount()       // マスター側LED数
total := controller.GetTotalCount()  // マスター + スレーブ合計

// KeyToLEDマッピング
ledIdx := controller.GetLEDIndex(row, col)  // -1 = マッピングなし

// リアクティブエフェクト通知
controller.NotifyKeyEvent(led.KeyEvent{
    Row: 0, Col: 1, LEDIndex: 1, Pressed: true, Tick: tick,
})
```

### Keyboard経由の操作

```go
// エフェクト設定
kb.SetLEDEffect(led.NewRainbowEffect(5))

// コントローラ取得
ctrl := kb.GetLEDController()
```

## エフェクト

### Effectインターフェース

```go
type Effect interface {
    Update(buffer []color.RGBA, tick uint32)
}
```

すべてのエフェクトはこのインターフェースを実装します。

### 基本エフェクト

#### Static（単色）

全LEDを指定した色で点灯します。

```go
effect := led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0})
effect.SetColor(color.RGBA{R: 0, G: 255, B: 0})  // 色変更
```

#### Breathing（呼吸）

サイン波で明るさが変化します。

```go
// speed: 値が小さいほどゆっくり（0でデフォルト10）
effect := led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255}, 10)
```

#### Rainbow（レインボー）

HSV色空間を使い、LEDが虹色に変化します。

```go
// speed: 値が小さいほどゆっくり（0でデフォルト5）
effect := led.NewRainbowEffect(5)
```

### リアクティブエフェクト

キー入力に反応するエフェクトです。`ReactiveEffect`インターフェースを実装します。

```go
type ReactiveEffect interface {
    Effect
    OnKeyEvent(event KeyEvent)
}

type KeyEvent struct {
    Row      uint8
    Col      uint8
    LEDIndex int8    // -1 = マッピングなし
    Pressed  bool
    Tick     uint32
}
```

#### KeyPress（キー押下点灯）

キーを押している間だけ対応するLEDが点灯し、離すと消灯します。

```go
effect := led.NewKeyPressEffect(color.RGBA{R: 255, G: 255, B: 255})
```

#### FadeOut（フェードアウト）

キーを押すと対応するLEDが点灯し、時間経過で徐々に消灯します。

```go
// duration: フェードアウトtick数（0でデフォルト60）
effect := led.NewFadeOutEffect(color.RGBA{R: 0, G: 255, B: 128}, 60)
```

#### Ripple（波紋）

キーを押した位置から波紋が広がるように周囲のLEDが点灯します。

```go
// posX, posY: 各LEDの物理座標
// duration: 波紋の持続tick数（0でデフォルト90）
// 同時最大8つの波紋
posX := []uint8{0, 1, 2, 3, 0, 1, 2, 3}
posY := []uint8{0, 0, 0, 0, 1, 1, 1, 1}
effect := led.NewRippleEffect(color.RGBA{R: 128, G: 0, B: 255}, posX, posY, 90)
```

## エフェクトレジストリ

複数のエフェクトを登録し、切り替えて使用できます。最大16エフェクトまで登録可能です。

```go
registry := led.NewRegistry()

// エフェクト登録
id1 := registry.Register(led.NewStaticEffect(color.RGBA{R: 255, G: 0, B: 0}))
id2 := registry.Register(led.NewRainbowEffect(5))
id3 := registry.Register(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255}, 10))

// 切り替え
registry.SetCurrent(id2)       // IDで指定
effect := registry.Next()      // 次のエフェクトへ（循環）
effect = registry.Previous()   // 前のエフェクトへ（循環）

// 情報取得
current := registry.Current()      // 現在のエフェクト
currentID := registry.CurrentID()  // 現在のエフェクトID
count := registry.Count()          // 登録数
```

## 分割キーボード両側LED制御

分割キーボードでは2つのモードでLEDを制御できます。

### Unifiedモード（デフォルト）

マスターが両側のLEDバッファを統合管理します。

```go
cfg := &engine.BoardConfig{
    LED: &led.Config{
        Count:         10,  // マスター側LED数
        SlaveLEDCount: 10,  // スレーブ側LED数
        // ...
    },
    Split: &engine.SplitConfig{
        LEDSyncMode: engine.LEDSyncUnified,  // デフォルト
        // ...
    },
}
```

- マスターが全LED（マスター + スレーブ）のエフェクトを計算
- `GetSlaveColors()`でスレーブ側のカラーデータを取得し、UART経由で転送
- スレーブは`SetColorsFromBytes()`で受信データを反映
- スレーブ側のキーイベントもマスターに転送され、リアクティブエフェクトが両側で動作

### Independentモード

各側が独立してエフェクトを実行します。

```go
cfg := &engine.BoardConfig{
    Split: &engine.SplitConfig{
        LEDSyncMode: engine.LEDSyncIndependent,
        // ...
    },
}
```

- 各側が独立してローカルのエフェクトを動作
- スレーブ側も自身のキーイベントに対してリアクティブエフェクトを実行
- `MsgTypeLEDMode`プロトコルメッセージでモード通知

## 使用例

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
        // ... マトリクス設定等 ...
        LED: &led.Config{
            Pin:      machine.GP16,
            Count:    8,
            KeyToLED: keyToLED,
            MaxCols:  4,
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // エフェクト設定
    kb.SetLEDEffect(led.NewRainbowEffect(5))

    kb.Run()
}
```

## 制約・制限

- LED数は最大128個
- エフェクトレジストリの登録上限は16個
- Rippleエフェクトの同時波紋数は最大8
- WS2812のGRB順序に対応
- Unifiedモードでは全LEDのカラーデータがUART経由で毎フレーム転送されるため、LED数が多いと帯域に注意
