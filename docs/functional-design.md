# 機能設計書

> **対象外（本プロジェクトに UI がないため）:** 画面遷移・ワイヤーフレーム
> **対象外（RDB を持たないため）:** ER 図
> これらを除いた項目を以下に記載する。

---

## システム構成図

### パッケージ依存関係

```
┌─────────────────────────────────────────────┐
│  examples/zero-kb02/                        │
│  エントリポイント（別モジュール）            │
└────────────────┬────────────────────────────┘
                 │ 呼び出し
┌────────────────▼────────────────────────────┐
│  engine/                                    │
│  メインループ（スキャン → 解決 → HID 送信） │
│  Peripheral インターフェースで周辺機器を統合 │
└──┬──────────┬──────────────┬────────────────┘
   │          │              │
┌──▼────┐ ┌──▼────────┐ ┌───▼──────────────┐
│matrix/│ │ keycode/  │ │ peripheral/      │
│スキャン│ │ 定数定義  │ │ ├ encoder/       │
│デバウンス│ └──────────┘ │ ├ led/           │
└──┬────┘              │ └ oled/           │
   │                   └──────────┬────────┘
   │    ※ engine/ も machine に直接依存（USB HID）
┌──▼──────────────────────────────▼───────────┐
│  machine + tinygo.org/x/drivers             │
│  GPIO / USB HID / WS2812 / SSD1306 / I2C   │
└─────────────────────────────────────────────┘
       │
┌──────▼──────────────────────────────────────┐
│  RP2040 ハードウェア（zero-kb02）            │
│  マトリクス GPIO / USB / LED / OLED / Enc   │
└─────────────────────────────────────────────┘
```

### USB 接続構成（Phase 1）

```
zero-kb02（RP2040）
  └─ USB Full-Speed（12Mbps）
        └─ macOS / Windows / Linux ホスト
             └─ HID Usage Page 0x01 / Usage 0x06（Keyboard）
                  └─ 6KRO キーボードレポート
```

### 周辺機器構成（Phase 4 で実装済み）

```
engine/
  ├── peripheral.go        ← Peripheral インターフェース定義
  ├── matrix/
  ├── keycode/
  └── peripheral/
        ├── encoder/       ← ロータリーエンコーダー（GP3/GP4）
        ├── led/           ← RGB LED WS2812/SK6812（GP1、12 LED）
        ├── oled/          ← OLED SSD1306（GP12/GP13、I2C）
        └── joystick/      ← Phase 5 で追加予定
```

`engine` → `peripheral` の依存は `engine.Peripheral` インターフェース経由。

---

## 機能ごとのアーキテクチャ

### マトリクススキャン（matrix パッケージ）

COL2ROW 方式。列を 1 本ずつ High にして全行を読み取る。

```
Col ピン（出力）: GP5, GP6, GP7, GP8
Row ピン（入力プルダウン）: GP9, GP10, GP11

スキャン手順（1列あたり）:
  1. 対象 Col を High
  2. 10µs 待機（GPIO 安定化）
  3. 全 Row 読み取り
  4. Col を Low に戻す
  → 全 4 列分繰り返す
```

スキャン結果は `[Rows][Cols]bool` の固定配列に書き込む（ヒープ割り当てなし）。

### デバウンス（matrix パッケージ内）

カウンタ方式。ホットパスが `machine` 非依存になるよう `matrix.go` と分離する。

```
raw（生スキャン結果）
  ↓ カウンタインクリメント
  閾値（5 サイクル）に達したら stable に反映
  ↓
stable（確定済みキー状態）
```

`debounce.go` は `//go:build tinygo` タグを持たないため、標準 Go でユニットテスト可能。

### USB HID キーボード（engine パッケージ内）

TinyGo 標準の `machine/usb/hid/keyboard` を使う。自前の HID 記述子は書かない。
HID 送信ロジックは独立パッケージではなく `engine/keyboard.go` に統合している。

```
起動シーケンス:
  1. engine.New(scanner, keymap) で Keyboard を生成
  2. kb.Init()（GPIO 初期化 + 500ms 待機でエニュメレーション完了を待つ）
  3. メインループ: kb.Tick() を 1ms 周期で呼び出す

レポート送信（6KRO）:
  - キー状態に変化があった場合のみ hidkb.Keyboard.Down/Up を呼び出す
  - Modifier byte（1 バイト）
  - Keycode 6 バイト（同時押し最大 6 キー）
```

### キーコード定義（keycode パッケージ）

HID Usage Page 7（Keyboard/Keypad）の Usage ID を定数として定義する。
`machine` 非依存のため標準 Go でテスト可能。

### メインループ（engine パッケージ）

1 goroutine で動かす。全作業用配列は初期化時に確保し、ループ内でのヒープ割り当てを防ぐ。

```
初期化（kb.Init）
  ├─ scanner.Init()（GPIO ピン設定）
  └─ USB エニュメレーション完了待ち（500ms）

ループ（1ms 周期、examples/zero-kb02/main.go が time.Sleep で制御）
  └─ kb.Tick()
       ├─ scanner.Scan() → state, changed
       ├─ changed == false なら即 return
       └─ 変化キーを検出 → hidkb.Keyboard.Down/Up 呼び出し
```

---

## データ構造定義

### 主要な型

```go
// キーコード（HID Usage ID ベース）
// 0x0000: 無効、0xE0xx: 修飾キー、0xF0xx: 通常キー
type Keycode uint16

// マトリクスサイズ
const RowCount = 3
const ColCount = 4

// デバウンス状態
type Debouncer struct {
    debounced [RowCount][ColCount]bool  // 確定済みキー状態
    counter   [RowCount][ColCount]uint8 // 連続サイクルカウンタ
}

// HID レポートの送信は machine/usb/hid/keyboard が担う（独自型なし）
```

### データフロー

```
物理キー押下（GPIO）
    │
    ▼ [RowCount][ColCount]bool
matrix.Scan()（内部でデバウンス処理）
    │ state [RowCount][ColCount]bool, changed bool
    ▼
engine.tick()
    ├─ tap.Advance()          ← 全 pending キーのカウンタ++
    ├─ タイムアウトチェック    ← pending → holding 遷移時にレイヤー有効化
    ├─ resolver.Resolve()     ← 最上位アクティブレイヤーから走査
    ├─ handlePress()          ← MO/TG/LT/TT/通常キー分岐
    └─ handleRelease()        ← タップ判定 → activeKeys に基づく解除
    │
    ▼
USB HID レポート → OS
```

### レイヤー解決

```
resolver.Resolve(row, col):
  for layer = MaxLayers-1 downto 0:
    if !active[layer]: continue
    kc = Layers[layer][row][col]
    if kc != TRNS: return kc
  return None
```

- 最上位のアクティブレイヤーから下方向に走査
- TRNS（透過キー）は下位レイヤーにフォールバック
- すべて TRNS の場合は None を返す

### タップ/ホールド判定

```
tapIdle ──[Press(LT/TT)]──→ tapPending ──[閾値超過(200ms)]──→ tapHolding
  ▲                              │                                │
  └──────[Release=タップ]────────┘                                │
  └──────────────────────[Release=ホールド解除]───────────────────┘

タップ時:
  LT(n, kc): kc を Down → Up（即時送信）
  TT(n):     レイヤー n をトグル

ホールド時:
  LT/TT:     レイヤー n を Activate（リリースで Deactivate）
```

### キーマップ構造

```go
const MaxLayers = 4

type Keymap struct {
    Layers [MaxLayers][RowCount][ColCount]Keycode
}
```

- `Layers[0]` がベースレイヤー（常に有効）
- 番号が大きいほど優先度が高い
- 未使用レイヤーはゼロ値（全キー None）のまま

---

## コンポーネント設計

### keycode パッケージ

**責務:** HID キーコードの型・定数定義。他パッケージに依存しない。

```go
// 公開定数（抜粋）
const (
    None Keycode = 0x0000

    // 修飾キー
    ModLeftCtrl  Keycode = 0xE000
    ModLeftShift Keycode = 0xE001
    // ...

    // 通常キー（HID Usage ID + 0xF000）
    A Keycode = 0xF004
    B Keycode = 0xF005
    // ...
)
```

**テスト:** `go test ./keycode/...` で実行可能。

---

### matrix パッケージ

**責務:** GPIO スキャンとデバウンスの 2 責務を持つ。ファイルで分離する。

| ファイル | ビルドタグ | 責務 |
|---|---|---|
| `const.go` | なし | マトリクスサイズ定数 |
| `debounce.go` | なし | デバウンスロジック（machine 非依存） |
| `debounce_test.go` | なし | ユニットテスト |
| `matrix.go` | `//go:build tinygo` | GPIO スキャン実装 |

```go
// Scanner 生成
s := matrix.New(cols, rows)
s.Init()

// スキャン（メインループで呼び出す）
state, changed := s.Scan()
// state:   [RowCount][ColCount]bool
// changed: 状態変化があった場合 true
```

---

### engine パッケージ

**責務:** 全コンポーネントを統合するメインループと HID 送信。goroutine は 1 つのみ使用する。
HID 送信は `machine/usb/hid/keyboard` を直接利用し、独立した `hid/` パッケージは持たない。

```go
// Keyboard を生成する
func New(scanner *matrix.Scanner, km *Keymap) *Keyboard

// ハードウェア初期化（GPIO 初期化 + USB エニュメレーション待機）
func (kb *Keyboard) Init()

// 1 スキャンサイクルを実行する（メインループから毎回呼び出す）
func (kb *Keyboard) Tick()
```

**ビルドタグ:** `engine/keyboard.go` は `//go:build tinygo`（標準 Go ビルドから除外）

---

## ユースケース

本プロジェクトはファームウェアであるため、従来の画面遷移やワイヤーフレームは存在しない。
以下にキーボードとしての主要なユースケースを記述する。

### UC-01: キー入力

```
アクター: ユーザー
事前条件: zero-kb02 が USB 接続され OS にキーボードとして認識済み

基本フロー:
  1. ユーザーが物理キーを押下する
  2. マトリクススキャンが押下を検出する
  3. デバウンスが 5 サイクル連続検出で状態を確定する
  4. レイヤー 0 のキーマップから対応するキーコードを解決する
  5. HID レポートを USB 経由で OS に送信する
  6. OS のフォーカスアプリに文字が入力される

代替フロー（チャタリング）:
  3a. 5 サイクル未満の場合は状態を確定しない → UC-01 終了
```

### UC-02: キーリリース

```
アクター: ユーザー
事前条件: UC-01 完了済み（キーが押下された状態）

基本フロー:
  1. ユーザーが物理キーを離す
  2. マトリクススキャンが解放を検出する
  3. デバウンスが 5 サイクル連続検出で解放状態を確定する
  4. 空の HID レポートを送信する
  5. OS がキーリリースイベントを受け取る
```

### UC-03: 複数キー同時押し（Phase 1）

```
アクター: ユーザー
制約: 6KRO（最大 6 キー同時押し）。7 キー以上は先着 6 キーのみ有効。

基本フロー:
  1. ユーザーが複数のキーを押下する
  2. 各キーのデバウンスが独立して確定する
  3. stableState から押下中キーを最大 6 つ収集する
  4. 6KRO HID レポートを送信する
```

---

## API 設計

各パッケージの公開 API 一覧。詳細は `docs/packages/` の各仕様書を参照。

### keycode パッケージ

| 識別子 | 種別 | 説明 |
|---|---|---|
| `Keycode` | 型（`uint16`） | HID キーコード |
| `None` | 定数 | 無効（未割り当て） |
| `ModLeftCtrl` … | 定数 | 修飾キー |
| `A` … `Z` | 定数 | アルファベットキー |
| `Num0` … `Num9` | 定数 | 数字キー |
| `Enter`, `Space` … | 定数 | 基本操作キー |
| `F1` … `F12` | 定数 | ファンクションキー |

### matrix パッケージ

| 識別子 | 種別 | 説明 |
|---|---|---|
| `RowCount`, `ColCount` | 定数 | マトリクスサイズ（3, 4） |
| `Scanner` | 構造体 | スキャナー本体（GPIO + Debouncer を内包） |
| `New(cols, rows)` | コンストラクタ | Scanner を生成して返す |
| `(*Scanner).Init()` | メソッド | GPIO ピン初期化 |
| `(*Scanner).Scan()` | メソッド | スキャン実行、`(state [RowCount][ColCount]bool, changed bool)` を返す |
| `Debouncer` | 構造体 | デバウンス状態管理（Scanner に内包） |
| `(*Debouncer).Update(raw)` | メソッド | デバウンス処理を実行し `(state, changed)` を返す |

### engine パッケージ

| 識別子 | 種別 | 説明 |
|---|---|---|
| `Keyboard` | 構造体 | キーボードエンジン本体 |
| `Keymap` | 構造体 | 複数レイヤーキーマップ（`Layers [MaxLayers][RowCount][ColCount]Keycode`） |
| `MaxLayers` | 定数 | 最大レイヤー数（4） |
| `Resolver` | 構造体 | レイヤー解決（最上位アクティブレイヤーから走査） |
| `TapDetector` | 構造体 | タップ/ホールド判定（LT/TT キー用） |
| `Peripheral` | インターフェース | 周辺機器共通（`Init()`, `Tick()`, `OnLayerChange()`) |
| `MaxPeripherals` | 定数 | 登録可能な周辺機器の最大数（4） |
| `Config` | 構造体 | エンジン設定（Scanner, Keymap, ProductName, Peripherals） |
| `New(cfg)` | コンストラクタ | Config から Keyboard を生成して返す |
| `(*Keyboard).Run()` | メソッド | キーボード起動（初期化 + 無限ループ、戻らない） |
| `ErrKeyOverflow` | エラー | 6KRO 上限超過（7 キー以上同時押し） |
