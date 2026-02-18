# ロータリーエンコーダー

keygoardはロータリーエンコーダーをサポートし、回転操作やクリックをキーコードにマッピングできます。

[English](encoder.md) | 日本語

## 概要

- 最大2個のロータリーエンコーダーに対応
- 時計回り / 反時計回りのイベント検出
- クリック（プッシュスイッチ）対応
- 任意のキーコードへのマッピング
- Gray codeルックアップテーブルによる高精度な回転検出

## 設定

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/peripheral/encoder"
    "machine"
)

cfg := &engine.BoardConfig{
    Encoders: []*encoder.Config{
        {
            PinA:     machine.GP10,         // エンコーダAピン（プルアップ入力）
            PinB:     machine.GP11,         // エンコーダBピン（プルアップ入力）
            PinClick: machine.GP12,         // クリックボタンピン
            CW:       keycode.KC_VOLU,      // 時計回り → 音量アップ
            CCW:      keycode.KC_VOLD,      // 反時計回り → 音量ダウン
            Click:    keycode.KC_MUTE,      // クリック → ミュート
        },
        {
            PinA:     machine.GP13,
            PinB:     machine.GP14,
            PinClick: machine.NoPin,        // クリック不要の場合
            CW:       keycode.KC_PGDN,
            CCW:      keycode.KC_PGUP,
        },
    },
}
```

### Config構造体

| フィールド | 型 | 説明 |
|---|---|---|
| `PinA` | `machine.Pin` | エンコーダAピン（プルアップ入力） |
| `PinB` | `machine.Pin` | エンコーダBピン（プルアップ入力） |
| `PinClick` | `machine.Pin` | クリックボタンピン（不要なら`machine.NoPin`） |
| `CW` | `keycode.Keycode` | 時計回り回転に割り当てるキーコード |
| `CCW` | `keycode.Keycode` | 反時計回り回転に割り当てるキーコード |
| `Click` | `keycode.Keycode` | クリックに割り当てるキーコード（任意） |

## イベント

```go
type Event uint8

const (
    EventNone  Event = iota  // 変化なし
    EventCW                  // 時計回り
    EventCCW                 // 反時計回り
    EventClick               // ボタンクリック
)
```

## API

### Encoder

```go
// コンストラクタ
enc := encoder.New(cfg)

// 更新（毎1ms呼び出し、エンジンが自動実行）
event := enc.Update()

// イベントからキーコードへ変換
kc := enc.GetKeycode(event)

// 累積回転カウント
pos := enc.GetPosition()

// カウントリセット
enc.ResetPosition()
```

### エンジンとの統合

`BoardConfig.Encoders`に設定を渡すだけで、エンジンが自動的に以下を行います：

1. 毎スキャンサイクルで`Update()`を呼び出し
2. イベントが発生した場合、対応するキーコードを通常のキー入力として処理
3. HIDレポートに反映

## 使用例

```go
package main

import (
    "machine"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/peripheral/encoder"
)

func main() {
    cfg := &engine.BoardConfig{
        // ... マトリクス設定等 ...
        Encoders: []*encoder.Config{
            {
                PinA:     machine.GP10,
                PinB:     machine.GP11,
                PinClick: machine.GP12,
                CW:       keycode.KC_VOLU,
                CCW:      keycode.KC_VOLD,
                Click:    keycode.KC_MUTE,
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## 制約・制限

- 最大2個のエンコーダーに対応
- ピンはプルアップ入力として設定される
- 1ms間隔でのポーリングにより回転を検出（割り込みではない）
- 高速回転時にステップが飛ぶ場合がある
