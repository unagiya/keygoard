# led パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/peripheral/led`（非公開: 外部 import 不可）

## 概要

WS2812 互換の RGB LED 制御を提供するパッケージです。レイヤー変更時に全 LED の色を切り替えることで、現在のアクティブレイヤーを視覚的に表示します。`engine.Peripheral` インターフェースを実装しています。

利用者は `engine.LEDConfig` 経由で設定し、`engine.New()` が内部でこのパッケージを呼び出します。

## ファイル構成

| ファイル  | machine 依存 | 役割                              |
|-----------|:------------:|-----------------------------------|
| `led.go`  | あり         | `LED` 構造体・WS2812/SK6812 制御  |

## API

### DeviceType

```go
type DeviceType uint8

const (
    TypeWS2812 DeviceType = iota // WS2812 / SK6812MINI-E（RGB、3 バイト/LED）
    TypeSK6812                   // SK6812（RGBW、4 バイト/LED）
)
```

利用者は内部の `DeviceType` を直接扱いません。`engine.LEDType`（`engine.WS2812` / `engine.SK6812`）を `LEDConfig.Type` に指定すると、engine が内部で変換します。

### LED

```go
// engine/keyboard.go から呼び出される
l := led.New(pin machine.Pin, count uint8, deviceType DeviceType, layerColors [MaxLayerColors]color.RGBA)
```

`engine.Peripheral` インターフェースを実装:

| メソッド               | 動作                                            |
|------------------------|-------------------------------------------------|
| `Init()`               | ピンを出力に設定、デバイス種別に応じた初期化、Layer 0 の色で全点灯 |
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
