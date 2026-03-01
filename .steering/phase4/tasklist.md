# Phase 4: 周辺機器 — タスクリスト

---

## feature/peripheral-interface — Peripheral インターフェース・engine 統合

### 実装

- [x]`engine/peripheral.go` を新規作成（ビルドタグなし）
  - `Peripheral` インターフェース（`Init()`, `Tick()`, `OnLayerChange()`）
  - `MaxPeripherals` 定数（4）
- [x]`engine/config.go` に `Peripherals []Peripheral` フィールドを追加
- [x]`engine/keyboard.go` の `Keyboard` 構造体に追加
  - `peripherals [MaxPeripherals]Peripheral`
  - `peripheralCount int`
  - `prevTopLayer int`
- [x]`engine/keyboard.go` の `New()` で Peripherals をコピー
- [x]`engine/keyboard.go` に `topLayer()` メソッドを追加
- [x]`engine/keyboard.go` に `notifyLayerChange()` メソッドを追加
- [x]`engine/keyboard.go` の `setup()` に周辺機器初期化・初期レイヤー通知を追加
- [x]`engine/keyboard.go` の `tick()` に周辺機器 Tick 呼び出し・レイヤー変更検出を追加

### テスト

- [x]`make test` 全パス（既存テストの回帰確認）

### 検証

- [x]`make lint` パス
- [x]`make build` パス
- [x]実機書き込み・動作確認（周辺機器なしで Phase 3 同等の動作 — 回帰テスト）

### PR・マージ

- [x]`feature/peripheral-interface` → `phase4` の PR 作成
- [x]PR マージ

---

## feature/encoder — ロータリーエンコーダー

### 実装

- [x]`engine/peripheral/encoder/quadrature.go` を新規作成（ビルドタグなし）
  - `Direction` 型（`DirNone`, `DirCW`, `DirCCW`）
  - `Decoder` 構造体（`prev uint8`）
  - `Update(a, b bool) Direction`（状態遷移テーブルによる O(1) 判定）
- [x]`engine/peripheral/encoder/encoder.go` を新規作成（`//go:build tinygo`）
  - `Config` 構造体（`PinA`, `PinB`, `KeyCW`, `KeyCCW`）
  - `Encoder` 構造体
  - `New(cfg *Config) *Encoder`
  - `Init()`（GPIO ピンをプルアップ入力に設定）
  - `Tick() keycode.Keycode`（A/B 読み取り → Decoder 判定 → キーコード返却）
  - `OnLayerChange(layer int)`（no-op）
- [x]`examples/zero-kb02/config.go` にエンコーダーピン定義を追加
- [x]`examples/zero-kb02/main.go` にエンコーダーの初期化・登録を追加

### テスト

- [x]`engine/peripheral/encoder/quadrature_test.go` を新規作成
  - CW 方向の状態遷移テスト（00→01→11→10→00）
  - CCW 方向の状態遷移テスト（00→10→11→01→00）
  - 状態変化なしで DirNone を返すテスト
  - 不正遷移（スキップ）で DirNone を返すテスト
- [x]`make test` 全パス

### ドキュメント

- [x]`docs/packages/encoder.md` を新規作成

### 検証

- [x]`make lint` パス
- [x]`make build` パス
- [x]実機書き込み
- [x]エンコーダー CW 回転で UpArrow が送信されることを確認
- [x]エンコーダー CCW 回転で DownArrow が送信されることを確認
- [x]既存キーボード機能に影響がないことを確認

### PR・マージ

- [x]`feature/encoder` → `phase4` の PR 作成
- [x]PR マージ

---

## feature/led — RGB LED（WS2812）

### 実装

- [x]`go.mod` / `examples/zero-kb02/go.mod` に `tinygo-org/drivers` 依存を追加
- [x]`engine/peripheral/led/led.go` を新規作成（`//go:build tinygo`）
  - `MaxLEDs` 定数（12）
  - `MaxLayerColors` 定数（4）
  - `Config` 構造体（`Pin`, `Count`, `LayerColors`）
  - `LED` 構造体
  - `New(cfg *Config) *LED`
  - `Init()`（WS2812 デバイス初期化、デフォルト色で全点灯）
  - `Tick() keycode.Keycode`（keycode.None を返す）
  - `OnLayerChange(layer int)`（全 LED をレイヤー色に更新）
- [x]`examples/zero-kb02/config.go` に LED ピン定義を追加
- [x]`examples/zero-kb02/main.go` に LED の初期化・登録を追加

### テスト

- [x]`make test` 全パス

### ドキュメント

- [x]`docs/packages/led.md` を新規作成

### 検証

- [x]`make lint` パス
- [x]`make build` パス
- [x]実機書き込み
- [x]起動時にデフォルト色（緑）で全 LED が点灯することを確認
- [x]MO(1) ホールド中に LED が青に変わることを確認
- [x]MO(1) リリースで LED が緑に戻ることを確認

### PR・マージ

- [x]`feature/led` → `phase4` の PR 作成
- [x]PR マージ

---

## feature/oled — OLED ディスプレイ（SSD1306）

### 実装

- [x]`engine/peripheral/oled/font.go` を新規作成（ビルドタグなし）
  - `charWidth`, `charHeight` 定数
  - `fontData` — 'L', '0', '1', '2', '3' のビットマップ定義
- [x]`engine/peripheral/oled/oled.go` を新規作成（`//go:build tinygo`）
  - `Config` 構造体（`Bus`, `SDA`, `SCL`, `Address`, `Width`, `Height`）
  - `OLED` 構造体
  - `New(cfg *Config) *OLED`
  - `Init()`（I2C 初期化、SSD1306 初期化、初期表示 "L0"）
  - `Tick() keycode.Keycode`（keycode.None を返す）
  - `OnLayerChange(layer int)`（"L0"〜"L3" を表示）
- [x]`examples/zero-kb02/config.go` に OLED ピン定義を追加
- [x]`examples/zero-kb02/main.go` に OLED の初期化・登録を追加

### テスト

- [x]`make test` 全パス

### ドキュメント

- [x]`docs/packages/oled.md` を新規作成

### 検証

- [x]`make lint` パス
- [x]`make build` パス
- [x]実機書き込み
- [x]起動時に "L0" が OLED に表示されることを確認
- [x]MO(1) ホールド中に "L1" に変わることを確認
- [x]MO(1) リリースで "L0" に戻ることを確認

### PR・マージ

- [x]`feature/oled` → `phase4` の PR 作成
- [x]PR マージ

---

## Phase 4 完了作業

- [x] `docs/functional-design.md` を更新（周辺機器統合のデータフロー・システム構成図追加）
- [x] `docs/repository-structure.md` を更新（新ファイル・ディレクトリ追記）
- [x] `docs/glossary.md` を更新（周辺機器関連の用語追加）
- [x] `ROADMAP.md` の Phase 4 チェックボックスを更新
- [x] `README.md` の開発フェーズ表を更新
- [ ] Phase 4 完了コミット・PR 作成・マージ
- [ ] `v0.4.0` タグを作成
