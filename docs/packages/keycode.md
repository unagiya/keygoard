# keycode パッケージ仕様

## 概要

HID キーコードの型と定数を定義します。`machine` パッケージに依存しないため、標準 Go でテスト可能です。

## Keycode 型

`uint16` の型定義です。`machine/usb/hid/keyboard` の `Keycode` と同じエンコーディングを採用しているため、`keyboard.Keycode(kc)` のキャスト 1 つで HID 送信できます。

| 範囲            | 意味                           |
|-----------------|--------------------------------|
| `0x0000`        | 無効（キー未割り当て）         |
| `0xE000-0xE0FF` | 修飾キー（Ctrl/Shift/Alt/GUI） |
| `0xF000-0xFFFF` | 通常キー（HID Usage Page 7）   |

## Phase 1 で定義する定数

| グループ       | 定数名                                                                                     |
|----------------|--------------------------------------------------------------------------------------------|
| 無効           | `None`                                                                                     |
| 修飾キー       | `ModLeftCtrl`, `ModLeftShift`, `ModLeftAlt`, `ModLeftGUI`, `ModRight*`                     |
| アルファベット | `A`〜`Z`                                                                                   |
| 数字           | `Num0`〜`Num9`                                                                             |
| 基本操作       | `Enter`, `Escape`, `Backspace`, `Tab`, `Space`                                             |
| 記号           | `Minus`, `Equal`, `LBracket`, `RBracket`, `Backslash`, `Semicolon`, `Quote`, `Grave`, `Comma`, `Dot`, `Slash` |
| ファンクション | `F1`〜`F12`, `CapsLock`                                                                    |
| ナビゲーション | `Insert`, `Home`, `PageUp`, `Delete`, `End`, `PageDown`, 矢印 4 キー                      |

## テスト

```bash
go test ./keycode/...
```
