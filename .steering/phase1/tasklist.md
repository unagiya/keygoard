# Phase 1 タスクリスト

## 完了条件

- [ ] 全タスク完了
- [x] `tinygo build -target=waveshare-rp2040-zero` 成功
- [x] `go test ./keycode/... ./matrix/...` 全パス
- [x] macOS でキーボードデバイスとして認識される
- [x] 実際にキー入力が届くことを実機確認
- [ ] `v0.1.0` タグを `main` に打つ

---

## 実装メモ

> 設計時は `feature/matrix-scan` / `feature/debounce` / `feature/hid-basic` / `feature/keymap-layer0` の
> 4 ブランチ分割を想定していたが、実際には `feature/matrix-scan` ブランチ 1 本で
> Phase 1 の全コードを実装した。
> また `hid/` を独立パッケージとする設計を変更し、HID 送信は `engine/keyboard.go` に統合した。

---

## feature/matrix-scan（実装済み）

**ブランチ:** `feature/matrix-scan`
**目標:** zero-kb02 のピン配置で 3×4 マトリクスをスキャンできる

- [x] `matrix/const.go` を作成（RowCount / ColCount 定数）
- [x] `matrix/matrix.go` を作成（`//go:build tinygo` タグ付き）
  - [x] `Init()`: Col/Row ピン初期化
  - [x] `Scan()`: COL2ROW スキャン実装、デバウンス呼び出し
- [x] `docs/packages/matrix.md` を作成（スキャン仕様を記述）

---

## feature/debounce（実装済み）

**ブランチ:** `feature/matrix-scan`（統合実装）
**目標:** チャタリング除去のデバウンスロジックをユニットテスト済みで実装する

- [x] `matrix/debounce.go` を作成（`machine` 非依存）
  - [x] `Debouncer` 構造体定義
  - [x] `Update(raw [RowCount][ColCount]bool)` 実装
  - [x] `debounceThr` 定数定義（5 サイクル）
- [x] `matrix/debounce_test.go` を作成
  - [x] 閾値未満のパルスでは状態変化しないテスト
  - [x] 閾値以上継続で状態確定するテスト
  - [x] キーリリースのデバウンステスト
- [x] `go test ./matrix/...` 全パス
- [x] `docs/packages/matrix.md` にデバウンス仕様を追記

---

## feature/hid-basic（実装済み・設計変更あり）

**ブランチ:** `feature/matrix-scan`（統合実装）
**目標:** TinyGo 標準 HID で USB キーボードとして OS に認識される

> 設計変更: 独立した `hid/` パッケージは作成せず、`engine/keyboard.go` に統合した。

- [x] `engine/keyboard.go` に HID 送信を実装（`//go:build tinygo` タグ付き）
  - [x] `Init()`: USB エニュメレーション待機（500ms）
  - [x] `Tick()`: キー変化時に `hidkb.Keyboard.Down/Up` 送信
- [x] 実機書き込みで macOS にキーボードデバイスとして表示されることを確認

---

## feature/keymap-layer0（実装済み）

**ブランチ:** `feature/matrix-scan`（統合実装）
**目標:** キーを押したら対応する文字が入力される

- [x] `keycode/keycode.go` を作成（HID Usage ID キーコード定数）
- [x] `keycode/keycode_test.go` を作成（定数値の検証）
- [x] `engine/keymap.go` を作成（`Keymap` 構造体・Layer0 フィールド）
- [x] `engine/errors.go` を作成（`ErrKeyOverflow` 定義）
- [x] `engine/keyboard.go` にメインループを実装
  - [x] matrix.Scan → debounce（Scanner 内部） → キーコード収集 → hid.Send
- [x] `examples/zero-kb02/main.go` を作成（`engine.New` + `kb.Tick()` ループ）
- [x] `examples/zero-kb02/config.go` を作成（ピン定義）
- [x] `examples/zero-kb02/keymap.go` を作成（デフォルトキーマップ Q〜V）
- [x] `docs/packages/keycode.md` を作成
- [x] `go test ./keycode/... ./matrix/...` 全パス
- [x] `tinygo build -target=waveshare-rp2040-zero` 成功
- [x] 実機書き込みで文字入力できることを確認

---

## Phase 1 完了処理

- [ ] `phase1` → `main` へ PR 作成・マージ
- [ ] `main` に `v0.1.0` タグを打つ
- [x] ROADMAP.md の Phase 1 チェックリストを全て完了にする
- [ ] `.steering/phase1/tasklist.md` の完了条件を全てチェック
