# Encoder Test サンプル

ロータリーエンコーダー機能のテスト用サンプルです（Phase 2）。

[English README](README.md)

## 特徴

- 4x4マトリックスキーボード
- 2つのロータリーエンコーダー（クリックボタン付き）
- OLEDディスプレイ
- ゲーミング最適化タイミング

## 必要なハードウェア

- RP2040ボード（例：Raspberry Pi Pico）
- 4x4キースイッチマトリックス
- 2x ロータリーエンコーダー（EC11等）
- 128x32 OLEDディスプレイ（I2C）

## ピン配置

### マトリックス
- 行（Row）: GP0, GP1, GP2, GP3
- 列（Col）: GP4, GP5, GP6, GP7

### エンコーダー 1（音量制御）
- Pin A: GP8
- Pin B: GP9
- クリック: GP10
- 時計回り（CW）: Volume Up (Keypad +)
- 反時計回り（CCW）: Volume Down (Keypad -)
- クリック: ミュート (Keypad Enter)

### エンコーダー 2（ナビゲーション）
- Pin A: GP11
- Pin B: GP12
- クリック: GP13
- 時計回り（CW）: Down Arrow
- 反時計回り（CCW）: Up Arrow
- クリック: Enter

### OLED
- I2C0（デフォルトピン: GP4=SDA, GP5=SCL）
- アドレス: 0x3C（デフォルト）

## 配線図

```
Rotary Encoder (EC11)
┌─────────┐
│   A  ───┼── GP8/GP11
│   B  ───┼── GP9/GP12
│  GND ───┼── GND
│   C  ───┼── GP10/GP13 (click)
│  C+ ───┼── 3.3V
└─────────┘
```

## ビルド

```bash
tinygo build -target=pico -o firmware.uf2 .
```

## フラッシュ

1. BOOTSELボタンを押しながらRP2040をUSBに接続
2. firmware.uf2をRPI-RP2ドライブにコピー：
   ```bash
   cp firmware.uf2 /Volumes/RPI-RP2/
   ```

## 使い方

### エンコーダー 1（音量制御）
- 時計回りに回す: 音量アップ
- 反時計回りに回す: 音量ダウン
- クリック: ミュート

### エンコーダー 2（ナビゲーション）
- 時計回りに回す: 下矢印
- 反時計回りに回す: 上矢印
- クリック: Enter

## カスタマイズ

`config.go` を編集してエンコーダーのピン割り当てとキーマッピングを変更できます：

```go
Encoders: []*encoder.Config{
    {
        PinA:     machine.GP8,
        PinB:     machine.GP9,
        PinClick: machine.GP10,
        CW:       keycode.KC_PPLS, // 好きなキーに変更
        CCW:      keycode.KC_PMNS,
        Click:    keycode.KC_PENT,
    },
}
```

## トラブルシューティング

### エンコーダーが反応しない
- 配線を確認（A, B, GND）
- config.goのピン設定を確認
- テスターでエンコーダーを確認（抵抗値が変化するはず）

### エンコーダーのステップが飛ぶ
- A-GND間とB-GND間に0.1μFのコンデンサを追加してみる
- 接続の緩みを確認

### 回転方向が逆
- 設定でPinAとPinBを入れ替える

### クリックボタンが動作しない
- PinClickが接続されているか確認
- プルアップ設定を確認（内蔵プルアップが有効化されています）

## ライセンス

MITライセンス
