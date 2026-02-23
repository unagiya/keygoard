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

## Phase 2: レイヤーシステム（Phase 1 完了後）

- 複数レイヤー（MO / TG）
- KC_TRNS（透過キー）
- TT / LT（タップ検出）

---

## Phase 3: 周辺機器（Phase 2 完了後）

- ロータリーエンコーダー（GP3/GP4）
- RGB LED WS2812（GP1、12 LED）
- OLED SSD1306（GP12/GP13）

---

## Phase 4: Joystick + Gamepad（Phase 3 完了後）

- アナログジョイスティック（GP28/GP29）
- Composite HID（キーボード＋ゲームパッド）← Composite はここまで遅らせる

---

## Phase 5: Split 対応（Phase 4 完了後）

- UART 双方向通信
- マスター/スレーブ分離

---

## Phase 6: 高度な機能（Phase 5 完了後）

- マクロ / コンボ
- Flash 永続化（キーマップ・設定）
- NKRO
- CLI ツール（`keygoard init`）
