# Phase 3: レイヤーシステム — 設計書

## 実装アプローチ

### 方針

レイヤーロジックは `engine` パッケージ内にビルドタグなしのファイルとして追加する。
これにより `machine` 非依存でユニットテスト可能な構造を維持する。

新規パッケージは作らず、`engine/layer.go` と `engine/tap.go` に分離して配置する。
`keyboard.go`（`//go:build tinygo`）がこれらを統合する。

---

## キーコードエンコーディング拡張

### 現在のエンコーディング（`keycode.Keycode = uint16`）

| 範囲 | 意味 |
|------|------|
| `0x0000` | None（無効） |
| `0xE000-0xE0FF` | 修飾キー |
| `0xF000-0xFFFF` | 通常キー（HID Usage Page 7） |

### Phase 3 で追加するエンコーディング

| 範囲 | 意味 | エンコーディング |
|------|------|------------------|
| `0x0001` | TRNS（透過） | 固定値 |
| `0xA000-0xA00F` | MO(layer) | `0xA000 \| layer` |
| `0xA010-0xA01F` | TG(layer) | `0xA010 \| layer` |
| `0xA020-0xA02F` | TT(layer) | `0xA020 \| layer` |
| `0xB000-0xBFFF` | LT(layer, kc) | `0xB000 \| (layer << 8) \| usage_id` |

**LT のビットレイアウト:**

```
0xBLUU
  │││
  ││└─ UU: HID Usage ID（0x00-0xFF）
  │└── L:  レイヤー番号（0x0-0xF）
  └─── B:  LT プレフィックス
```

**例:** `LT(1, A)` = `0xB000 | (1 << 8) | 4` = `0xB104`

### keycode パッケージへの追加

```go
// 透過キー
const TRNS Keycode = 0x0001

// レイヤーアクション生成関数
func MO(layer uint8) Keycode  // Momentary
func TG(layer uint8) Keycode  // Toggle
func TT(layer uint8) Keycode  // Tap-Toggle
func LT(layer uint8, kc Keycode) Keycode  // Layer-Tap

// 判定関数
func (kc Keycode) IsMO() bool
func (kc Keycode) IsTG() bool
func (kc Keycode) IsTT() bool
func (kc Keycode) IsLT() bool
func (kc Keycode) IsLayerAction() bool  // MO|TG|TT|LT のいずれか
func (kc Keycode) IsTapAction() bool    // TT|LT のいずれか

// 情報抽出関数
func (kc Keycode) Layer() int           // レイヤー番号を返す
func (kc Keycode) TapKeycode() Keycode  // LT のタップ時キーコードを返す
```

---

## レイヤー解決（engine/layer.go）

### 定数

```go
const MaxLayers = 4
```

### Resolver 構造体

```go
// Resolver はアクティブレイヤーからキーコードを解決します。
type Resolver struct {
    keymap *Keymap
    active [MaxLayers]bool  // active[0] は常に true
}
```

### 解決アルゴリズム

最上位のアクティブレイヤーから下方向に走査し、最初の非 TRNS キーコードを返す。

```
Resolve(row, col):
  for layer = MaxLayers-1 downto 0:
    if !active[layer]: continue
    kc = keymap.Layers[layer][row][col]
    if kc != TRNS: return kc
  return None
```

### メソッド

| メソッド | 動作 |
|----------|------|
| `NewResolver(km *Keymap) *Resolver` | Resolver を生成（Layer 0 を active に初期化） |
| `Resolve(row, col int) Keycode` | キーコード解決 |
| `Activate(layer int)` | レイヤーを有効化（layer > 0 のみ） |
| `Deactivate(layer int)` | レイヤーを無効化（layer > 0 のみ） |
| `Toggle(layer int)` | レイヤーの ON/OFF を反転 |
| `IsActive(layer int) bool` | レイヤーの状態を返す |

---

## タップ検出（engine/tap.go）

### 定数

```go
const tapThreshold = 200  // ticks（= 200ms）
```

### 状態遷移

```
tapIdle ──[押下(LT/TT)]-→ tapPending ──[閾値超過]-→ tapHolding
  ▲                          │                         │
  └──────[リリース]──────────┘                         │
  └──────────────────────[リリース]─────────────────────┘

tapPending → リリース = タップ（LT: キーコード送信 / TT: レイヤートグル）
tapHolding → リリース = ホールド解除（レイヤー無効化）
```

### TapDetector 構造体

```go
type tapPhase uint8

const (
    tapIdle    tapPhase = iota
    tapPending
    tapHolding
)

// TapDetector はタップ/ホールド判定を管理します。
// 全フィールドは固定サイズ配列でヒープ割り当てを回避します。
type TapDetector struct {
    phase   [matrix.RowCount][matrix.ColCount]tapPhase
    counter [matrix.RowCount][matrix.ColCount]uint16
    kc      [matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

### メソッド

| メソッド | 動作 |
|----------|------|
| `Press(row, col int, kc Keycode)` | pending 状態に入る |
| `Release(row, col int) (wasPending bool, kc Keycode)` | リリース処理。タップなら `wasPending=true` |
| `Advance()` | 全 pending キーのカウンタをインクリメント |
| `CheckTimeout(row, col int) (timedOut bool, kc Keycode)` | 閾値超過チェック。超過時 holding に遷移 |
| `Phase(row, col int) tapPhase` | 指定位置の現在フェーズを返す |
| `Reset(row, col int)` | 状態をリセット |

### メモリ使用量

```
phase:   3 × 4 × 1 byte  = 12 bytes
counter: 3 × 4 × 2 bytes = 24 bytes
kc:      3 × 4 × 2 bytes = 24 bytes
合計: 60 bytes
```

---

## Keymap 構造体の変更

### Before（Phase 2）

```go
type Keymap struct {
    Layer0 [matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

### After（Phase 3）

```go
type Keymap struct {
    Layers [MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

### メモリ影響

- Before: 3 × 4 × 2 = 24 bytes
- After: 4 × 3 × 4 × 2 = 96 bytes（+72 bytes）

---

## Keyboard 構造体の変更

### 追加フィールド

```go
type Keyboard struct {
    scanner     *matrix.Scanner
    keymap      *Keymap
    productName string
    prevState   [matrix.RowCount][matrix.ColCount]bool
    activeKeys  [matrix.RowCount][matrix.ColCount]keycode.Keycode  // 追加: 押下中のキーコード/アクション
    resolver    *Resolver                                          // 追加: レイヤー解決
    tap         TapDetector                                        // 追加: タップ検出
}
```

### tick() の処理フロー

```
tick():
  1. tap.Advance()
     全 pending キーのカウンタをインクリメント

  2. タイムアウトチェック（全キーポジション走査）
     pending かつ閾値超過 → holding に遷移
     → レイヤーを Activate

  3. state, changed = scanner.Scan()
     changed == false なら return

  4. 変化キーの処理
     for each (row, col):
       if curr == prev: continue

       if 押下:
         kc = resolver.Resolve(row, col)
         handlePress(row, col, kc)

       if リリース:
         handleRelease(row, col)

  5. prevState = state
```

### handlePress(row, col, kc)

```
switch:
  MO(n):
    resolver.Activate(n)
    activeKeys[row][col] = kc

  TG(n):
    resolver.Toggle(n)
    activeKeys[row][col] = kc

  LT(n, tapKc) / TT(n):
    tap.Press(row, col, kc)
    // activeKeys はまだ設定しない

  通常キー / 修飾キー:
    HID Down 送信
    activeKeys[row][col] = kc

  None:
    何もしない
```

### handleRelease(row, col)

```
1. タップ判定チェック
   wasPending, kc = tap.Release(row, col)
   if wasPending:
     LT: タップキーコードを Down → Up（即時送信）
     TT: レイヤートグル
     activeKeys クリア、return

2. activeKeys[row][col] に基づく処理
   switch:
     MO(n): resolver.Deactivate(n)
     LT(n, _) / TT(n): resolver.Deactivate(n)  // ホールド中だった
     TG(n): 何もしない（トグルは押下時に処理済み）
     通常キー / 修飾キー: HID Up 送信

3. activeKeys[row][col] = None
```

---

## 変更するコンポーネント一覧

### 新規ファイル

| ファイル | ビルドタグ | 役割 |
|----------|:----------:|------|
| `engine/layer.go` | なし | Resolver 構造体・レイヤー解決ロジック |
| `engine/layer_test.go` | なし | Resolver のユニットテスト |
| `engine/tap.go` | なし | TapDetector 構造体・タップ/ホールド判定 |
| `engine/tap_test.go` | なし | TapDetector のユニットテスト |

### 変更ファイル

| ファイル | 変更内容 |
|----------|----------|
| `keycode/keycode.go` | TRNS 定数、MO/TG/TT/LT 生成関数、判定・抽出関数 |
| `keycode/keycode_test.go` | 新規キーコードのテスト追加 |
| `engine/keymap.go` | `Layer0` → `Layers [MaxLayers][...]` に変更 |
| `engine/keyboard.go` | Resolver・TapDetector 統合、tick() 書き換え |
| `engine/config.go` | 変更なし（Config 構造体は互換性維持） |
| `examples/zero-kb02/keymap.go` | 新 Keymap 構造体に対応 |

### ドキュメント更新

| ファイル | 変更内容 |
|----------|----------|
| `docs/packages/engine.md` | Resolver・TapDetector・Keymap 変更を反映 |
| `docs/packages/keycode.md` | 新エンコーディング範囲を追加 |
| `docs/functional-design.md` | レイヤー解決のデータフロー追加 |
| `docs/repository-structure.md` | 新ファイルの追記 |
| `ROADMAP.md` | Phase 3 の完了チェック |

---

## 影響範囲の分析

### 後方互換性

| 項目 | 互換性 | 備考 |
|------|:------:|------|
| `Config` 構造体 | ✅ 互換 | フィールド変更なし |
| `engine.New()` / `Run()` | ✅ 互換 | シグネチャ変更なし |
| `Keymap` 構造体 | ❌ 破壊 | `Layer0` → `Layers` に変更 |
| `keycode.Keycode` 型 | ✅ 互換 | 既存値のエンコーディングは不変 |
| `matrix` パッケージ | ✅ 無変更 | 影響なし |

`Keymap` 構造体の変更は破壊的だが、フレームワーク利用側（`examples/zero-kb02`）の修正は
フィールド名変更のみで軽微。

### パフォーマンス影響

| 処理 | オーバーヘッド |
|------|----------------|
| レイヤー解決 | 最大 4 レイヤー × 1 比較 = 4 回の条件分岐。1ms 周期に影響なし |
| タップカウンタ | 12 キー × 1 インクリメント = 12 回の加算。1ms 周期に影響なし |
| タイムアウトチェック | 12 キー × 1 比較 = 12 回の条件分岐。1ms 周期に影響なし |
| RAM 増加 | Keymap +72B, TapDetector +60B, Resolver +4B, activeKeys +24B = +160B |

---

## feature ブランチ構成

```
main
└── phase3
    ├── feature/layer-keycode     ← keycode 拡張
    ├── feature/layer-resolve     ← Resolver・Keymap 変更
    ├── feature/mo-tg             ← MO/TG の engine 統合
    └── feature/tap-detect        ← TapDetector・LT/TT の engine 統合
```

各ブランチの詳細スコープは `tasklist.md` で定義する。

---

## examples/zero-kb02 の更新例

### Before

```go
var defaultKeymap = &engine.Keymap{
    Layer0: [matrix.RowCount][matrix.ColCount]keycode.Keycode{
        {keycode.Q, keycode.W, keycode.E, keycode.R},
        {keycode.A, keycode.S, keycode.D, keycode.F},
        {keycode.Z, keycode.X, keycode.C, keycode.V},
    },
}
```

### After

```go
var defaultKeymap = &engine.Keymap{
    Layers: [engine.MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode{
        // Layer 0: ベースレイヤー
        {
            {keycode.Q, keycode.W, keycode.E, keycode.R},
            {keycode.A, keycode.S, keycode.D, keycode.F},
            {keycode.Z, keycode.X, keycode.C, keycode.MO(1)},
        },
        // Layer 1: 数字・記号レイヤー
        {
            {keycode.Num1, keycode.Num2, keycode.Num3, keycode.Num4},
            {keycode.Num5, keycode.Num6, keycode.Num7, keycode.Num8},
            {keycode.Num9, keycode.Num0, keycode.TRNS, keycode.TRNS},
        },
        // Layer 2-3: 未使用（ゼロ値 = None）
    },
}
```
