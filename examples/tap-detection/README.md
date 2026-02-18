# Tap Detection Example

TT (Tap Toggle) と LT (Layer Tap) のタップ検出機能を示すサンプルです。Phase 2 M4で実装されました。

## Features

### TT (Tap Toggle)
- **タップ**: レイヤーをトグル (ON/OFF切り替え)
- **ホールド**: レイヤーを一時的に有効化 (離すと無効化)
- **判定時間**: 150ms以内にリリースでタップ、それ以上でホールド

### LT (Layer Tap)
- **タップ**: 指定されたキーを送信
- **ホールド**: レイヤーを一時的に有効化
- **判定時間**: 150ms以内にリリースでタップ、それ以上でホールド

## Hardware

- RP2040/RP2350開発ボード
- 3x3キーマトリックス
- ジョイスティック (ADC0)
- SSD1306 OLED (128x32, I2C0)

## Wiring

### Matrix
- Rows: GP0, GP1, GP2
- Cols: GP3, GP4, GP5

### Peripherals
- Joystick: ADC0 (GP26)
- OLED: I2C0 (SDA: GP4, SCL: GP5)

## Keymap

### Layer 0 (デフォルト)
```
A      B       C
D      E       F
TT(1)  LT(2,Space)  Enter
```

- **TT(1)**:
  - タップ: Layer 1をトグル (再度タップで解除)
  - ホールド: Layer 1を一時有効化

- **LT(2, Space)**:
  - タップ: Spaceキーを送信
  - ホールド: Layer 2を一時有効化

### Layer 1 (TT(1)でトグル)
```
1  2  3
4  5  6
-  7  8
```

### Layer 2 (LT(2, Space)でホールド)
```
F1  F2  F3
F4  F5  F6
F7  -   F8
```

## Usage Examples

### TT(1)キーの使い方

#### タップ (Layer 1固定)
1. TT(1)キーを150ms以内に押して離す
2. Layer 1が有効になる
3. 数字キー (1-8) が使える
4. もう一度TT(1)をタップでLayer 0に戻る

#### ホールド (Layer 1一時)
1. TT(1)キーを150ms以上押し続ける
2. Layer 1が一時的に有効になる
3. TT(1)を離すとLayer 0に戻る

### LT(2, Space)キーの使い方

#### タップ (Spaceキー)
1. LT(2, Space)を150ms以内に押して離す
2. Spaceキーが送信される
3. Layer 0のまま

#### ホールド (Layer 2一時)
1. LT(2, Space)を150ms以上押し続ける
2. Layer 2が一時的に有効になる
3. Fキー (F1-F8) が使える
4. キーを離すとLayer 0に戻る

## Timing Configuration

タップ検出のタイミングは `engine/tap.go` で設定されています:

```go
type TapConfig struct {
    TappingTerm time.Duration // Tap detection window (default: 200ms)
    HoldTerm    time.Duration // Hold threshold (default: 150ms)
}
```

### デフォルト設定
- **HoldTerm**: 150ms - これより短い場合タップ判定
- **TappingTerm**: 200ms - タップ検出のウィンドウ時間

### カスタマイズ

タイミングをカスタマイズするには、以下のようにTapConfigを作成:

```go
customConfig := &engine.TapConfig{
    TappingTerm: 250 * time.Millisecond,
    HoldTerm:    180 * time.Millisecond,
}
```

## Joystick Calibration

このサンプルはジョイスティックキャリブレーション機能も含んでいます。

### 起動時の自動キャリブレーション
- 起動時にジョイスティックの中心位置を自動的に読み取ります
- ジョイスティックを中央に保って起動してください

### 手動キャリブレーション

キャリブレーションAPIを使用して、実行時に再キャリブレーション可能:

```go
// ジョイスティックを中央に保って実行
kb.CalibrateJoystick()

// 現在のキャリブレーション値を取得
centerX, centerY := kb.GetJoystickCenter()

// 生のADC値を取得 (キャリブレーションUI用)
rawX, rawY := kb.GetJoystickRaw()
```

### デッドゾーン調整

ジョイスティックのデッドゾーンは調整可能です:

```go
if joy := kb.GetJoystick(); joy != nil {
    // デッドゾーンを10%に設定
    joy.SetDeadzone(4095 * 10 / 100)

    // 現在のデッドゾーンを取得
    dz := joy.GetDeadzone()
}
```

デフォルトは5% (約205) です。

## Building

```bash
tinygo build -target=pico -o firmware.uf2 .
cp firmware.uf2 /Volumes/RPI-RP2/
```

## Testing

### TT(1)のテスト

1. **短押し (タップ)**
   - TT(1)を素早く押して離す
   - Layer 1に切り替わる (OLED表示で確認)
   - 'E'キーを押すと '5' が入力される
   - 再度TT(1)をタップでLayer 0に戻る

2. **長押し (ホールド)**
   - TT(1)を200ms以上押し続ける
   - Layer 1が一時有効になる
   - TT(1)を離すとLayer 0に戻る

### LT(2, Space)のテスト

1. **短押し (タップ)**
   - LT(2, Space)を素早く押して離す
   - Spaceキーが入力される
   - Layer 0のまま

2. **長押し (ホールド)**
   - LT(2, Space)を200ms以上押し続ける
   - Layer 2が一時有効になる
   - 'E'キーを押すと 'F5' が入力される
   - キーを離すとLayer 0に戻る

## Implementation Details

### Tap Detection State Machine

各タップ可能なキーは以下の状態を持ちます:

1. **Idle**: キーが押されていない
2. **Pressed**: キーが押された直後、タップかホールドか判定中
3. **Holding**: ホールド判定、レイヤーが有効化されている
4. **Tapped**: タップ判定、キーコードを送信済み

### Processing Flow

1. キーが押される → Pressed状態
2. 150ms以内にリリース → Tapped (タップ処理)
3. 150ms経過 → Holding (ホールド処理、レイヤー有効化)
4. リリース → Idle

## Phase 2 M4 Status

✅ TT (Tap Toggle) 実装
✅ LT (Layer Tap) 実装
✅ タップ検出ステートマシン
✅ タイミング設定可能
✅ ジョイスティックキャリブレーション
✅ 手動再キャリブレーション
✅ デッドゾーン調整

## Known Limitations

- 現在のバージョンでは、LTのタップキーは8ビット (0x00-0xFF) に制限されています
- モディファイヤーキーとの組み合わせは Phase 3 で改善予定
- キャリブレーション値の永続化は Phase 3 で実装予定

## Next Steps

Phase 3 では以下の機能が追加予定:
- より高度なLED効果 (レイヤー反応型)
- NKRO (N-Key Rollover) 対応
- キャリブレーション値の保存/復元
- パフォーマンス最適化

## Notes

- タップ検出は全てのキーで動作しますが、TT/LTキーコードのみが特別な処理を受けます
- ジョイスティックキャリブレーションは起動時に自動実行されます
- OLED表示でレイヤー状態を確認できます
