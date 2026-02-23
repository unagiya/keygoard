# Phase 1 要求仕様

## 概要

RP2040（zero-kb02）を USB キーボードとして OS に認識させ、実際にキー入力が届くことを確認する。
これが keygoard の最初の動作確認マイルストーンとなる。

---

## ユーザーストーリー

1. **キーボードとして認識される**
   - zero-kb02 を Mac に接続したとき、macOS の「システム情報」や USB Prober にキーボードデバイスとして表示される。

2. **キーが打てる**
   - キーマトリクス上の物理キーを押すと、対応する HID キーコードが OS に届き、テキストエディタ等で文字として入力できる。

---

## 受け入れ条件

| # | 条件 |
|---|---|
| 1 | `tinygo build -target=waveshare-rp2040-zero` がエラーなく成功する |
| 2 | macOS の USB Prober / System Information にキーボードデバイスとして表示される |
| 3 | 3×4 キーマトリクスのいずれかのキーを押すと、対応する文字がテキストエディタに入力される |
| 4 | デバウンス処理が有効であること（チャタリングによる誤入力が発生しない） |
| 5 | `go test ./keycode/... ./matrix/...` が全てパスする |

---

## 制約

### スコープ内（Phase 1 で実装する）

- マトリクススキャン（COL2ROW、zero-kb02 ピン配置）
- デバウンス処理
- 基本 6KRO キーボード HID（TinyGo 標準の `machine/usb/hid` を使う）
- レイヤー 0 のみのキーマップ（3×4 固定）

### スコープ外（Phase 1 では実装しない）

- Split 通信 / OLED / RGB LED / ロータリーエンコーダー / ジョイスティック
- 複数レイヤー / MO / TG / TT / LT / KC_TRNS
- Gamepad / Composite HID
- マクロ / コンボ / NKRO / Flash 永続化

### 技術制約

- TinyGo 制約を全て遵守すること（`reflect` 禁止、`fmt.Sprintf` 禁止、ヒープ割り当て最小化）
- サードパーティパッケージの追加はユーザーへの確認が必要
- RAM 使用量 128KB 以下

---

## ハードウェア前提

対象ボード: zero-kb02（Waveshare RP2040-Zero ベース）

| 機能 | GPIO |
|------|------|
| Matrix Col 0–3（出力） | GP5, GP6, GP7, GP8 |
| Matrix Row 0–2（入力） | GP9, GP10, GP11 |
| Matrix 方式 | COL2ROW（列を High にして行を Read） |
