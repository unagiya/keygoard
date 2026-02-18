# マクロ・コンボキー

keygoardはキーシーケンスマクロと複数キー同時押しによるコンボキーをサポートしています。

[English](macro-combo.md) | 日本語

## 概要

### マクロ
- 最大16個のマクロを定義可能
- 各マクロは最大32ステップ
- キー押下 / 離上 / タップ / 遅延のステップタイプ
- `BoardConfig`での宣言的定義と`RegisterMacro()` APIによる動的登録

### コンボキー
- 最大16個のコンボを定義可能
- 各コンボは2〜4キーの同時押し
- 検出ウィンドウは設定可能（デフォルト50ms）
- `BoardConfig`での宣言的定義と`RegisterCombo()` APIによる動的登録

## マクロ

### マクロステップ

```go
type MacroStepType uint8

const (
    MacroStepEnd     MacroStepType = 0x00  // 終端
    MacroStepKeyDown MacroStepType = 0x01  // キー押下
    MacroStepKeyUp   MacroStepType = 0x02  // キー離上
    MacroStepKeyTap  MacroStepType = 0x03  // 押下 + 即離上
    MacroStepDelay   MacroStepType = 0x04  // 遅延（Valueにミリ秒を指定）
)

type MacroStep struct {
    Type  MacroStepType
    Value uint16  // キーコード値 または 遅延ミリ秒
}
```

### BoardConfigでの定義

```go
cfg := &engine.BoardConfig{
    Macros: []engine.MacroDef{
        {
            ID: 0,  // マクロID（0〜15）
            Steps: []engine.MacroStep{
                // Ctrl+C を送信
                {Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
                {Type: engine.MacroStepKeyTap,  Value: uint16(keycode.KC_C)},
                {Type: engine.MacroStepKeyUp,   Value: uint16(keycode.KC_LCTL)},
            },
        },
        {
            ID: 1,
            Steps: []engine.MacroStep{
                // "gg" と入力して100ms待機
                {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_G)},
                {Type: engine.MacroStepDelay,  Value: 100},
                {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_G)},
            },
        },
    },
}
```

### キーマップでの割り当て

マクロはキーコード `KC_MACRO0` 〜 `KC_MACRO15`（0x7010〜0x701F）で起動します。

```go
keymap := &engine.Keymap{
    Layers: [16][][]keycode.Keycode{
        { // レイヤー0
            {keycode.KC_MACRO0, keycode.KC_MACRO1, keycode.KC_A, keycode.KC_B},
        },
    },
}
```

### 動的登録API

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

err := kb.RegisterMacro(2, []engine.MacroStep{
    {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_H)},
    {Type: engine.MacroStepKeyTap, Value: uint16(keycode.KC_I)},
})
```

### MacroManager API

```go
mm := engine.NewMacroManager()

// マクロ登録
mm.RegisterMacro(id, steps)

// マクロ起動
mm.Trigger(id)

// 実行状態確認
running := mm.IsRunning()

// 毎スキャン呼び出し（エンジンが自動実行）
kc, stepType, active := mm.Update()
```

## コンボキー

### BoardConfigでの定義

```go
cfg := &engine.BoardConfig{
    Combos: []engine.ComboDef{
        {
            Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_A, keycode.KC_B},
            Count:  2,              // 構成キー数
            Output: keycode.KC_ESC, // 成立時に発行するキーコード
        },
        {
            Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_J, keycode.KC_K, keycode.KC_L},
            Count:  3,
            Output: keycode.KC_ENT,
        },
    },
    ComboWindow: 50 * time.Millisecond,  // 検出ウィンドウ（デフォルト50ms）
}
```

### ComboDef構造体

| フィールド | 型 | 説明 |
|---|---|---|
| `Keys` | `[4]keycode.Keycode` | 構成キー（`KC_NO`で終端） |
| `Count` | `uint8` | 構成キー数（2〜4） |
| `Output` | `keycode.Keycode` | コンボ成立時に発行するキーコード |

### 動的登録API

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

err := kb.RegisterCombo(
    []keycode.Keycode{keycode.KC_D, keycode.KC_F},
    keycode.KC_TAB,
)
```

### ComboDetector API

```go
cd := engine.NewComboDetector(50 * time.Millisecond)

// コンボ登録
cd.RegisterCombo(keys, output)

// キー押下処理
outputKC, consumed := cd.ProcessKeyPress(kc)
// consumed=true → 元のキーはコンボ判定のため保留される

// キー離上処理
cd.ProcessKeyRelease(kc)

// タイムアウト処理（毎スキャン呼び出し）
pendingKeys, count := cd.Update()
// タイムアウト時に保留されていたキーが返却される
```

### コンボの動作

1. コンボに含まれるキーが押されると、そのキーは一時的に保留される
2. 検出ウィンドウ内にすべての構成キーが押されると、コンボが成立し`Output`キーコードが発行される
3. ウィンドウがタイムアウトすると、保留されていたキーが通常のキー入力として処理される

## 使用例

```go
package main

import (
    "time"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    cfg := &engine.BoardConfig{
        // ... マトリクス設定等 ...
        Macros: []engine.MacroDef{
            {
                ID: 0,
                Steps: []engine.MacroStep{
                    // Ctrl+Z（元に戻す）
                    {Type: engine.MacroStepKeyDown, Value: uint16(keycode.KC_LCTL)},
                    {Type: engine.MacroStepKeyTap,  Value: uint16(keycode.KC_Z)},
                    {Type: engine.MacroStepKeyUp,   Value: uint16(keycode.KC_LCTL)},
                },
            },
        },
        Combos: []engine.ComboDef{
            {
                Keys:   [engine.MaxComboKeys]keycode.Keycode{keycode.KC_A, keycode.KC_S},
                Count:  2,
                Output: keycode.KC_ESC,
            },
        },
        ComboWindow: 50 * time.Millisecond,
    }

    keymap := &engine.Keymap{
        Layers: [16][][]keycode.Keycode{
            {
                {keycode.KC_MACRO0, keycode.KC_A, keycode.KC_S, keycode.KC_D},
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## 制約・制限

### マクロ
- 最大16個のマクロ
- 各マクロは最大32ステップ
- マクロは同時に1つのみ実行可能
- マクロ実行中は毎スキャンサイクルで1ステップずつ処理される

### コンボキー
- 最大16個のコンボ
- 各コンボは2〜4キーで構成
- コンボ検出ウィンドウのデフォルトは50ms
- コンボのキーが保留される間、わずかな入力遅延が発生する
