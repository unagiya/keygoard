# Phase 2 タスクリスト

## 完了条件

- [x] 全タスク完了
- [x] `go test ./keycode/... ./matrix/...` 全パス
- [x] `tinygo build -target=waveshare-rp2040-zero` 成功
- [x] macOS でキーボードデバイスとして認識される
- [x] 設定した Product Name が OS に表示される
- [x] 実際にキー入力が届くことを実機確認
- [ ] `v0.2.0` タグを `main` に打つ

---

## Phase ブランチ作成

- [x] `main` を最新に更新（`git checkout main && git pull`）
- [x] `phase2` ブランチを作成（`git checkout -b phase2`）
- [x] `phase2` をリモートにプッシュ（`git push -u origin phase2`）

---

## feature/engine-run

**ブランチ:** `feature/engine-run`（`phase2` から分岐 → `phase2` へ PR）
**目標:** `engine.Run()` でメインループをフレームワーク内部に隠蔽する

> **振り返り:** 本ブランチで `feature/usb-product-name` のスコープ（`ProductName` フィールド・
> `usb.Product` 設定）も実装してしまった。ブランチスコープルールの策定前だったため、
> 今回はこのまま進めるが、今後は tasklist の該当セクションのみを実装すること。

### 1. ブランチ作成

- [x] `phase2` から分岐（`git checkout -b feature/engine-run phase2`）

### 2. 実装

- [x] `engine/config.go` を作成（`Config` 構造体定義）
  - [x] `Scanner` / `Keymap` フィールド
  - [x] `ProductName` フィールド（※本来は feature/usb-product-name のスコープ）
- [x] `engine/keyboard.go` を変更
  - [x] `New(cfg *Config) *Keyboard` にシグネチャ変更
  - [x] `Run()` メソッドを追加（setup + 無限ループ）
  - [x] `Init()` を非公開化（`setup()` にリネーム）
  - [x] `Tick()` を非公開化（`tick()` にリネーム）
  - [x] `Run()` 冒頭で `usb.Product` を設定（※本来は feature/usb-product-name のスコープ）
- [x] `examples/zero-kb02/main.go` を更新
  - [x] `for` ループ・`time.Sleep` を削除
  - [x] `engine.New(&engine.Config{...})` + `kb.Run()` に変更
  - [x] `ProductName: "zero-kb02"` を設定（※本来は feature/usb-product-name のスコープ）

### 3. 検証

- [x] `go test ./keycode/... ./matrix/...` 全パス
- [x] `tinygo build -target=waveshare-rp2040-zero` 成功
- [x] 実機書き込み・macOS にキーボードデバイスとして認識される
- [x] キー入力が引き続き動作する

### 4. コミット・PR

- [x] コミット（`feat: engine.Run() を追加しメインループをフレームワークに内包`）
- [x] リモートにプッシュ（`git push -u origin feature/engine-run`）
- [x] `feature/engine-run` → `phase2` へ PR 作成（#3）
- [x] PR マージ

---

## feature/usb-product-name

**ステータス: スキップ（feature/engine-run に統合済み）**

`feature/engine-run` で以下の全項目が実装済みのため、本ブランチの作業は不要。

- [x] `engine/config.go` の `Config` に `ProductName` フィールドを追加
- [x] `engine/keyboard.go` の `Run()` 冒頭で `usb.Product` を設定
  - [x] `Config.ProductName` が空の場合はデフォルト値 `"keygoard"` を使用
- [x] `examples/zero-kb02/main.go` の `Config` に `ProductName: "zero-kb02"` を追加

---

## ドキュメント・実機確認

**ブランチ:** `phase2` 上で直接、または別途 feature ブランチ

### 1. ドキュメント

- [x] `docs/packages/engine.md` を新規作成（engine パッケージ仕様書）
- [x] コミット（`docs: engine パッケージ仕様書を追加`）

### 2. ビルド・実機確認

- [x] `tinygo build -target=waveshare-rp2040-zero` 成功
- [x] 実機書き込みで macOS にキーボードデバイスとして認識される
- [x] 設定した Product Name（"zero-kb02"）が OS に表示される
- [x] キー入力が引き続き動作する

---

## Phase 2 完了処理

### 1. ドキュメント更新

- [x] ROADMAP.md の Phase 2 チェックリストを全て完了にする
- [x] `.steering/phase2/tasklist.md` の完了条件を全てチェック
- [ ] ドキュメント更新をコミット

### 2. マージ・タグ

- [ ] `phase2` → `main` へ PR 作成・マージ
- [ ] `main` に `v0.2.0` タグを打つ
