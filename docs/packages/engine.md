# engine パッケージ仕様

## 概要

キーボードエンジンのコアパッケージです。マトリクススキャン・レイヤー解決・HID 送信を統合し、`Run()` 一つでキーボードを起動できるフレームワーク API を提供します。

## ファイル構成

| ファイル          | machine 依存 | 役割                             |
|-------------------|:------------:|----------------------------------|
| `config.go`       | あり         | `Config` 構造体定義              |
| `keyboard.go`     | あり         | エンジン本体（`New` / `Run`）    |
| `keymap.go`       | なし         | `Keymap` 構造体定義              |
| `layer.go`        | なし         | `Resolver` — レイヤー解決ロジック |
| `layer_test.go`   | なし         | Resolver のユニットテスト        |
| `tap.go`          | なし         | `TapDetector` — タップ/ホールド判定 |
| `tap_test.go`     | なし         | TapDetector のユニットテスト     |
| `errors.go`       | なし         | エラー変数定義                   |

## API

### Config

```go
type Config struct {
    Scanner     *matrix.Scanner // マトリクススキャナー
    Keymap      *Keymap         // キーマップ
    ProductName string          // USB Product Name（空の場合 "keygoard"）
}
```

- `ProductName` は ASCII のみ、最大 126 文字

### Keyboard

```go
// Config からキーボードエンジンを生成
kb := engine.New(&engine.Config{
    Scanner:     matrix.New(cols, rows),
    Keymap:      keymap,
    ProductName: "my-keyboard",
})

// キーボードを起動（この関数は戻らない）
kb.Run()
```

### Keymap

```go
const MaxLayers = 4

type Keymap struct {
    Layers [MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

- `Layers[0]` がベースレイヤー（常に有効）
- 番号が大きいほど優先度が高い
- 未使用レイヤーはゼロ値（全キー `None`）のまま

### Resolver

```go
// Resolver を生成（レイヤー 0 は常に有効）
r := NewResolver(km)

// キーコードを解決（最上位アクティブレイヤーから走査）
kc := r.Resolve(row, col)

// レイヤー操作
r.Activate(layer)
r.Deactivate(layer)
r.Toggle(layer)
r.IsActive(layer)
```

- 最上位のアクティブレイヤーから下方向に走査
- `TRNS`（透過キー）は下位レイヤーにフォールバック
- すべて `TRNS` の場合は `None` を返す
- レイヤー 0 は Deactivate / Toggle できない

### TapDetector

```go
const tapThreshold = 200 // ticks（= 200ms）

type TapDetector struct { /* 固定サイズ配列、ヒープ割り当てなし */ }

// キー押下時（LT/TT キーコード）に pending 状態に入る
td.Press(row, col, kc)

// キーリリース時。pending なら wasPending=true（タップ判定）
wasPending, kc := td.Release(row, col)

// 全 pending キーのカウンタをインクリメント（毎スキャンサイクル呼び出し）
td.Advance()

// 閾値超過チェック。超過時 holding に遷移して timedOut=true
timedOut, kc := td.CheckTimeout(row, col)

// 指定位置の状態をリセット
td.Reset(row, col)
```

**状態遷移:**

```
tapIdle ──[Press(LT/TT)]──→ tapPending ──[閾値超過]──→ tapHolding
  ▲                              │                         │
  └──────[Release=タップ]────────┘                         │
  └──────────────────────[Release=ホールド解除]─────────────┘
```

- `tapPending` でリリース → タップ（LT: キーコード送信 / TT: レイヤートグル）
- `tapHolding` でリリース → ホールド解除（レイヤー無効化）

## Run() の内部処理

```
1. USB Product Name を設定（usb.Product に代入）
2. scanner.Init()（GPIO 初期化）
3. USB エニュメレーション完了まで待機（500ms）
4. 無限ループ:
   a. tick()（1 スキャンサイクル実行）
   b. time.Sleep(1ms)
```

### tick() の動作

1. `tap.Advance()` で全 pending キーのカウンタをインクリメント
2. 全キーポジションのタイムアウトチェック（pending → holding 遷移時にレイヤー有効化）
3. `scanner.Scan()` でデバウンス済みキー状態を取得
4. 変化がなければ即リターン
5. 変化があったキーについて:
   - 押下: `resolver.Resolve()` でキーコードを解決 → `handlePress()`
   - リリース: `handleRelease()`（タップ判定 → activeKeys に基づく解除）
6. 前回状態を更新

## USB Product Name

TinyGo の `machine/usb.Product` パッケージ変数に値を代入することで設定する。

- `Run()` の冒頭で `usb.Product` に代入
- ホストが USB String Descriptor を要求した時点で遅延評価される
- `Config.ProductName` が空の場合、デフォルト値 `"keygoard"` を使用

## エラー

| 変数             | 意味                           |
|------------------|--------------------------------|
| `ErrKeyOverflow` | 6KRO の同時押し上限を超えた場合 |

## テスト

```bash
# Resolver・TapDetector のユニットテスト（machine 非依存）
go test ./engine/...
```
