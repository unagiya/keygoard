# engine パッケージ仕様

## 概要

キーボードエンジンのコアパッケージです。マトリクススキャン・キーマップ参照・HID 送信を統合し、`Run()` 一つでキーボードを起動できるフレームワーク API を提供します。

## ファイル構成

| ファイル        | machine 依存 | 役割                             |
|-----------------|:------------:|----------------------------------|
| `config.go`     | なし         | `Config` 構造体定義              |
| `keyboard.go`   | あり         | エンジン本体（`New` / `Run`）    |
| `keymap.go`     | なし         | `Keymap` 構造体定義              |
| `errors.go`     | なし         | エラー変数定義                   |

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
type Keymap struct {
    Layer0 [matrix.RowCount][matrix.ColCount]keycode.Keycode
}
```

Phase 1 では単一レイヤーのみ対応。

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
3. 変化があったキーについて `hidkb.Keyboard.Down()` / `Up()` で HID レポートを送信
4. 前回状態を更新

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

`keyboard.go` は `//go:build tinygo` タグ付きのため、標準 Go ではビルドできない。
テスト可能な検証は以下の通り。

```bash
# keymap / errors は machine 非依存だがテストファイル未作成
# 関連パッケージのテストで間接的に検証
go test ./keycode/... ./matrix/...
```
