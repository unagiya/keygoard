# Phase 5: パッケージ構造リファクタリング — タスクリスト

## ブランチ運用

- `phase5` ブランチを `main` から作成し、全作業をこのブランチで実施
- 完了後 `phase5` → `main` へ PR を作成してマージ

## タスク一覧

### 1. internal ディレクトリの作成とファイル移動

- [x] 1-1. `internal/matrix/` を作成し `matrix/` の全ファイルを移動
- [x] 1-2. `internal/layer/` を作成し `engine/layer.go` → `resolver.go`、`engine/layer_test.go` → `resolver_test.go` として移動
- [x] 1-3. `internal/tap/` を作成し `engine/tap.go` → `detector.go`、`engine/tap_test.go` → `detector_test.go` として移動
- [x] 1-4. `internal/peripheral/encoder/` を作成し `engine/peripheral/encoder/` の全ファイルを移動
- [x] 1-5. `internal/peripheral/led/` を作成し `engine/peripheral/led/` の全ファイルを移動
- [x] 1-6. `internal/peripheral/oled/` を作成し `engine/peripheral/oled/` の全ファイルを移動
- [x] 1-7. 旧ディレクトリ（`matrix/`、`engine/peripheral/`）を削除

### 2. internal パッケージのリネーム・修正

- [x] 2-1. `internal/layer`: package 名を `layer` に変更、`Resolver` が `*[MaxLayers][RowCount][ColCount]keycode.Keycode` を受け取るよう変更
- [x] 2-2. `internal/layer`: テストの import パスとパッケージ名を修正
- [x] 2-3. `internal/tap`: package 名を `tap` に変更、`TapDetector` → `Detector` にリネーム
- [x] 2-4. `internal/tap`: テストの import パスとパッケージ名を修正
- [x] 2-5. `internal/matrix`: import パスのみ変更（パッケージ名・ロジックは維持）
- [x] 2-6. `internal/peripheral/*`: import パスを `internal/` 配下に変更、Config 構造体を削除（engine 側で定義）

### 3. engine パッケージの Facade 化

- [x] 3-1. `engine/types.go` を新規作成（`Pin` / `I2CBus` / `Rotation` / `LEDType` 型と定数）
- [x] 3-2. `engine/config.go` を書き直し（`EncoderConfig` / `LEDConfig` / `OLEDConfig` を追加、`Scanner` フィールドを `ColPins`/`RowPins` に置換）
- [x] 3-3. `engine/keymap.go` を修正（`internal/matrix` から import、`RowCount`/`ColCount` を再エクスポート）
- [x] 3-4. `engine/keyboard.go` を修正（`New()` 内でマトリクス・ペリフェラルの内部生成、Pin 変換ヘルパー追加、import パスを `internal/*` に変更）

### 4. examples/zero-kb02 の書き直し

- [x] 4-1. `main.go` を新 API で書き直し（`engine` と `keycode` のみ import）
- [x] 4-2. `config.go` を削除
- [x] 4-3. `keymap.go` を修正（`matrix` import を除去、`engine.RowCount`/`engine.ColCount` を使用）

### 5. Makefile・ドキュメント更新

- [x] 5-1. `Makefile` のテスト対象パスを `./keycode/... ./internal/...` に変更
- [x] 5-2. `docs/repository-structure.md` を新しいディレクトリ構造に更新
- [x] 5-3. `docs/functional-design.md` のシステム構成図・API 設計を更新
- [x] 5-4. `docs/packages/` 配下の仕様書を新パス・パッケージ名に合わせて更新
- [x] 5-5. `docs/packages/layer.md` / `tap.md` を新規作成
- [x] 5-6. `docs/architecture.md` / `docs/glossary.md` のテストパス・命名セクションを更新

### 6. 検証

- [x] 6-1. `make test` が通ること（tinygo build + go test）
- [x] 6-2. `make lint` が通ること
- [x] 6-3. 実機書き込み（`make flash`）後、全機能が Phase 4 と同等に動作すること
      - キー入力（Layer 0）
      - レイヤー切り替え（MO / LT）
      - エンコーダー回転
      - LED のレイヤー色変更
      - OLED のレイヤー表示

## 完了条件

- [x] 利用者の `main.go` が `engine` と `keycode` の 2 パッケージのみで記述できている
- [x] 利用者の `main.go` が `machine` パッケージを直接 import していない
- [x] `internal/` 配下のパッケージが外部から import できない
- [x] `make test` が通る
- [x] 実機で Phase 4 同等の動作が確認できている
