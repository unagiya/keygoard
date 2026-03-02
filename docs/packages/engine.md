# engine パッケージ仕様

## 概要

キーボードエンジンの公開 Facade パッケージです。利用者は `engine` と `keycode` の 2 パッケージのみ import すれば、キーボードファームウェアを構成できます。マトリクススキャン・レイヤー解決・HID 送信・周辺機器を統合し、`Run()` 一つでキーボードを起動します。

内部実装（`internal/matrix`, `internal/layer`, `internal/tap`, `internal/peripheral/*`）は外部 import 不可です。

## パッケージパス

`github.com/unagiya/keygoard/engine`（公開パッケージ）

## ファイル構成

| ファイル          | machine 依存 | 役割                                         |
|-------------------|:------------:|----------------------------------------------|
| `keyboard.go`     | あり         | エンジン本体（`New` / `Run` / tick ループ）  |
| `config.go`       | あり         | `Config` / `EncoderConfig` / `LEDConfig` / `OLEDConfig` 定義 |
| `keymap.go`       | なし         | `Keymap` 構造体・`RowCount` / `ColCount` / `MaxLayers` 再エクスポート |
| `types.go`        | なし         | `Pin` / `I2CBus` / `Rotation` 型定義         |
| `peripheral.go`   | なし         | `Peripheral` インターフェース定義            |
| `errors.go`       | なし         | エラー変数定義                               |

## API

### ハードウェア抽象型

```go
// GPIO ピン番号（GPIO 番号をそのまま指定: GP5 → Pin(5)）
type Pin uint8

// I2C バス識別子
type I2CBus uint8
const (
    I2C0 I2CBus = 0
    I2C1 I2CBus = 1
)

// ディスプレイ回転角度（90° 単位）
type Rotation uint8
const (
    Rotation0   Rotation = 0 // 回転なし
    Rotation90  Rotation = 1 // 90° 時計回り
    Rotation180 Rotation = 2 // 180°
    Rotation270 Rotation = 3 // 270° 時計回り
)
```

### Config

```go
type Config struct {
    ProductName string        // USB Product Name（空の場合 "keygoard"）
    ColPins     []Pin         // マトリクス列ピン（出力）
    RowPins     []Pin         // マトリクス行ピン（入力プルダウン）
    Keymap      *Keymap       // キーマップ
    Encoder     *EncoderConfig // ロータリーエンコーダー（nil で無効）
    LED         *LEDConfig     // RGB LED（nil で無効）
    OLED        *OLEDConfig    // OLED ディスプレイ（nil で無効）
    Peripherals []Peripheral   // カスタム周辺機器
}
```

- `ProductName` は ASCII のみ、最大 126 文字
- `ColPins` / `RowPins` は GPIO 番号を `Pin` 型で指定する
- マトリクススキャナーは `engine.New()` 内部で自動生成される
- 標準ペリフェラル（Encoder/LED/OLED）も Config から内部生成される

### ペリフェラル設定

```go
type EncoderConfig struct {
    PinA   Pin             // A 信号ピン
    PinB   Pin             // B 信号ピン
    KeyCW  keycode.Keycode // 時計回りキーコード
    KeyCCW keycode.Keycode // 反時計回りキーコード
}

const MaxLayerColors = 4

type LEDConfig struct {
    Pin         Pin                          // データピン
    Count       uint8                        // LED 数
    Type        LEDType                      // デバイス種別（ゼロ値 = WS2812）
    LayerColors [MaxLayerColors]color.RGBA   // レイヤーごとの色
}

type OLEDConfig struct {
    Bus      I2CBus   // I2C バス
    SDA      Pin      // I2C データピン
    SCL      Pin      // I2C クロックピン
    Address  uint16   // I2C アドレス（通常 0x3C）
    Width    int16    // 画面幅（px）
    Height   int16    // 画面高（px）
    Rotation Rotation // ソフトウェア回転
}
```

### Keyboard

```go
// Config からキーボードエンジンを生成
kb := engine.New(&engine.Config{
    ProductName: "zero-kb02",
    ColPins:     []engine.Pin{5, 6, 7, 8},
    RowPins:     []engine.Pin{9, 10, 11},
    Keymap:      keymap,
    Encoder:     &engine.EncoderConfig{...},
    LED:         &engine.LEDConfig{...},
    OLED:        &engine.OLEDConfig{...},
})

// キーボードを起動（この関数は戻らない）
kb.Run()
```

### Keymap

```go
const MaxLayers = 4  // internal/layer.MaxLayers を再エクスポート
const RowCount  = 3  // internal/matrix.RowCount を再エクスポート
const ColCount  = 4  // internal/matrix.ColCount を再エクスポート

type Keymap struct {
    Layers [MaxLayers][RowCount][ColCount]keycode.Keycode
}
```

### Peripheral

```go
const MaxPeripherals = 4

type Peripheral interface {
    Init()
    Tick() keycode.Keycode
    OnLayerChange(layer int)
}
```

- 標準ペリフェラル + カスタムペリフェラルの合計が `MaxPeripherals` 以下であること
- カスタムペリフェラルは `Config.Peripherals` で渡す

## Run() の内部処理

```
1. USB Product Name を設定
2. scanner.Init()（GPIO 初期化）
3. 全周辺機器の Init()
4. USB エニュメレーション完了まで待機（500ms）
5. 初期レイヤー通知（Layer 0）
6. 無限ループ:
   a. tick()（1 スキャンサイクル実行）
   b. time.Sleep(1ms)
```

### tick() の動作

1. 周辺機器の `Tick()` を呼び出し（エンコーダー等のポーリング）
2. `tap.Advance()` で全 pending キーのカウンタをインクリメント
3. 全キーポジションのタイムアウトチェック（pending → holding 遷移時にレイヤー有効化）
4. `scanner.Scan()` でデバウンス済みキー状態を取得
5. 変化がなければレイヤー変更チェックのみ行い即リターン
6. 変化キーについて:
   - 押下: `resolver.Resolve()` → `handlePress()`
   - リリース: `handleRelease()`（タップ判定 → activeKeys に基づく解除）
7. 前回状態を更新
8. レイヤー変更があれば全周辺機器に通知

## 内部依存パッケージ

| パッケージ | 役割 |
|---|---|
| `internal/matrix` | マトリクススキャン・デバウンス |
| `internal/layer` | レイヤー解決 |
| `internal/tap` | タップ/ホールド判定 |
| `internal/peripheral/encoder` | ロータリーエンコーダー |
| `internal/peripheral/led` | RGB LED |
| `internal/peripheral/oled` | OLED ディスプレイ |

Pin → `machine.Pin` の変換は `keyboard.go` 内のヘルパー関数（`pinsToCols`, `pinsToRows`, `i2cBus`）が担います。

## エラー

| 変数             | 意味                           |
|------------------|--------------------------------|
| `ErrKeyOverflow` | 6KRO の同時押し上限を超えた場合 |

## テスト

engine パッケージは machine 依存のため、`make test` の TinyGo ビルドでコンパイル検証を行う。動作確認は実機で実施する。
