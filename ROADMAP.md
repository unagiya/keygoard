# ROADMAP

各フェーズは「実機で動作確認済み・main にマージ済み」になって初めて完了とする。
フェーズ完了時に `v0.N.0` タグを打つ。

## Feature Backlog

以下は未着手の機能候補。次の Phase 開始時にこの中から選定して取り込む。

| # | Feature | 概要 | 備考 |
|---|---------|------|------|
| 1 | Split 対応 | UART 双方向通信、マスター/スレーブ分離 | |
| 2 | Gamepad HID | Composite HID（keyboard + gamepad） | TinyGo Issue #3474 の解決が前提 |
| 3 | マクロ / コンボ | マクロ・コンボキー機能 | |
| 4 | Flash 永続化 | キーマップ・設定の Flash 永続化 | |
| 5 | NKRO | N-Key Rollover 対応 | |
| 6 | LED エフェクト拡張 | リアクティブライティング・個別 LED 制御 | |
| 7 | OLED 自由描画 | 任意テキスト・カスタム画像の表示 | |
| 8 | CLI ツール | `keygoard init` コマンド | |

---

## Completed Phases

### Phase 1: 最小 HID キーボード ✅ `v0.1.0`

RP2040 が USB キーボードとして OS に認識され、キー入力が届く。

- マトリクススキャン（COL2ROW）、デバウンス処理
- 基本 6KRO キーボード HID
- レイヤー 0 のみのキーマップ（3×4 固定）

### Phase 2: フレームワーク API ブラッシュアップ ✅ `v0.2.0`

ファームウェア側の記述を最小化し、フレームワークらしい API を提供。

- `engine.Run()` 導入
- USB Product Name のカスタマイズ対応

### Phase 3: レイヤーシステム ✅ `v0.3.0`

- 複数レイヤー（MO / TG）、KC_TRNS（透過キー）
- TT / LT（タップ検出）

### Phase 4: 周辺機器 ✅ `v0.4.0`

- Peripheral インターフェース定義・engine 統合
- ロータリーエンコーダー / RGB LED / OLED SSD1306

### Phase 5: パッケージ構造リファクタリング ✅ `v0.5.0`

フレームワーク利用者が `engine` と `keycode` の 2 パッケージだけで使える構造に移行。

- engine Facade 化、`internal/` へのパッケージ移動
- `engine.Pin` / `I2CBus` / `Rotation` / `LEDType` 型導入

### Phase 6: アナログジョイスティック — マウスポインティングデバイス ✅ `v0.6.0`

- アナログジョイスティック（ADC → マウスカーソル移動）
- Composite HID（keyboard + mouse）
- デッドゾーン・感度設定・軸反転・ボタン対応
