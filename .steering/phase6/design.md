# Phase 6 設計書: アナログジョイスティック（マウスポインティングデバイス）

## アーキテクチャ

```text
ADC GP29 (X) → uint16 → AxisMapper.MapDelta() → dx int8 → mouse.Move(dx, dy, 0)
ADC GP28 (Y) → uint16 →                       → dy int8 →
GP0 (Button) → bool   → エッジ検出             → keycode or mouse.Click()
```

## パッケージ構成

### 新規パッケージ: `internal/peripheral/joystick/`

```text
internal/peripheral/joystick/
├── axis.go           # 軸マッピングロジック（machine 非依存）
├── axis_test.go      # ユニットテスト
└── joystick.go       # ADC 読み取り + mouse.Move()（//go:build tinygo）
```

既存ペリフェラル（encoder/led/oled）と同じ分離パターンに従う:

- **axis.go** — ビルドタグなし。純粋な Go ロジック。`go test` で検証可能。
- **joystick.go** — `//go:build tinygo`。`machine.ADC` と `mouse` パッケージに依存。

## データフロー

### 軸マッピング（axis.go）

```
ADC 生値 (uint16: 0〜65535)
    │
    ▼
中心値からの偏差 = raw - 32768 (int32)
    │
    ▼
デッドゾーン判定: |偏差| <= deadZone → 0
    │
    ▼
有効偏差 = 偏差 - sign(偏差) * deadZone
    │
    ▼
正規化 = 有効偏差 * maxOutput / (32768 - deadZone)
    │
    ▼
クランプ: -128〜127 → int8
```

`maxOutput` は感度から算出: `sensitivity * 3`（感度 5 → maxOutput 15）

### ボタン処理（joystick.go）

```
GPIO ピン読み取り (Low = 押下)
    │
    ▼
エッジ検出: 前回 false → 今回 true（押下エッジ）
    │
    ├─ ButtonKey == keycode.None → mouse.Port().Click(mouse.Left)、return keycode.None
    │
    └─ ButtonKey != keycode.None → return ButtonKey（engine が HID 送信）
```

## コンポーネント設計

### AxisMapper 構造体（axis.go）

```go
// AxisMapper は ADC 生値をマウス移動量に変換します。
type AxisMapper struct {
    deadZone  int32 // デッドゾーン閾値
    maxOutput int32 // 最大出力値（感度から算出）
    invert    bool  // 軸反転フラグ
}

// NewAxisMapper は AxisMapper を生成します。
func NewAxisMapper(deadZone uint16, sensitivity uint8, invert bool) AxisMapper

// MapDelta は ADC 生値（0〜65535）をマウス移動量（int8）に変換します。
func (m *AxisMapper) MapDelta(raw uint16) int8
```

**設計判断:**

- `int32` で中間計算を行い、オーバーフローを防止
- `deadZone` を差し引いた有効範囲で正規化（デッドゾーン直後から移動開始）
- 感度 1〜10 → maxOutput 3〜30（`sensitivity * 3`）

### Joystick 構造体（joystick.go）

```go
// Joystick はアナログジョイスティックを制御します。
// engine.Peripheral インターフェースを実装します。
type Joystick struct {
    adcX      machine.ADC
    adcY      machine.ADC
    axisX     AxisMapper
    axisY     AxisMapper
    pinButton machine.Pin
    buttonKey keycode.Keycode
    hasButton bool
    prevBtn   bool // エッジ検出用
}
```

**Tick() の動作:**

1. ADC X/Y を読み取り
2. AxisMapper で dx, dy を算出
3. dx != 0 || dy != 0 なら `mouse.Move(dx, dy, 0)` を呼び出し
4. ボタン有効時: エッジ検出 → マウスクリック or キーコード返却
5. キーコードを返す場合以外は `keycode.None` を返す

**Init() の動作:**

1. ADC X/Y ピンを `machine.ADC` として Configure
2. ボタン有効時: ボタンピンを `PinInputPullup` に設定

**OnLayerChange():** 何もしない（エンコーダーと同じ）

## engine パッケージの変更

### JoystickConfig（config.go に追加）

```go
// JoystickConfig はアナログジョイスティックの設定です。
type JoystickConfig struct {
    // PinX は X 軸 ADC ピンの GPIO 番号です。
    PinX Pin
    // PinY は Y 軸 ADC ピンの GPIO 番号です。
    PinY Pin
    // PinButton はボタンピンの GPIO 番号です（EnableButton が true の場合のみ使用）。
    PinButton Pin
    // EnableButton はボタンを有効にするかどうかです。
    EnableButton bool
    // ButtonKey はボタンに割り当てるキーコードです。
    // keycode.None の場合はマウス左クリックを送信します。
    ButtonKey keycode.Keycode
    // Sensitivity はマウス移動の感度（1〜10）です。デフォルト: 5。
    Sensitivity uint8
    // DeadZone はデッドゾーン閾値です。デフォルト: 3000。
    DeadZone uint16
    // InvertX は X 軸を反転するかどうかです。
    InvertX bool
    // InvertY は Y 軸を反転するかどうかです。
    InvertY bool
}
```

### Config への追加（config.go）

```go
type Config struct {
    // ... 既存フィールド ...

    // Joystick はアナログジョイスティックの設定です。nil の場合はスキップします。
    Joystick *JoystickConfig
}
```

### keyboard.go の変更

`New()` 関数に joystick 生成ロジックを追加:

```go
if cfg.Joystick != nil {
    sens := cfg.Joystick.Sensitivity
    if sens == 0 {
        sens = 5
    }
    dz := cfg.Joystick.DeadZone
    if dz == 0 {
        dz = 3000
    }
    kb.addPeripheral(joystick.New(
        machine.Pin(cfg.Joystick.PinX),
        machine.Pin(cfg.Joystick.PinY),
        machine.Pin(cfg.Joystick.PinButton),
        cfg.Joystick.EnableButton,
        cfg.Joystick.ButtonKey,
        sens,
        dz,
        cfg.Joystick.InvertX,
        cfg.Joystick.InvertY,
    ))
}
```

### peripheral.go の変更

```go
const MaxPeripherals = 8  // 4→8 に増加
```

## examples/zero-kb02/main.go の変更

```go
import (
    _ "machine/usb/hid/keyboard"
    _ "machine/usb/hid/mouse"    // 追加: mouse HID 有効化
    // ...
)

func main() {
    kb := engine.New(&engine.Config{
        // ... 既存設定 ...
        Joystick: &engine.JoystickConfig{
            PinX:         29,
            PinY:         28,
            PinButton:    0,
            EnableButton: true,
            ButtonKey:    keycode.None, // マウス左クリック
            Sensitivity:  5,
        },
    })
    kb.Run()
}
```

## デフォルト値処理

| フィールド | ゼロ値 | デフォルト値 | 適用場所 |
|---|---|---|---|
| Sensitivity | 0 | 5 | keyboard.go New() |
| DeadZone | 0 | 3000 | keyboard.go New() |

## 影響範囲

### 変更なし

- `keycode/` — 変更なし
- `internal/matrix/` — 変更なし
- `internal/layer/` — 変更なし
- `internal/tap/` — 変更なし
- `internal/peripheral/encoder/` — 変更なし
- `internal/peripheral/led/` — 変更なし
- `internal/peripheral/oled/` — 変更なし

### 変更あり

- `engine/config.go` — JoystickConfig 追加
- `engine/keyboard.go` — joystick import + 生成ロジック
- `engine/peripheral.go` — MaxPeripherals 増加
- `examples/zero-kb02/main.go` — mouse import + Joystick 設定

### 新規作成

- `internal/peripheral/joystick/axis.go`
- `internal/peripheral/joystick/axis_test.go`
- `internal/peripheral/joystick/joystick.go`
