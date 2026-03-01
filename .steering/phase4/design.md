# Phase 4 設計書: 周辺機器（エンコーダー・LED・OLED）

## 実装アプローチ

### 全体方針

- 周辺機器ごとに独立パッケージを作成し、`engine` からはインターフェース経由で呼び出す
- machine 非依存のロジック（クアドラチャデコード、フォントビットマップ）はビルドタグなしで実装し、ユニットテスト可能にする
- `tinygo-org/drivers`（使用許可済み）の WS2812 / SSD1306 ドライバを使用する

---

## パッケージ構成

```
engine/
├── peripheral.go            ← Peripheral インターフェース定義（ビルドタグなし）
├── keyboard.go              ← Keyboard 構造体に周辺機器統合を追加（//go:build tinygo）
├── config.go                ← Config に Peripherals フィールドを追加（//go:build tinygo）
└── ...（既存ファイルは変更なし）

engine/peripheral/
├── encoder/
│   ├── encoder.go           ← Encoder 構造体・GPIO 読み取り（//go:build tinygo）
│   ├── quadrature.go        ← Decoder（クアドラチャ状態遷移ロジック、ビルドタグなし）
│   └── quadrature_test.go   ← Decoder ユニットテスト
├── led/
│   └── led.go               ← LED 構造体・WS2812 制御（//go:build tinygo）
└── oled/
    ├── oled.go              ← OLED 構造体・SSD1306 制御（//go:build tinygo）
    └── font.go              ← 最小ビットマップフォント（ビルドタグなし）
```

### パッケージ依存関係

```
examples/zero-kb02/
    │ import
    ▼
engine/ ──────────────────────────────── engine/peripheral/encoder/
    │ interface 定義                         │ import keycode
    │                                        ▼
    ├──── engine/peripheral/led/ ────── keycode/
    │         │ import
    │         ▼
    │    tinygo-org/drivers/ws2812
    │
    ├──── engine/peripheral/oled/ ───── keycode/
    │         │ import
    │         ▼
    │    tinygo-org/drivers/ssd1306
    │
    ├──── matrix/
    └──── keycode/
```

- `engine` は `engine/peripheral/*` を import しない（循環依存回避）
- `engine/peripheral/*` は `engine` を import しない
- `examples/zero-kb02` が全てのパッケージを import してワイヤリングする

---

## インターフェース設計

### Peripheral インターフェース

`engine/peripheral.go`（ビルドタグなし）に定義する。

```go
// Peripheral はエンジンに統合される周辺機器のインターフェースです。
type Peripheral interface {
    // Init はハードウェアを初期化します。
    Init()

    // Tick は 1 スキャンサイクルの処理を実行します。
    // 送信すべきキーコードがある場合はそのキーコードを返します。
    // 送信不要の場合は keycode.None を返します。
    Tick() keycode.Keycode

    // OnLayerChange はアクティブレイヤーが変化したときに呼び出されます。
    // layer は最上位のアクティブレイヤー番号です。
    OnLayerChange(layer int)
}
```

### 定数

```go
// MaxPeripherals はエンジンに登録できる周辺機器の最大数です。
const MaxPeripherals = 4
```

---

## コンポーネント設計

### 1. エンコーダー（`engine/peripheral/encoder`）

#### Decoder（クアドラチャデコード）

`quadrature.go`（ビルドタグなし）に配置。ユニットテスト可能。

```go
// Direction は回転方向を表します。
type Direction uint8

const (
    DirNone Direction = iota
    DirCW   // 時計回り
    DirCCW  // 反時計回り
)

// Decoder は 2 相エンコーダーの回転方向を判定します。
type Decoder struct {
    prev uint8 // 前回の状態（A<<1 | B）
}

// Update は現在の A/B 信号を受け取り、回転方向を返します。
func (d *Decoder) Update(a, b bool) Direction
```

**状態遷移テーブル（ルックアップ方式）:**

```
prev\curr  00    01    10    11
  00      None   CW   CCW   None
  01      CCW   None  None   CW
  10       CW   None  None  CCW
  11      None  CCW    CW   None
```

4×4 の固定配列で O(1) 判定。ホットパスに適した実装。

#### Encoder 構造体

`encoder.go`（`//go:build tinygo`）に配置。

```go
// Config はエンコーダーの設定です。
type Config struct {
    PinA   machine.Pin     // A 信号ピン（GP3）
    PinB   machine.Pin     // B 信号ピン（GP4）
    KeyCW  keycode.Keycode // 時計回りに割り当てるキーコード
    KeyCCW keycode.Keycode // 反時計回りに割り当てるキーコード
}

// Encoder はロータリーエンコーダーの入力を処理します。
type Encoder struct {
    pinA    machine.Pin
    pinB    machine.Pin
    keyCW   keycode.Keycode
    keyCCW  keycode.Keycode
    decoder Decoder
}

func New(cfg *Config) *Encoder
func (e *Encoder) Init()                    // GPIO ピンを入力（プルアップ）に設定
func (e *Encoder) Tick() keycode.Keycode    // A/B 読み取り → Decoder 判定 → キーコード返却
func (e *Encoder) OnLayerChange(layer int)  // no-op
```

**Tick() のフロー:**

```
1. pinA.Get(), pinB.Get() で現在の信号を読み取る
2. decoder.Update(a, b) で回転方向を判定
3. CW → keyCW を返す / CCW → keyCCW を返す / None → keycode.None
```

エンコーダーの回転は「タップ」として処理する（Down → Up を即時送信）。

#### HID 送信の注意

エンコーダーに割り当て可能なキーコードは現在の `keycode` パッケージに定義済みの通常キー・修飾キーのみ。
Consumer Control キー（音量 UP / DOWN 等）は Composite HID が必要なため Phase 5 スコープ。
Phase 4 の example では UpArrow / DownArrow を割り当てる。

---

### 2. RGB LED（`engine/peripheral/led`）

`led.go`（`//go:build tinygo`）に配置。

```go
// MaxLEDs は制御可能な LED の最大数です。
const MaxLEDs = 12

// MaxLayerColors はレイヤーごとの色設定の最大数です。
const MaxLayerColors = 4

// Config は LED の設定です。
type Config struct {
    Pin         machine.Pin                   // データピン（GP1）
    Count       uint8                         // LED 数
    LayerColors [MaxLayerColors]color.RGBA    // レイヤーごとの LED 色
}

// LED は WS2812 互換の RGB LED を制御します。
type LED struct {
    ws          ws2812.Device
    count       uint8
    layerColors [MaxLayerColors]color.RGBA
    buf         [MaxLEDs]color.RGBA
}

func New(cfg *Config) *LED
func (l *LED) Init()                    // WS2812 デバイス初期化、デフォルト色で全点灯
func (l *LED) Tick() keycode.Keycode    // keycode.None を返す（処理なし）
func (l *LED) OnLayerChange(layer int)  // 全 LED をレイヤーに対応する色に更新
```

**OnLayerChange() のフロー:**

```
1. layerColors[layer] から色を取得
2. buf の全要素に色を設定
3. ws.WriteColors(buf[:count]) で LED に書き込む
```

**デフォルト色の設定例:**

| レイヤー | 色 | RGBA |
|---|---|---|
| Layer 0 | 緑（dim） | `{R: 0, G: 8, B: 0}` |
| Layer 1 | 青（dim） | `{R: 0, G: 0, B: 8}` |
| Layer 2 | 赤（dim） | `{R: 8, G: 0, B: 0}` |
| Layer 3 | 白（dim） | `{R: 4, G: 4, B: 4}` |

輝度は低め（最大 255 のうち 8 程度）に設定し、消費電力と眩しさを抑える。

---

### 3. OLED ディスプレイ（`engine/peripheral/oled`）

#### 最小ビットマップフォント

`font.go`（ビルドタグなし）に配置。サードパーティ不要。

```go
// charWidth はフォントの文字幅（ピクセル）です。
const charWidth = 8

// charHeight はフォントの文字高（ピクセル）です。
const charHeight = 8

// fontData は最小ビットマップフォントです。
// 表示に必要な文字のみ定義します: 'L', '0', '1', '2', '3'
// 各文字は charHeight バイトで、各バイトの各ビットが 1 ピクセルに対応します。
var fontData = map[byte][charHeight]byte{
    'L': { ... },
    '0': { ... },
    '1': { ... },
    '2': { ... },
    '3': { ... },
}
```

**注意:** `map` はホットパスでは使用しない。フォント参照はレイヤー変更時のみ（OnLayerChange 内）のため許容。

#### OLED 構造体

`oled.go`（`//go:build tinygo`）に配置。

```go
// Config は OLED の設定です。
type Config struct {
    Bus     *machine.I2C     // I2C バス
    SDA     machine.Pin      // SDA ピン（GP12）
    SCL     machine.Pin      // SCL ピン（GP13）
    Address uint16           // I2C アドレス（0x3C）
    Width   int16            // 画面幅（128）
    Height  int16            // 画面高（64）
}

// OLED は SSD1306 OLED ディスプレイを制御します。
type OLED struct {
    dev     ssd1306.Device
    width   int16
    height  int16
}

func New(cfg *Config) *OLED
func (o *OLED) Init()                    // I2C 初期化、SSD1306 初期化、初期表示
func (o *OLED) Tick() keycode.Keycode    // keycode.None を返す（処理なし）
func (o *OLED) OnLayerChange(layer int)  // "L0"〜"L3" を表示
```

**OnLayerChange() のフロー:**

```
1. dev.ClearDisplay()
2. fontData から 'L' と digit のビットマップを取得
3. dev.SetPixel() で各ピクセルを描画
4. dev.Display() で画面を更新
```

表示位置は画面中央付近にレンダリングする。

---

## engine 統合設計

### Config の変更

```go
// Config（//go:build tinygo）
type Config struct {
    Scanner     *matrix.Scanner
    Keymap      *Keymap
    ProductName string
    Peripherals []Peripheral    // 追加: 周辺機器リスト（最大 MaxPeripherals）
}
```

### Keyboard 構造体の変更

```go
type Keyboard struct {
    // ... 既存フィールド ...
    peripherals     [MaxPeripherals]Peripheral  // 追加
    peripheralCount int                         // 追加
    prevTopLayer    int                         // 追加: レイヤー変更検出用
}
```

### New() の変更

```go
func New(cfg *Config) *Keyboard {
    kb := &Keyboard{ ... }
    // Peripherals をコピー（スライスから固定配列へ）
    for i := 0; i < len(cfg.Peripherals) && i < MaxPeripherals; i++ {
        kb.peripherals[i] = cfg.Peripherals[i]
        kb.peripheralCount++
    }
    return kb
}
```

### setup() の変更

```go
func (kb *Keyboard) setup() {
    kb.scanner.Init()
    // 周辺機器を初期化
    for i := 0; i < kb.peripheralCount; i++ {
        kb.peripherals[i].Init()
    }
    // USB エニュメレーション完了待ち
    _ = machine.USBDev
    time.Sleep(500 * time.Millisecond)
    // 初期レイヤー通知
    kb.notifyLayerChange(0)
}
```

### tick() の変更

```go
func (kb *Keyboard) tick() {
    prevTop := kb.topLayer()

    // 周辺機器の Tick（エンコーダーのスキャン）
    for i := 0; i < kb.peripheralCount; i++ {
        if kc := kb.peripherals[i].Tick(); kc != keycode.None {
            // エンコーダー等からのキーコードはタップとして送信
            hidkb.Keyboard.Down(hidkb.Keycode(kc))
            hidkb.Keyboard.Up(hidkb.Keycode(kc))
        }
    }

    // タップカウンタ・タイムアウト（既存処理）
    kb.tap.Advance()
    // ... 既存のキースキャン・レイヤー解決 ...

    // レイヤー変更検出
    currTop := kb.topLayer()
    if currTop != prevTop {
        kb.notifyLayerChange(currTop)
        kb.prevTopLayer = currTop
    }
}
```

### 新規メソッド

```go
// topLayer は最上位のアクティブレイヤー番号を返します。
func (kb *Keyboard) topLayer() int {
    for l := MaxLayers - 1; l >= 0; l-- {
        if kb.resolver.IsActive(l) {
            return l
        }
    }
    return 0
}

// notifyLayerChange は全周辺機器にレイヤー変更を通知します。
func (kb *Keyboard) notifyLayerChange(layer int) {
    for i := 0; i < kb.peripheralCount; i++ {
        kb.peripherals[i].OnLayerChange(layer)
    }
}
```

---

## examples/zero-kb02 の変更

### config.go に追加

```go
var (
    encoderPinA = machine.GPIO3
    encoderPinB = machine.GPIO4
    ledPin      = machine.GPIO1
    oledSDA     = machine.GPIO12
    oledSCL     = machine.GPIO13
)
```

### main.go の変更

```go
func main() {
    enc := encoder.New(&encoder.Config{
        PinA:   encoderPinA,
        PinB:   encoderPinB,
        KeyCW:  keycode.UpArrow,
        KeyCCW: keycode.DownArrow,
    })

    leds := led.New(&led.Config{
        Pin:   ledPin,
        Count: 12,
        LayerColors: [led.MaxLayerColors]color.RGBA{
            {G: 8},           // Layer 0: 緑
            {B: 8},           // Layer 1: 青
            {R: 8},           // Layer 2: 赤
            {R: 4, G: 4, B: 4}, // Layer 3: 白
        },
    })

    display := oled.New(&oled.Config{
        Bus:     machine.I2C0,
        SDA:     oledSDA,
        SCL:     oledSCL,
        Address: 0x3C,
        Width:   128,
        Height:  64,
    })

    kb := engine.New(&engine.Config{
        Scanner:     matrix.New(matrixCols, matrixRows),
        Keymap:      defaultKeymap,
        ProductName: "zero-kb02",
        Peripherals: []engine.Peripheral{enc, leds, display},
    })
    kb.Run()
}
```

---

## 影響範囲

### 変更するファイル

| ファイル | 変更内容 |
|---|---|
| `engine/config.go` | `Peripherals` フィールド追加 |
| `engine/keyboard.go` | 周辺機器統合（setup/tick/topLayer/notifyLayerChange） |
| `examples/zero-kb02/main.go` | 周辺機器の初期化・登録 |
| `examples/zero-kb02/config.go` | エンコーダー・LED・OLED のピン定義追加 |
| `examples/zero-kb02/go.mod` | `tinygo-org/drivers` 依存追加 |
| `go.mod` | `tinygo-org/drivers` 依存追加（peripheral パッケージが使用） |

### 新規作成するファイル

| ファイル | 内容 |
|---|---|
| `engine/peripheral.go` | Peripheral インターフェース・MaxPeripherals 定数 |
| `engine/peripheral/encoder/encoder.go` | Encoder 構造体（`//go:build tinygo`） |
| `engine/peripheral/encoder/quadrature.go` | Decoder（ビルドタグなし） |
| `engine/peripheral/encoder/quadrature_test.go` | Decoder ユニットテスト |
| `engine/peripheral/led/led.go` | LED 構造体（`//go:build tinygo`） |
| `engine/peripheral/oled/oled.go` | OLED 構造体（`//go:build tinygo`） |
| `engine/peripheral/oled/font.go` | 最小ビットマップフォント（ビルドタグなし） |
| `docs/packages/encoder.md` | エンコーダーパッケージ仕様 |
| `docs/packages/led.md` | LED パッケージ仕様 |
| `docs/packages/oled.md` | OLED パッケージ仕様 |

### 変更しないファイル

- `engine/layer.go` — Resolver は変更なし（`IsActive()` を既存で利用）
- `engine/tap.go` — TapDetector は変更なし
- `engine/keymap.go` — Keymap 構造は変更なし
- `matrix/` — マトリクス関連は変更なし
- `keycode/` — キーコード定義は変更なし

---

## メモリ影響見積もり

| コンポーネント | 追加 RAM |
|---|---|
| Encoder（Decoder + ピン + キーコード） | ~20 bytes |
| LED（WS2812 デバイス + 色バッファ 12×4 bytes） | ~100 bytes |
| OLED（SSD1306 デバイス + フレームバッファ 128×64/8） | ~1,100 bytes |
| Keyboard 追加フィールド（Peripheral 配列 + count + prevTopLayer） | ~40 bytes |
| **合計** | **~1.3 KB** |

128KB 制約内で十分許容範囲。
