# NKRO（Nキーロールオーバー）

keygoardは6KRO（6キーロールオーバー）とNKRO（全キー同時押し対応）の動的切り替えをサポートしています。

[English](nkro.md) | 日本語

## 概要

- 6KROとNKROの動的切り替え
- ビットマップベースのNKROレポート（最大120キー同時押し）
- キーコードによるモード切り替え（トグル / オン / オフ）
- CompositeHIDによるキーボード + ゲームパッドの統合管理

## 設定

### BoardConfig

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/hid"
)

cfg := &engine.BoardConfig{
    DefaultHIDMode: hid.HIDModeNKRO,  // 起動時からNKROモード
    // ...
}
```

### HIDMode

```go
type HIDMode uint8

const (
    HIDMode6KRO HIDMode = 0  // 6KRO（デフォルト）
    HIDModeNKRO HIDMode = 1  // NKRO
)
```

## キーコードによる切り替え

キーマップに以下のキーコードを配置することで、ユーザーがモードを切り替えられます。

| キーコード | 値 | 説明 |
|---|---|---|
| `KC_NKRO_TOGGLE` | `0x7000` | NKRO トグル（6KRO ↔ NKRO） |
| `KC_NKRO_ON` | `0x7001` | NKRO オン |
| `KC_NKRO_OFF` | `0x7002` | NKRO オフ |

```go
keymap := &engine.Keymap{
    Layers: [16][][]keycode.Keycode{
        { // レイヤー0
            {keycode.KC_A, keycode.KC_B, keycode.KC_C, keycode.KC_NKRO_TOGGLE},
        },
    },
}
```

## API

### UnifiedKeyboardHID

6KROとNKROを透過的に切り替えるラッパーです。

```go
uhid := hid.NewUnifiedKeyboardHID(hid.HIDMode6KRO)

// キー操作（現在のモードに応じて内部で適切なレポートを使用）
uhid.AddKey(kc)
uhid.SetModifier(mask)
uhid.SendReport()
uhid.Clear()

// モード切り替え
uhid.SetMode(hid.HIDModeNKRO)
uhid.ToggleMode()

// 状態取得
mode := uhid.GetMode()
report := uhid.GetReport()       // 6KROレポート
nkroReport := uhid.GetNKROReport()  // NKROレポート
```

### NKROReport

```go
type NKROReport struct {
    Modifier uint8     // モディファイアバイト
    Reserved uint8
    Keys     [15]byte  // ビットマップ（0x04〜0x7Bの120キー）
}

report := &hid.NKROReport{}
report.Clear()
report.AddKey(0x04)     // キーを追加（0x04〜0x7Bのみ）
report.RemoveKey(0x04)  // キーを削除
report.HasKey(0x04)     // キーが含まれるか確認
data := report.ToBytes()  // 17バイトのレポートデータ
```

### CompositeHID

キーボードとゲームパッドのHIDレポートを統合管理します。

```go
chid := hid.NewCompositeHID()                       // 6KROデフォルト
chid := hid.NewCompositeHIDWithMode(hid.HIDModeNKRO)  // NKROで開始

// 各デバイスへのアクセス
keyboard := chid.Keyboard()  // *UnifiedKeyboardHID
gamepad := chid.Gamepad()    // *GamepadHID

// レポート送信
chid.SendReports()
chid.Clear()
```

## 6KROとNKROの違い

| 項目 | 6KRO | NKRO |
|---|---|---|
| 同時押しキー数 | 最大6キー + モディファイア | 最大120キー + モディファイア |
| レポートサイズ | 8バイト | 17バイト |
| 互換性 | すべてのOSで動作 | 一部のBIOSで非対応の場合あり |
| デフォルト | はい | いいえ |

## 使用例

```go
package main

import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/hid"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    cfg := &engine.BoardConfig{
        DefaultHIDMode: hid.HIDModeNKRO,
        // ... その他の設定 ...
    }

    keymap := &engine.Keymap{
        Layers: [16][][]keycode.Keycode{
            { // レイヤー0: 通常キー + NKROトグル
                {keycode.KC_A, keycode.KC_B, keycode.KC_C, keycode.KC_NKRO_TOGGLE},
            },
            { // レイヤー1: NKRO ON/OFF
                {keycode.KC_NKRO_ON, keycode.KC_NKRO_OFF, keycode.KC_TRNS, keycode.KC_TRNS},
            },
        },
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)
    kb.Run()
}
```

## 制約・制限

- NKROモードでもTinyGoのUSBスタックの制約上、USBレポートは6KRO互換形式で送信される（ビットマップから最大6キーを抽出）。完全なNKROはTinyGoのカスタムHID記述子対応後に対応予定
- NKROレポートの対応キー範囲は0x04〜0x7B（120キー）
- 一部のBIOSやブートローダーはNKROレポートを認識しない場合がある
- Storageが有効な場合、HIDモードの状態は永続化される
