# Phase 2 要求仕様

## 概要

Phase 1 で実装した最小 HID キーボードの API をブラッシュアップし、フレームワークとしての使いやすさを向上させる。

---

## ユーザーストーリー

1. **メインループを意識せずにファームウェアを書ける**
   - ファームウェア開発者は `engine.New()` → `engine.Run()` の 2 行でキーボードを起動でき、`for` ループや `time.Sleep` を自分で書く必要がない。

2. **キーボードの名前を自分でつけられる**
   - ファームウェア開発者は USB Product Name を設定でき、OS のデバイス一覧に自分が付けた名前（例: "zero-kb02"）が表示される。

---

## 受け入れ条件

| # | 条件 |
|---|---|
| 1 | `examples/zero-kb02/main.go` から `for` ループ・`time.Sleep` が消え、`engine.Run()` のみで動作する |
| 2 | `engine.Run()` 内部でスキャンループが実行される |
| 3 | USB Product Name を設定する API が提供されている |
| 4 | macOS のシステム情報に設定した Product Name が表示される |
| 5 | `go test ./keycode/... ./matrix/...` が引き続き全パスする |
| 6 | `tinygo build -target=waveshare-rp2040-zero` が成功する |
| 7 | 実機でキー入力が引き続き動作する |

---

## 制約

### スコープ内（Phase 2 で実装する）

- `engine.Run()` メソッドの追加（`Init()` + メインループを内包）
- USB Product Name の設定 API
- `examples/zero-kb02/` の更新（新 API に対応）

### スコープ外（Phase 2 では実装しない）

- スキャン間隔の外部設定化（デフォルト値のままで良い）
- USB Vendor Name / Serial Number 等のカスタマイズ（Product Name のみ）
- Phase 1 のスコープ外に挙げた全機能

### 技術制約

- Phase 1 と同じ TinyGo 制約を全て遵守すること
- `engine.Run()` は戻らない（無限ループ）ため、呼び出し前に全ての設定を完了する設計とする
