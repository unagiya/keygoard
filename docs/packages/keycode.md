# keycode パッケージ仕様

## 概要

HID キーコードの型と定数を定義します。`machine` パッケージに依存しないため、標準 Go でテスト可能です。

## Keycode 型

`uint16` の型定義です。`machine/usb/hid/keyboard` の `Keycode` と同じエンコーディングを採用しているため、`keyboard.Keycode(kc)` のキャスト 1 つで HID 送信できます。

| 範囲            | 意味                           |
|-----------------|--------------------------------|
| `0x0000`        | 無効（キー未割り当て）         |
| `0x0001`        | 透過（下位レイヤー参照）       |
| `0xA000-0xA00F` | MO(layer) — Momentary レイヤー |
| `0xA010-0xA01F` | TG(layer) — Toggle レイヤー    |
| `0xA020-0xA02F` | TT(layer) — Tap-Toggle レイヤー|
| `0xB000-0xBFFF` | LT(layer, kc) — Layer-Tap      |
| `0xE000-0xE0FF` | 修飾キー（Ctrl/Shift/Alt/GUI） |
| `0xF000-0xFFFF` | 通常キー（HID Usage Page 7）   |

## 定数

| グループ       | 定数名                                                                                     |
|----------------|--------------------------------------------------------------------------------------------|
| 無効           | `None`                                                                                     |
| 透過           | `TRNS`                                                                                     |
| 修飾キー       | `ModLeftCtrl`, `ModLeftShift`, `ModLeftAlt`, `ModLeftGUI`, `ModRight*`                     |
| アルファベット | `A`〜`Z`                                                                                   |
| 数字           | `Num0`〜`Num9`                                                                             |
| 基本操作       | `Enter`, `Escape`, `Backspace`, `Tab`, `Space`                                             |
| 記号           | `Minus`, `Equal`, `LBracket`, `RBracket`, `Backslash`, `Semicolon`, `Quote`, `Grave`, `Comma`, `Dot`, `Slash` |
| ファンクション | `F1`〜`F12`, `CapsLock`                                                                    |
| ナビゲーション | `Insert`, `Home`, `PageUp`, `Delete`, `End`, `PageDown`, 矢印 4 キー                      |

## レイヤーアクション生成関数

| 関数 | 説明 |
|------|------|
| `MO(layer uint8) Keycode` | Momentary — ホールド中のみレイヤーを有効化 |
| `TG(layer uint8) Keycode` | Toggle — 押下のたびにレイヤーを ON/OFF |
| `TT(layer uint8) Keycode` | Tap-Toggle — タップでトグル、ホールドで MO |
| `LT(layer uint8, kc Keycode) Keycode` | Layer-Tap — タップでキーコード送信、ホールドでレイヤー有効化 |

## 判定メソッド

| メソッド | 説明 |
|----------|------|
| `(kc Keycode) IsMO() bool` | MO キーコードか |
| `(kc Keycode) IsTG() bool` | TG キーコードか |
| `(kc Keycode) IsTT() bool` | TT キーコードか |
| `(kc Keycode) IsLT() bool` | LT キーコードか |
| `(kc Keycode) IsLayerAction() bool` | MO/TG/TT/LT のいずれか |
| `(kc Keycode) IsTapAction() bool` | TT/LT のいずれか（タップ検出が必要） |

## 情報抽出メソッド

| メソッド | 説明 |
|----------|------|
| `(kc Keycode) Layer() int` | レイヤー番号を抽出 |
| `(kc Keycode) TapKeycode() Keycode` | LT のタップ時キーコードを抽出（LT 以外は None） |

## テスト

```bash
go test ./keycode/...
```
