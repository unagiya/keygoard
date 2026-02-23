# Phase 1 設計

## 実装アプローチ

最小構成で「USB キーボードとして認識・文字入力できる」状態を最短で実現する。
構造を作りすぎず、各 feature ブランチで 1 機能を完成させてから次へ進む。

ブランチ順序: `feature/matrix-scan` → `feature/debounce` → `feature/hid-basic` → `feature/keymap-layer0`

---

## パッケージ構成

```
keygoard/
├── main.go               ← エントリポイント（engine を起動するだけ）
├── engine/
│   └── engine.go         ← メインループ（スキャン → キーコード解決 → HID 送信）
├── matrix/
│   ├── matrix.go         ← マトリクススキャン（//go:build tinygo）
│   ├── debounce.go       ← デバウンスロジック（machine 非依存）
│   └── debounce_test.go  ← デバウンスユニットテスト
├── hid/
│   └── hid.go            ← HID レポート送信（machine/usb/hid ラッパー）
└── keycode/
    ├── keycode.go        ← キーコード定数定義
    ├── keymap.go         ← レイヤー 0 キーマップ（3×4 固定）
    └── keycode_test.go   ← キーコードユニットテスト
```

---

## 各コンポーネントの設計

### matrix パッケージ

**matrix.go**（`//go:build tinygo` タグ付き）

- `Init()`: Col ピンを OutputHigh、Row ピンを InputPulldown に設定
- `Scan(state *[Rows][Cols]bool)`: 各 Col を High にして Row を読み取り、`state` に書き込む
  - ホットパス内でヒープ割り当てなし（固定配列を引数で受け取る）

定数:

```go
const (
    Rows = 3
    Cols = 4
)
// Col: GP5, GP6, GP7, GP8
// Row: GP9, GP10, GP11
```

**debounce.go**（`machine` 非依存、標準 Go でテスト可能）

- カウンタ方式デバウンス（連続 N サイクル同じ状態が続いたら確定）
- `Debounce` 構造体: `[Rows][Cols]uint8` カウンタ + `[Rows][Cols]bool` 確定状態
- `Update(raw *[Rows][Cols]bool, stable *[Rows][Cols]bool)`: raw を読んで stable を更新
- デバウンス閾値: `DebounceThreshold = 5`（定数）

### hid パッケージ

**hid.go**（`//go:build tinygo` タグ付き）

- TinyGo 標準の `machine/usb/hid` を使う（自前 HID 記述子を書かない）
- `Init()`: USB HID キーボードを初期化
- `Send(keys [6]uint8, mods uint8) error`: 6KRO レポートを送信
- USB エニュメレーション完了待ち: 起動後 500ms スリープしてからスキャン開始

### keycode パッケージ

**keycode.go**

- HID Usage ID に基づくキーコード定数（`KC_A`, `KC_B`, ... `KC_NONE` 等）
- `machine` 非依存のため標準 Go でテスト可能

**keymap.go**

- 3×4 固定配列 `var Layer0 [matrix.Rows][matrix.Cols]uint8`
- Phase 1 の割り当て例:

```
KC_Q  KC_W  KC_E  KC_R
KC_A  KC_S  KC_D  KC_F
KC_Z  KC_X  KC_C  KC_V
```

### engine パッケージ

**engine.go**

メインループ（1 goroutine のみ）:

```
1. USB エニュメレーション完了まで待機（500ms）
2. ループ:
   a. matrix.Scan(&rawState)
   b. debouncer.Update(&rawState, &stableState)
   c. stableState から押下キーを最大 6 つ収集
   d. hid.Send(keys, mods)
   e. time.Sleep(1ms)
```

ホットパス内での割り当てなし: 全ての作業用配列は `engine.go` のパッケージ変数として初期化時に確保。

---

## データフロー

```
物理キー押下
    ↓
matrix.Scan()        → rawState [3][4]bool
    ↓
debounce.Update()    → stableState [3][4]bool
    ↓
keymap.Layer0[][]    → keycodes [6]uint8
    ↓
hid.Send()          → USB HID レポート
    ↓
OS（macOS）
```

---

## 影響範囲・依存関係

```
engine → matrix, hid, keycode
matrix → （tinygo のみ）machine
hid    → （tinygo のみ）machine/usb/hid
keycode → なし（他に依存しない）
```

- `keycode` と `matrix/debounce.go` は標準 Go でテスト可能
- `matrix.go` と `hid.go` は `//go:build tinygo` タグで標準 Go ビルドから除外

---

## 未解決事項・リスク

| 事項 | 対応方針 |
|---|---|
| USB エニュメレーション完了タイミング | 500ms 固定スリープで様子見。認識されない場合は延長を検討 |
| デバウンス閾値 | 5 サイクルから開始。チャタリングが出た場合は増やす |
| TinyGo の `machine/usb/hid` API 詳細 | feature/hid-basic ブランチで調査・実装時に確定 |
