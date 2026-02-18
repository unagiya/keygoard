# キーコードリファレンス

keygoardで使用できる全キーコードの一覧です。

## キーコード体系

keygoardは独自のキーコード体系を採用しています（QMK互換ではありません）。

### 範囲

| 範囲 | 用途 |
|------|------|
| 0x0000-0x00FF | 基本キー（USB HID準拠） |
| 0x0100-0x01FF | モディファイア |
| 0x0200-0x02FF | レイヤー操作 |
| 0x6000-0x6FFF | ゲームパッドボタン |
| 0x7000-0xFFFF | カスタム/将来拡張 |

## 基本キー

### 英字

```go
keycode.KC_A  // A
keycode.KC_B  // B
keycode.KC_C  // C
keycode.KC_D  // D
keycode.KC_E  // E
keycode.KC_F  // F
keycode.KC_G  // G
keycode.KC_H  // H
keycode.KC_I  // I
keycode.KC_J  // J
keycode.KC_K  // K
keycode.KC_L  // L
keycode.KC_M  // M
keycode.KC_N  // N
keycode.KC_O  // O
keycode.KC_P  // P
keycode.KC_Q  // Q
keycode.KC_R  // R
keycode.KC_S  // S
keycode.KC_T  // T
keycode.KC_U  // U
keycode.KC_V  // V
keycode.KC_W  // W
keycode.KC_X  // X
keycode.KC_Y  // Y
keycode.KC_Z  // Z
```

### 数字

```go
keycode.KC_1  // 1
keycode.KC_2  // 2
keycode.KC_3  // 3
keycode.KC_4  // 4
keycode.KC_5  // 5
keycode.KC_6  // 6
keycode.KC_7  // 7
keycode.KC_8  // 8
keycode.KC_9  // 9
keycode.KC_0  // 0
```

### 特殊キー

```go
keycode.KC_ENT   // Enter
keycode.KC_ESC   // Escape
keycode.KC_BSPC  // Backspace
keycode.KC_TAB   // Tab
keycode.KC_SPC   // Space
keycode.KC_MINS  // - (マイナス)
keycode.KC_EQL   // = (イコール)
keycode.KC_LBRC  // [ (左角括弧)
keycode.KC_RBRC  // ] (右角括弧)
keycode.KC_BSLS  // \ (バックスラッシュ)
keycode.KC_SCLN  // ; (セミコロン)
keycode.KC_QUOT  // ' (クォート)
keycode.KC_GRV   // ` (グレイブ)
keycode.KC_COMM  // , (カンマ)
keycode.KC_DOT   // . (ピリオド)
keycode.KC_SLSH  // / (スラッシュ)
keycode.KC_CAPS  // Caps Lock
```

### ファンクションキー

```go
keycode.KC_F1   // F1
keycode.KC_F2   // F2
keycode.KC_F3   // F3
keycode.KC_F4   // F4
keycode.KC_F5   // F5
keycode.KC_F6   // F6
keycode.KC_F7   // F7
keycode.KC_F8   // F8
keycode.KC_F9   // F9
keycode.KC_F10  // F10
keycode.KC_F11  // F11
keycode.KC_F12  // F12
```

### システムキー

```go
keycode.KC_PSCR  // Print Screen
keycode.KC_SLCK  // Scroll Lock
keycode.KC_PAUS  // Pause
keycode.KC_INS   // Insert
keycode.KC_HOME  // Home
keycode.KC_PGUP  // Page Up
keycode.KC_DEL   // Delete
keycode.KC_END   // End
keycode.KC_PGDN  // Page Down
```

### 矢印キー

```go
keycode.KC_RGHT  // →（右矢印）
keycode.KC_LEFT  // ←（左矢印）
keycode.KC_DOWN  // ↓（下矢印）
keycode.KC_UP    // ↑（上矢印）
```

### テンキー

```go
keycode.KC_NLCK  // Num Lock
keycode.KC_PSLS  // /（テンキー）
keycode.KC_PAST  // *（テンキー）
keycode.KC_PMNS  // -（テンキー）
keycode.KC_PPLS  // +（テンキー）
keycode.KC_PENT  // Enter（テンキー）
keycode.KC_P1    // 1（テンキー）
keycode.KC_P2    // 2（テンキー）
keycode.KC_P3    // 3（テンキー）
keycode.KC_P4    // 4（テンキー）
keycode.KC_P5    // 5（テンキー）
keycode.KC_P6    // 6（テンキー）
keycode.KC_P7    // 7（テンキー）
keycode.KC_P8    // 8（テンキー）
keycode.KC_P9    // 9（テンキー）
keycode.KC_P0    // 0（テンキー）
keycode.KC_PDOT  // .（テンキー）
```

### その他

```go
keycode.KC_APP   // Application（メニュー）
keycode.KC_LGUI  // 左GUI（Windows/Command）
keycode.KC_RGUI  // 右GUI（Windows/Command）
```

## モディファイア

```go
keycode.KC_LCTL  // 左Control
keycode.KC_LSFT  // 左Shift
keycode.KC_LALT  // 左Alt
keycode.KC_LGUI  // 左GUI（Windows/Command）
keycode.KC_RCTL  // 右Control
keycode.KC_RSFT  // 右Shift
keycode.KC_RALT  // 右Alt
keycode.KC_RGUI  // 右GUI（Windows/Command）
```

### 使用例

```go
// 通常の使用
{keycode.KC_LCTL, keycode.KC_C}  // Ctrl+C

// モディファイアは他のキーと組み合わせて使用
```

## レイヤー操作

### MO - Momentary（モーメンタリ）

押している間だけレイヤーが有効になります。

```go
keycode.MO(0)   // レイヤー0へ一時切り替え
keycode.MO(1)   // レイヤー1へ一時切り替え
keycode.MO(15)  // レイヤー15へ一時切り替え
```

**使用例：**
```go
// 押している間だけファンクションレイヤーに切り替え
{keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.MO(1)}
```

### TG - Toggle（トグル）

押すたびにレイヤーのON/OFFが切り替わります。

```go
keycode.TG(0)   // レイヤー0をトグル
keycode.TG(1)   // レイヤー1をトグル
keycode.TG(15)  // レイヤー15をトグル
```

**使用例：**
```go
// 押すとゲーミングレイヤーに切り替わり、もう一度押すと戻る
{keycode.KC_ESC, keycode.KC_1, keycode.KC_2, keycode.TG(2)}
```

### TT - Tap Toggle（タップトグル）

- タップ：トグル動作（ON/OFF切り替え）
- 長押し：モーメンタリ動作（押している間のみ）

```go
keycode.TT(0)   // レイヤー0へタップトグル
keycode.TT(1)   // レイヤー1へタップトグル
keycode.TT(15)  // レイヤー15へタップトグル
```

**使用例：**
```go
// タップで固定、長押しで一時的にレイヤー切り替え
{keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.TT(1)}
```

### LT - Layer Tap（レイヤータップ）

- タップ：通常のキーとして動作
- 長押し：レイヤー切り替え

```go
keycode.LT(1, keycode.KC_SPC)   // タップでSpace、長押しでレイヤー1
keycode.LT(2, keycode.KC_ENT)   // タップでEnter、長押しでレイヤー2
```

**使用例：**
```go
// Spaceキーを長押しすると記号レイヤーに切り替わる
{keycode.KC_Z, keycode.KC_X, keycode.LT(1, keycode.KC_SPC), keycode.KC_ENT}
```

**注意：** Phase 1ではタップ検出が未実装のため、長押し動作（モーメンタリ）のみ機能します。

## ゲームパッドボタン

ゲームパッドとして認識されるボタンです。

```go
keycode.GP_BTN1   // ゲームパッドボタン1
keycode.GP_BTN2   // ゲームパッドボタン2
keycode.GP_BTN3   // ゲームパッドボタン3
keycode.GP_BTN4   // ゲームパッドボタン4
keycode.GP_BTN5   // ゲームパッドボタン5
// ... 中略 ...
keycode.GP_BTN32  // ゲームパッドボタン32
```

### 使用例

```go
// ゲーミングキーボードレイアウト
var keymapLayers = [16][][]keycode.Keycode{
    {
        // WASD + ゲームパッドボタン
        {keycode.KC_W, keycode.KC_A, keycode.KC_S, keycode.KC_D},
        {keycode.GP_BTN1, keycode.GP_BTN2, keycode.GP_BTN3, keycode.GP_BTN4},
    },
}
```

## 特殊キーコード

### KC_TRNS - Transparent（透過）

下のレイヤーのキーをそのまま使用します。

```go
keycode.KC_TRNS
```

**使用例：**
```go
// レイヤー1で一部のキーだけ変更し、残りは透過
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー0
    {
        {keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
    },
    // レイヤー1：最初の2つだけ変更
    {
        {keycode.KC_1, keycode.KC_2, keycode.KC_TRNS, keycode.KC_TRNS},
    },
}
```

### KC_NO - No Operation（何もしない）

キーが押されても何も起こりません。

```go
keycode.KC_NO
```

**使用例：**
```go
// 物理的にキーはあるが、使用しない場合
{keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_NO}
```

## キーコードの組み合わせ例

### ゲーミングキーボード

```go
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー0: ゲーミング
    {
        {keycode.KC_ESC, keycode.KC_1, keycode.KC_2, keycode.KC_3},
        {keycode.KC_TAB, keycode.KC_Q, keycode.KC_W, keycode.KC_E},
        {keycode.KC_LSFT, keycode.KC_A, keycode.KC_S, keycode.KC_D},
        {keycode.KC_LCTL, keycode.KC_Z, keycode.KC_X, keycode.MO(1)},
    },
    // レイヤー1: ファンクション
    {
        {keycode.KC_GRV, keycode.KC_F1, keycode.KC_F2, keycode.KC_F3},
        {keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_UP, keycode.KC_TRNS},
        {keycode.KC_TRNS, keycode.KC_LEFT, keycode.KC_DOWN, keycode.KC_RGHT},
        {keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS},
    },
}
```

### マクロパッド

```go
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー0: ショートカット
    {
        {keycode.KC_LCTL, keycode.KC_C},  // Ctrl+C
        {keycode.KC_LCTL, keycode.KC_V},  // Ctrl+V
        {keycode.KC_LCTL, keycode.KC_Z},  // Ctrl+Z
        {keycode.MO(1), keycode.KC_ENT},  // レイヤー切替 + Enter
    },
}
```

## Phase 2以降で追加予定

- カスタムキーコード（0x7000-0xFFFF）
- マクロ機能
- タップダンス
- コンボキー

## 関連ドキュメント

- [レイヤーシステム](layers_ja.md) - レイヤーの詳細な使い方
- [クイックスタート](getting-started_ja.md) - 基本的な使用方法

## ライセンス

MITライセンス
