# ROADMAP

各フェーズは「実機で動作確認済み・main にマージ済み」になって初めて完了とする。
フェーズ完了時に `v0.N.0` タグを打つ。

## Phase 1（現在の目標）: 最小 HID キーボード

**ゴール: RP2040 が USB キーボードとして OS に認識され、キー入力が届く**

> コード実装は `feature/matrix-scan` ブランチで完了済み。
> 残りは実機ビルド・動作確認のみ。

- [x] マトリクススキャン（COL2ROW、zero-kb02 のピン配置で動作）
- [x] デバウンス処理
- [x] 基本 6KRO キーボード HID（TinyGo 標準の `machine/usb/hid` を使う）
- [x] レイヤー 0 のみのキーマップ（3×4 固定）
- [x] `tinygo build -target=waveshare-rp2040-zero` でビルド成功
- [x] macOS でキーボードとして認識されることを確認
- [x] 実際に文字が入力できることを確認

**Phase 1 スコープ外（追加しない）:**
- Split / OLED / LED / Encoder / Joystick
- 複数レイヤー / MO / TG / TT / LT / KC_TRNS
- Gamepad / Composite HID
- マクロ / コンボ / NKRO / Flash 永続化

---

## Phase 2: フレームワーク API ブラッシュアップ（Phase 1 完了後）

**ゴール: ファームウェア側の記述を最小化し、フレームワークらしい API を提供する**

- [x] `engine.Run()` 導入（メインループをエンジン内部に隠蔽）
- [x] USB Product Name のカスタマイズ対応

---

## Phase 3: レイヤーシステム（Phase 2 完了後）

- [x] 複数レイヤー（MO / TG）
- [x] KC_TRNS（透過キー）
- [x] TT / LT（タップ検出）

---

## Phase 4: 周辺機器（Phase 3 完了後）

- [x] Peripheral インターフェース定義・engine 統合
- [x] ロータリーエンコーダー（GP3/GP4）
- [x] RGB LED WS2812/SK6812（GP1、12 LED）
- [x] OLED SSD1306（GP12/GP13、ソフトウェア回転対応）

---

## Phase 5: パッケージ構造リファクタリング（Phase 4 完了後）

**ゴール: フレームワーク利用者が `engine` と `keycode` の 2 パッケージだけで使えるようにする**

- [x] engine パッケージを Facade 化（matrix・ペリフェラルの設定を Config に統合）
- [x] `engine.Pin` / `I2CBus` / `Rotation` / `LEDType` 型導入（machine 隠蔽）
- [x] `internal/` に layer / tap / matrix を移動（外部 import 不可）
- [x] `internal/peripheral/` に encoder / led / oled を移動
- [x] ドキュメント・テストの整合性を維持
- [x] README.md をユーザー向けドキュメントとして整備

---

## Phase 6: Joystick + Gamepad（Phase 5 完了後）

- アナログジョイスティック（GP28/GP29）
- Composite HID（キーボード＋ゲームパッド）← Composite はここまで遅らせる

---

## Phase 7: Split 対応（Phase 6 完了後）

- UART 双方向通信
- マスター/スレーブ分離

---

## Phase 8 以降: 高度な機能（Phase 7 完了後）

以下の機能は優先度・依存関係に応じて Phase 8 以降で段階的に実装する。

- マクロ / コンボ
- Flash 永続化（キーマップ・設定）
- NKRO
- LED エフェクト拡張（リアクティブライティング・個別 LED 制御）
- OLED 自由描画（任意テキスト・カスタム画像の表示）
- CLI ツール（`keygoard init`）
