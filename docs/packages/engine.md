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

1. `scanner.Scan()` でデバウンス済みキー状態を取得
2. 変化がなければ即リターン
3. 変化があったキーについて `resolver.Resolve()` でキーコードを解決
4. `hidkb.Keyboard.Down()` / `Up()` で HID レポートを送信
5. 前回状態を更新

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
# Resolver のユニットテスト（machine 非依存）
go test ./engine/...
```
