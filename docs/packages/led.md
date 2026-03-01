# led パッケージ仕様

## 概要

WS2812 互換の RGB LED 制御を提供するパッケージです。レイヤー変更時に全 LED の色を切り替えることで、現在のアクティブレイヤーを視覚的に表示します。`engine.Peripheral` インターフェースを実装しています。

## パッケージパス

`github.com/unagiya/keygoard/engine/peripheral/led`

## ファイル構成

| ファイル  | machine 依存 | 役割                              |
|-----------|:------------:|-----------------------------------|
| `led.go`  | あり         | `LED` 構造体・WS2812 制御         |

## API

### DeviceType

```go
type DeviceType uint8

const (
    WS2812 DeviceType = iota // WS2812 / SK6812MINI-E（RGB、3 バイト/LED）
    SK6812                   // SK6812（RGBW、4 バイト/LED）
)
```

### Config

```go
type Config struct {
    Pin         machine.Pin                   // データピン
    Count       uint8                         // LED 数
    Type        DeviceType                    // LED デバイス種別（デフォルト: WS2812）
    LayerColors [MaxLayerColors]color.RGBA    // レイヤーごとの LED 色
}
```

### LED

```go
leds := led.New(&led.Config{
    Pin:   machine.GPIO1,
    Count: 12,
    Type:  led.WS2812,
    LayerColors: [led.MaxLayerColors]color.RGBA{
        {G: 8},              // Layer 0: 緑
        {B: 8},              // Layer 1: 青
        {R: 8},              // Layer 2: 赤
        {R: 4, G: 4, B: 4},  // Layer 3: 白
    },
})
```

`engine.Peripheral` インターフェースを実装:

| メソッド               | 動作                                            |
|------------------------|-------------------------------------------------|
| `Init()`               | ピンを出力に設定、WS2812 デバイス初期化、Layer 0 の色で全点灯 |
| `Tick() Keycode`       | keycode.None を返す（毎サイクルの処理なし）     |
| `OnLayerChange(int)`   | 全 LED をレイヤーに対応する色に更新             |

### 定数

| 定数             | 値  | 説明                         |
|------------------|-----|------------------------------|
| `MaxLEDs`        | 12  | 制御可能な LED の最大数      |
| `MaxLayerColors` | 4   | レイヤーごとの色設定の最大数 |

## 依存パッケージ

- `tinygo.org/x/drivers/ws2812`（使用許可済み）

## テスト

LED パッケージは machine 依存のため、`make test` の TinyGo ビルドでコンパイル検証を行う。動作確認は実機で実施する。
