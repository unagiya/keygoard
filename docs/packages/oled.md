# oled パッケージ仕様

## 概要

SSD1306 OLED ディスプレイの制御を提供するパッケージです。レイヤー変更時に現在のレイヤー番号（"L0"〜"L3"）を画面中央に描画します。ソフトウェアによる 90° 単位の画面回転をサポートしています。`engine.Peripheral` インターフェースを実装しています。

## パッケージパス

`github.com/unagiya/keygoard/engine/peripheral/oled`

## ファイル構成

| ファイル        | machine 依存 | 役割                                  |
|-----------------|:------------:|---------------------------------------|
| `oled.go`       | あり         | `OLED` 構造体・SSD1306 I2C 制御      |
| `font.go`       | なし         | 8x8 ビットマップフォント定義          |
| `font_test.go`  | なし         | フォントデータのユニットテスト        |

## API

### Rotation

```go
type Rotation uint8

const (
    Rotation0   Rotation = 0 // 回転なし（デフォルト）
    Rotation90  Rotation = 1 // 90° 時計回り
    Rotation180 Rotation = 2 // 180° 回転
    Rotation270 Rotation = 3 // 270° 時計回り
)
```

### Config

```go
type Config struct {
    Bus      *machine.I2C   // I2C バス
    SDA      machine.Pin    // I2C データピン
    SCL      machine.Pin    // I2C クロックピン
    Address  uint16         // I2C アドレス（通常 0x3C）
    Width    int16          // 画面幅（ピクセル）
    Height   int16          // 画面高（ピクセル）
    Rotation Rotation       // ソフトウェア画面回転（90° 単位）
}
```

### OLED

```go
display := oled.New(&oled.Config{
    Bus:      machine.I2C0,
    SDA:      machine.GPIO12,
    SCL:      machine.GPIO13,
    Address:  0x3C,
    Width:    128,
    Height:   64,
    Rotation: oled.Rotation180,
})
```

`engine.Peripheral` インターフェースを実装:

| メソッド               | 動作                                            |
|------------------------|-------------------------------------------------|
| `Init()`               | ディスプレイをクリアし、"L0" を描画             |
| `Tick() Keycode`       | keycode.None を返す（毎サイクルの処理なし）     |
| `OnLayerChange(int)`   | レイヤー番号（"L0"〜"L3"）を画面中央に描画     |

## フォント

8x8 ビットマップフォントを fontScale=3 で拡大描画します。レイヤー表示に必要な最小限の文字（'L', '0'〜'3'）のみ定義しています。

| 定数         | 値 | 説明                   |
|--------------|----|------------------------|
| `charWidth`  | 8  | フォントの文字幅（px） |
| `charHeight` | 8  | フォントの文字高（px） |
| `fontScale`  | 3  | 表示時の拡大倍率       |

## 依存パッケージ

- `tinygo.org/x/drivers/ssd1306`（使用許可済み）

## テスト

- `font_test.go`: `glyphIndex` の正当性とフォント定数・グリフデータの検証
- machine 依存部分（`oled.go`）は `make test` の TinyGo ビルドでコンパイル検証を行い、動作確認は実機で実施する
