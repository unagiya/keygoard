# keygoard

TinyGo ベースのカスタムキーボードファームウェアフレームワークです。

`engine` と `keycode` の 2 パッケージだけで、RP2040 ベースの自作キーボードファームウェアを構成できます。`machine` パッケージを直接扱う必要はありません。

---

## 特徴

- **シンプルな API** — GPIO ピン番号を整数で指定するだけ。ハードウェア抽象層を意識する必要なし
- **2 パッケージで完結** — `engine`（設定・起動）と `keycode`（キーコード定数）のみ
- **レイヤーシステム** — 最大 4 レイヤー。MO / TG / LT / TT によるレイヤー切り替え
- **周辺機器サポート** — ロータリーエンコーダー、RGB LED（WS2812/SK6812）、OLED（SSD1306）
- **ゼロ割り当てホットパス** — スキャンループ内でヒープ割り当てなし

---

## 必要なもの

| ツール | バージョン |
|---|---|
| [Go](https://golang.org/) | 1.25 以上 |
| [TinyGo](https://tinygo.org/) | 0.40.1 以上 |

対応 MCU: RP2040（Waveshare RP2040-Zero 等）

---

## クイックスタート

### 1. プロジェクトを作成

```sh
mkdir my-keyboard && cd my-keyboard
go mod init my-keyboard
go get github.com/unagiya/keygoard
```

### 2. main.go を作成

最小構成（マトリクスのみ、周辺機器なし）の例です。

```go
package main

import (
    _ "machine/usb/hid/keyboard"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    kb := engine.New(&engine.Config{
        ProductName: "my-keyboard",
        ColPins:     []engine.Pin{5, 6, 7, 8},
        RowPins:     []engine.Pin{9, 10, 11},
        Keymap:      defaultKeymap,
    })
    kb.Run()
}

var defaultKeymap = &engine.Keymap{
    Layers: [engine.MaxLayers][engine.RowCount][engine.ColCount]keycode.Keycode{
        // Layer 0
        {
            {keycode.Q, keycode.W, keycode.E, keycode.R},
            {keycode.A, keycode.S, keycode.D, keycode.F},
            {keycode.Z, keycode.X, keycode.C, keycode.Space},
        },
    },
}
```

> `_ "machine/usb/hid/keyboard"` は USB HID ハンドラの登録に必要です。

### 3. ビルドしてフラッシュ

```sh
tinygo build -target=waveshare-rp2040-zero .
tinygo flash -target=waveshare-rp2040-zero .
```

---

## キーマップの定義

キーマップは `engine.Keymap` 構造体で定義します。最大 4 レイヤーをサポートし、`Layers[0]` がベースレイヤー（常に有効）です。番号が大きいレイヤーほど優先度が高くなります。

```go
var myKeymap = &engine.Keymap{
    Layers: [engine.MaxLayers][engine.RowCount][engine.ColCount]keycode.Keycode{
        // Layer 0: ベース
        {
            {keycode.Q, keycode.W, keycode.E, keycode.R},
            {keycode.A, keycode.S, keycode.D, keycode.F},
            {keycode.Z, keycode.X, keycode.C, keycode.MO(1)},
        },
        // Layer 1: MO(1) ホールド中に有効
        {
            {keycode.Num1, keycode.Num2, keycode.Num3, keycode.Num4},
            {keycode.Num5, keycode.Num6, keycode.Num7, keycode.Num8},
            {keycode.Num9, keycode.Num0, keycode.TRNS, keycode.TRNS},
        },
    },
}
```

### レイヤー操作キー

| キーコード | 動作 |
|---|---|
| `keycode.MO(n)` | ホールド中のみレイヤー n を有効化 |
| `keycode.TG(n)` | 押すたびにレイヤー n を ON/OFF |
| `keycode.LT(n, kc)` | タップで kc を入力、ホールドでレイヤー n を有効化 |
| `keycode.TT(n)` | タップでレイヤー n をトグル、ホールドで MO と同じ |
| `keycode.TRNS` | 透過（下位レイヤーのキーコードを参照） |
| `keycode.None` | 何も割り当てない |

> LT / TT のタップ/ホールド判定閾値は 200ms です。

---

## 周辺機器の設定

周辺機器はすべてオプションです。`Config` の対応フィールドを `nil` のままにすれば無効になります。

### ロータリーエンコーダー

```go
Encoder: &engine.EncoderConfig{
    PinA:   3,              // A 信号の GPIO 番号
    PinB:   4,              // B 信号の GPIO 番号
    KeyCW:  keycode.UpArrow,   // 時計回りで送信するキー
    KeyCCW: keycode.DownArrow, // 反時計回りで送信するキー
},
```

### RGB LED（WS2812 / SK6812）

```go
LED: &engine.LEDConfig{
    Pin:   1,       // データピンの GPIO 番号
    Count: 12,      // LED の数
    Type:  engine.WS2812, // engine.WS2812 または engine.SK6812（省略時: WS2812）
    LayerColors: [engine.MaxLayerColors]color.RGBA{
        {G: 8},              // Layer 0: 緑
        {B: 8},              // Layer 1: 青
        {R: 8},              // Layer 2: 赤
        {R: 4, G: 4, B: 4},  // Layer 3: 白
    },
},
```

> `image/color` パッケージの import が必要です。

### OLED ディスプレイ（SSD1306）

```go
OLED: &engine.OLEDConfig{
    Bus:      engine.I2C0, // I2C バス（I2C0 または I2C1）
    SDA:      12,          // I2C データピンの GPIO 番号
    SCL:      13,          // I2C クロックピンの GPIO 番号
    Address:  0x3C,        // I2C アドレス
    Width:    128,          // 画面幅（px）
    Height:   64,           // 画面高（px）
    Rotation: engine.Rotation180, // 画面回転（0 / 90 / 180 / 270）
},
```

レイヤー切り替え時に現在のレイヤー番号（"L0"〜"L3"）を画面中央に表示します。

---

## Config リファレンス

```go
engine.Config{
    ProductName string           // USB デバイス名（省略時: "keygoard"）
    ColPins     []engine.Pin     // マトリクス列ピン（GPIO 番号）
    RowPins     []engine.Pin     // マトリクス行ピン（GPIO 番号）
    Keymap      *engine.Keymap   // キーマップ
    Encoder     *EncoderConfig   // ロータリーエンコーダー（nil で無効）
    LED         *LEDConfig       // RGB LED（nil で無効）
    OLED        *OLEDConfig      // OLED ディスプレイ（nil で無効）
    Peripherals []Peripheral     // カスタム周辺機器（上級者向け）
}
```

`engine.Pin` は GPIO ピン番号を整数で指定する型です。GP5 なら `5` と書きます。

---

## ビルドとフラッシュ

```sh
# ビルド確認
tinygo build -target=waveshare-rp2040-zero .

# 実機への書き込み
tinygo flash -target=waveshare-rp2040-zero .
```

`-target` にはお使いのボードの TinyGo ターゲット名を指定してください。

---

## キーコード一覧

代表的なキーコードです。全定数は [`keycode/keycode.go`](keycode/keycode.go) を参照してください。

**アルファベット:** `keycode.A` 〜 `keycode.Z`

**数字:** `keycode.Num0` 〜 `keycode.Num9`

**修飾キー:**

| キーコード | キー |
|---|---|
| `keycode.ModLeftCtrl` | 左 Ctrl |
| `keycode.ModLeftShift` | 左 Shift |
| `keycode.ModLeftAlt` | 左 Alt |
| `keycode.ModLeftGUI` | 左 GUI（Command / Win） |
| `keycode.ModRightCtrl` | 右 Ctrl |
| `keycode.ModRightShift` | 右 Shift |
| `keycode.ModRightAlt` | 右 Alt |
| `keycode.ModRightGUI` | 右 GUI |

**基本操作:**

| キーコード | キー |
|---|---|
| `keycode.Enter` | Enter |
| `keycode.Escape` | Escape |
| `keycode.Backspace` | Backspace |
| `keycode.Tab` | Tab |
| `keycode.Space` | Space |

**矢印キー:** `keycode.UpArrow`, `keycode.DownArrow`, `keycode.LeftArrow`, `keycode.RightArrow`

**ファンクションキー:** `keycode.F1` 〜 `keycode.F12`

---

## 実装例

[`examples/zero-kb02/`](examples/zero-kb02/) に、エンコーダー・LED・OLED を含むフル機能の実装例があります。

---

## ロードマップ

今後の開発計画は [ROADMAP.md](ROADMAP.md) を参照してください。Joystick + Gamepad 対応、Split キーボード対応などを予定しています。

---

## 開発者向け情報

フレームワーク自体の開発に参加する場合は以下を参照してください。

### 必要な追加ツール

- `goimports`（フォーマット用）
- `golangci-lint`（Lint 用）

### 開発用コマンド

```sh
make build   # TinyGo ビルド確認
make test    # TinyGo ビルド + ユニットテスト
make flash   # 実機書き込み
make fmt     # コードフォーマット
make lint    # Lint
```

### ドキュメント

- [`docs/`](docs/) — 設計ドキュメント（アーキテクチャ・機能設計・パッケージ仕様）