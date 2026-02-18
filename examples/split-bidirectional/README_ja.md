# Split Bidirectional サンプル

分割キーボードの双方向通信サンプルです。Phase 2 M3で実装された機能を示します。

[English README](README.md)

## 特徴

このサンプルは以下の機能を実装しています:

### マスター → スレーブ通信
- **LED同期**: マスター側のLEDエフェクトがスレーブ側にリアルタイムで同期
- **OLED同期**: マスター側の表示内容がスレーブ側に送信される

### スレーブ → マスター通信
- **マトリックス状態**: スレーブ側のキー入力をマスターに送信
- **ジョイスティック**: スレーブ側のジョイスティック(Y軸)をマスターに送信

## 必要なハードウェア

### マスター側
- RP2040/RP2350開発ボード
- 4x3キーマトリックス
- ジョイスティック (X軸: ADC0)
- SSD1306 OLED (128x32, I2C0)
- WS2812 RGB LED x12 (GP16)
- UART1 (TX/RX)

### スレーブ側
- RP2040/RP2350開発ボード
- 4x3キーマトリックス
- ジョイスティック (Y軸: ADC0)
- SSD1306 OLED (128x32, I2C0)
- WS2812 RGB LED x12 (GP16)
- UART1 (TX/RX)

## 配線

### マスター

#### マトリックス
- 行（Row）: GP0, GP1, GP2, GP3
- 列（Col）: GP4, GP5, GP6

#### 周辺機器
- ジョイスティック X: ADC0 (GP26)
- OLED: I2C0 (SDA: GP4, SCL: GP5)
- LED: GP16

#### UART（スレーブへ）
- TX: UART1 TX (GP8)
- RX: UART1 RX (GP9)
- GND: 共通グランド

### スレーブ

#### マトリックス
- 行（Row）: GP0, GP1, GP2, GP3
- 列（Col）: GP4, GP5, GP6

#### 周辺機器
- ジョイスティック Y: ADC0 (GP26)
- OLED: I2C0 (SDA: GP4, SCL: GP5)
- LED: GP16

#### UART（マスターへ）
- TX: UART1 TX (GP8)
- RX: UART1 RX (GP9)
- GND: 共通グランド

**重要**: マスターとスレーブのGNDを共通接続してください。

## プロトコル

### ボーレート
460800 bps（高速通信）

### メッセージタイプ

#### スレーブ → マスター
- `0x01`: マトリックス状態（ビットパッキング）
- `0x02`: ジョイスティックデータ (X, Y: int8)

#### マスター → スレーブ
- `0x10`: LED同期（RGB色データ）
- `0x11`: OLED同期（表示コマンド）

### メッセージフォーマット
```
[Header: 0xFF][Type: 1byte][Length: 1byte][Data: N bytes][Checksum: 1byte]
```

チェックサムは単純なXORで計算されます。

## ビルド

### マスター
```bash
cd examples/split-bidirectional/master
tinygo build -target=pico -o master.uf2 .
cp master.uf2 /Volumes/RPI-RP2/
```

### スレーブ
```bash
cd examples/split-bidirectional/slave
tinygo build -target=pico -o slave.uf2 .
cp slave.uf2 /Volumes/RPI-RP2/
```

## 使い方

1. マスターとスレーブをUARTで接続
2. 両方に電源を投入
3. マスター側がUSBホストに接続される
4. キー入力、ジョイスティック、LEDが両側で動作
5. OLED表示が両側に同期される

## キーマップ

### レイヤー 0（デフォルト）

#### マスター側（左半分）
```
Q  W  E
A  S  D
Z  X  C
Ctrl MO(1) Space
```

#### スレーブ側（右半分）
```
R  T  Y
F  G  H
V  B  N
Enter MO(1) Shift
```

### レイヤー 1（MO(1)押下中）

#### マスター側
```
1  2  3
4  5  6
7  8  9
-  -  0
```

#### スレーブ側
```
F1  F2  F3
F4  F5  F6
F7  F8  F9
F10 -  F11
```

## LEDエフェクト

マスター側で設定したLEDエフェクトが自動的にスレーブ側に同期されます。

- **デフォルト**: Rainbow effect（速度5）
- 両側の12個のLEDが同じパターンで表示される

エフェクトを変更するには、master/main.goの以下の行を編集:
```go
kb.SetLEDEffect(led.NewRainbowEffect(5))
```

他のエフェクト例:
```go
// 静的な青色
kb.SetLEDEffect(led.NewStaticEffect(color.RGBA{R: 0, G: 0, B: 255, A: 255}))

// ブリージング (緑)
kb.SetLEDEffect(led.NewBreathingEffect(color.RGBA{R: 0, G: 255, B: 0, A: 255}, 10))
```

## OLEDディスプレイ

マスター側のOLED表示がスレーブ側にも同期されます:

- **レイヤー表示**: 現在のレイヤー番号
- **ジョイスティック値**: X軸（マスター）とY軸（スレーブ）の値
- **Lock状態**: Caps Lock, Num Lock

## 同期レート

- **LED同期**: 1000Hz（USBレポートレート毎）
- **OLED同期**: 1000Hz（USBレポートレート毎）
- **マトリックス同期**: 1000Hz（スキャンインターバル毎）

## パフォーマンス

- **スキャンレート**: 1000Hz（1ms間隔）
- **UART速度**: 460800 bps
- **レイテンシー**: ~2ms（マスター側）、~3-4ms（スレーブ側）

## トラブルシューティング

### スレーブが応答しない
- UART接続を確認（TX/RX、GND）
- ボーレートが一致しているか確認（460800 bps）
- タイムアウト設定を確認（デフォルト5ms）

### LEDが同期しない
- スレーブ側のLED設定を確認
- LED数が一致しているか確認（両側12個）
- UART通信が正常か確認

### OLEDが同期しない
- I2Cアドレスを確認（0x3C）
- スレーブ側のOLED初期化を確認
- UART受信コールバックが設定されているか確認

## Phase 2 M3 ステータス

✅ マスター → スレーブ LED同期
✅ マスター → スレーブ OLED同期
✅ 双方向UARTプロトコル
✅ コールバックベースの実装
✅ 460800 bps高速通信

## 今後の予定

Phase 2 M4では以下の機能が追加予定:
- タップ検出 (TT/LT keycodes)
- ジョイスティックキャリブレーションモード
- より高度なOLED同期（グラフィック描画）

## 備考

- スレーブ側のkeymapは使用されません（マスター側で処理）
- スレーブ側のHIDデバイスは初期化されません
- LED/OLEDはマスター側から制御されます
- ジョイスティックは両側で独立して動作します（2軸モード）
