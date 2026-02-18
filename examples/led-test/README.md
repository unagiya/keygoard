# LED Test Example

WS2812 RGB LEDのテスト用サンプルです。

## Hardware

- RP2040/RP2350開発ボード
- 4x3キーマトリックス
- WS2812 RGB LED (12個)

## Wiring

### Matrix
- Rows: GP0, GP1, GP2, GP3
- Cols: GP4, GP5, GP6

### LED
- Data: GP16
- VCC: 5V
- GND: GND

**注意**: WS2812は5V動作ですが、RP2040の3.3V信号でも動作することが多いです。
確実に動作させるには、レベルシフターまたは74HCT245などを使用してください。

## Features

この例では3つのレイヤーで異なるLED効果を示します:

### Layer 0 (デフォルト)
- 静的な赤色LED
- 数字キー (1-9, 0)
- MO(1), MO(2)でレイヤー切り替え

### Layer 1
- ブリージング効果 (青色)
- QWERTYキー配列
- MO(1)を押している間有効

### Layer 2
- レインボー効果
- ASDFGHキー配列
- MO(2)を押している間有効

## Building

```bash
# Build for RP2040
tinygo build -target=pico -o firmware.uf2 .

# Flash to device
cp firmware.uf2 /Volumes/RPI-RP2/
```

## LED Effects

このサンプルでは以下のエフェクトを使用します:

1. **StaticEffect**: 単色表示
2. **BreathingEffect**: フェードイン/アウト
3. **RainbowEffect**: 虹色グラデーション

## Customization

エフェクトの変更は、レイヤーごとに異なるエフェクトを設定できます。
詳細は`peripheral/led/effect.go`を参照してください。

### カスタムエフェクトの例

```go
// 静的な緑色
staticGreen := led.NewStaticEffect(color.RGBA{R: 0, G: 255, B: 0, A: 255})

// ゆっくりブリージング (赤)
breathingRed := led.NewBreathingEffect(color.RGBA{R: 255, G: 0, B: 0, A: 255}, 20)

// 速いレインボー
fastRainbow := led.NewRainbowEffect(3)
```

## Phase 2 M2 Status

✅ 基本LED制御
✅ 3種類のエフェクト (Static, Breathing, Rainbow)
✅ エンジンへの統合
⏳ レイヤー連動エフェクト切り替え (Phase 2 M3で改善予定)

## Notes

- 現在のバージョンでは、エフェクトは起動時に設定されます
- レイヤー変更時の自動エフェクト切り替えは、Phase 2 M3で実装予定です
- Split keyboardでのLED制御は、Phase 3で対応予定です
