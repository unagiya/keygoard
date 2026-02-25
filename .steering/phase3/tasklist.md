# Phase 3: レイヤーシステム — タスクリスト

---

## feature/layer-keycode — keycode パッケージ拡張

### 実装

- [x] `keycode/keycode.go` に `TRNS` 定数を追加（`0x0001`）
- [x] `keycode/keycode.go` にレイヤーアクション生成関数を追加
  - `MO(layer uint8) Keycode`
  - `TG(layer uint8) Keycode`
  - `TT(layer uint8) Keycode`
  - `LT(layer uint8, kc Keycode) Keycode`
- [x] `keycode/keycode.go` に判定メソッドを追加
  - `IsMO() bool`
  - `IsTG() bool`
  - `IsTT() bool`
  - `IsLT() bool`
  - `IsLayerAction() bool`
  - `IsTapAction() bool`
- [x] `keycode/keycode.go` に情報抽出メソッドを追加
  - `Layer() int`
  - `TapKeycode() Keycode`

### テスト

- [x] `keycode/keycode_test.go` に TRNS のテストを追加
- [x] `keycode/keycode_test.go` に MO/TG/TT/LT 生成・判定・抽出のテストを追加
- [x] `go test ./keycode/...` 全パス

### ドキュメント

- [x] `docs/packages/keycode.md` を更新（新エンコーディング範囲を追記）

### 検証

- [x] `make lint` パス
- [x] `make build` パス（TinyGo ビルド）
- [x] 実機書き込み・動作確認（Phase 2 同等の動作 — 回帰テスト）

### PR・マージ

- [x] `feature/layer-keycode` → `phase3` の PR 作成
- [x] PR マージ

---

## feature/layer-resolve — Resolver・Keymap 変更

### 実装

- [x] `engine/keymap.go` の `Keymap` 構造体を変更
  - `Layer0` フィールドを `Layers [MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode` に変更
- [x] `engine/layer.go` を新規作成（ビルドタグなし）
  - `MaxLayers` 定数
  - `Resolver` 構造体（`keymap *Keymap`, `active [MaxLayers]bool`）
  - `NewResolver(km *Keymap) *Resolver`
  - `Resolve(row, col int) keycode.Keycode`
  - `Activate(layer int)`
  - `Deactivate(layer int)`
  - `Toggle(layer int)`
  - `IsActive(layer int) bool`
- [x] `engine/keyboard.go` の `Keyboard` 構造体に `resolver` フィールドを追加
- [x] `engine/keyboard.go` の `New()` で `NewResolver` を呼び出す
- [x] `engine/keyboard.go` の `tick()` でキーコード解決を `resolver.Resolve()` 経由に変更
- [x] `examples/zero-kb02/keymap.go` を新 Keymap 構造体に対応

### テスト

- [x] `engine/layer_test.go` を新規作成
  - 単一レイヤー解決テスト
  - 複数レイヤー解決テスト（上位レイヤー優先）
  - KC_TRNS 透過テスト（下位レイヤーにフォールバック）
  - 全レイヤー TRNS の場合 None を返すテスト
  - Activate / Deactivate / Toggle テスト
  - Layer 0 は Deactivate できないテスト
- [x] `go test ./engine/...` 全パス
- [x] `go test ./keycode/...` 全パス（既存テスト維持）

### ドキュメント

- [x] `docs/packages/engine.md` を更新（Keymap・Resolver を反映）

### 検証

- [x] `make lint` パス
- [x] `make build` パス
- [x] 実機書き込み・動作確認（レイヤー 0 のみで Phase 2 同等の動作 — 回帰テスト）

### PR・マージ

- [x] `feature/layer-resolve` → `phase3` の PR 作成
- [x] PR マージ

---

## feature/mo-tg — MO/TG の engine 統合

### 実装

- [x] `engine/keyboard.go` に `activeKeys [RowCount][ColCount]keycode.Keycode` フィールドを追加
- [x] `engine/keyboard.go` に `handlePress(row, col int, kc keycode.Keycode)` メソッドを追加
  - MO(n): `resolver.Activate(n)` + activeKeys 記録
  - TG(n): `resolver.Toggle(n)` + activeKeys 記録
  - 通常キー/修飾キー: HID Down + activeKeys 記録
  - None: 何もしない
- [x] `engine/keyboard.go` に `handleRelease(row, col int)` メソッドを追加
  - MO(n): `resolver.Deactivate(n)` + activeKeys クリア
  - TG(n): 何もしない + activeKeys クリア
  - 通常キー/修飾キー: HID Up + activeKeys クリア
- [x] `engine/keyboard.go` の `tick()` を handlePress/handleRelease 呼び出しに書き換え
- [x] `examples/zero-kb02/keymap.go` に MO(1) を含むキーマップを設定

### テスト

- [x] `go test ./engine/...` 全パス
- [x] `go test ./keycode/...` 全パス

### 検証

- [x] `make lint` パス
- [x] `make build` パス
- [x] 実機書き込み
- [x] MO キーの動作確認
  - MO(1) ホールド中にレイヤー 1 のキーコードが送信される
  - MO(1) リリースでレイヤー 0 に戻る
- [x] TG キーの動作確認
  - TG(1) でレイヤー 1 が有効になる
  - 再度 TG(1) でレイヤー 0 に戻る
- [x] KC_TRNS の動作確認
  - レイヤー 1 の TRNS キーがレイヤー 0 のキーコードを送信する

### PR・マージ

- [x] `feature/mo-tg` → `phase3` の PR 作成
- [x] PR マージ

---

## feature/tap-detect — TapDetector・LT/TT の engine 統合

### 実装

- [x] `engine/tap.go` を新規作成（ビルドタグなし）
  - `tapThreshold` 定数（200）
  - `tapPhase` 型（`tapIdle` / `tapPending` / `tapHolding`）
  - `TapDetector` 構造体
  - `Press(row, col int, kc keycode.Keycode)`
  - `Release(row, col int) (wasPending bool, kc keycode.Keycode)`
  - `Advance()`
  - `CheckTimeout(row, col int) (timedOut bool, kc keycode.Keycode)`
  - `Phase(row, col int) tapPhase`
  - `Reset(row, col int)`
- [x] `engine/keyboard.go` に `tap TapDetector` フィールドを追加
- [x] `engine/keyboard.go` の `tick()` に tap.Advance() とタイムアウトチェックを追加
- [x] `handlePress` に LT/TT 分岐を追加
  - `tap.Press(row, col, kc)` を呼び出す
- [x] `handleRelease` にタップ判定を追加
  - LT タップ: タップキーコードを Down → Up 送信
  - TT タップ: レイヤートグル
  - LT/TT ホールド解除: `resolver.Deactivate(n)`
- [x] `examples/zero-kb02/keymap.go` に LT を含むキーマップを設定（任意）

### テスト

- [x] `engine/tap_test.go` を新規作成
  - Press → Advance (< 閾値) → Release = タップ判定テスト
  - Press → Advance (>= 閾値) → CheckTimeout = ホールド判定テスト
  - Press → Release (即座) = タップ判定テスト
  - 複数キー同時の独立動作テスト
  - Reset テスト
- [x] `go test ./engine/...` 全パス
- [x] `go test ./keycode/...` 全パス

### ドキュメント

- [x] `docs/packages/engine.md` を更新（TapDetector を反映）

### 検証

- [x] `make lint` パス
- [x] `make build` パス
- [x] 実機書き込み
- [x] LT キーの動作確認
  - 短押し（< 200ms）で通常キーが入力される
  - 長押し（>= 200ms）でレイヤーが有効になる
  - 長押し中に他キーを押すとそのレイヤーのキーコードが送信される
  - リリースでレイヤーが無効に戻る
- [ ] TT キーの動作確認
  - 短押しでレイヤーがトグルされる
  - 長押しで MO として動作する

### PR・マージ

- [ ] `feature/tap-detect` → `phase3` の PR 作成
- [ ] PR マージ

---

## Phase 3 完了作業

- [ ] `docs/functional-design.md` を更新（レイヤー解決のデータフロー追加）
- [ ] `docs/repository-structure.md` を更新（新ファイル追記）
- [ ] `ROADMAP.md` の Phase 3 チェックボックスを更新
- [ ] `phase3` → `main` の PR 作成・マージ
- [ ] `v0.3.0` タグを作成
