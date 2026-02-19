# waveshare-rp2040-zero サンプル

[sago35/keyboards](https://github.com/sago35/keyboards) の **zero-kb02** ボード向けサンプルです。

[English README](README.md)

## 対象ハードウェア

- Waveshare RP2040-Zero（MCU）
- 3行4列キーマトリックス（12キー）
- WS2812 RGB LED × 12個（各キーに1:1対応）
- ロータリーエンコーダ × 1個（クリックボタン付き）
- アナログジョイスティック × 1個
- SSD1306 OLED（128×32）

## ピン配置

| 機能 | ピン |
|------|------|
| マトリックス行 | GPIO9, GPIO10, GPIO11 |
| マトリックス列 | GPIO5, GPIO6, GPIO7, GPIO8 |
| WS2812 LED | GPIO1 |
| エンコーダ A | GPIO3 |
| エンコーダ B | GPIO4 |
| エンコーダボタン | GPIO2 |
| ジョイスティック X | GPIO29 (ADC3) |
| ジョイスティック Y | GPIO28 (ADC2) |
| OLED SDA | GPIO12 (I2C0) |
| OLED SCL | GPIO13 (I2C0) |

## キーレイアウト

### レイヤー 0（デフォルト）

```
┌───┬───┬───┬───┐
│ Q │ W │ E │ R │
├───┼───┼───┼───┤
│ A │ S │ D │ F │
├───┼───┼───┼───┤
│ Z │ X │ C │MO1│
└───┴───┴───┴───┘
```

### レイヤー 1（MO1を押しながら）

```
┌─────┬────┬────┬────┐
│ ESC │ F1 │ F2 │ F3 │
├─────┼────┼────┼────┤
│ TAB │ F4 │ F5 │ F6 │
├─────┼────┼────┼────┤
│     │ F7 │ F8 │    │
└─────┴────┴────┴────┘
```

### エンコーダ

| 操作 | キーコード |
|------|-----------|
| 時計回り | KC_PGUP（Page Up）|
| 反時計回り | KC_PGDN（Page Down）|
| クリック | KC_ENT（Enter）|

## LEDエフェクト

起動時はキー押下で白く光る **KeyPressEffect** が有効です。
`main.go` の `SetLEDEffect` を変更することで他のエフェクトに切り替えられます。

```go
// ブリージングエフェクト（青）
kb.SetLEDEffect(led.NewBreathingEffect(color.RGBA{R: 0, G: 0, B: 255, A: 255}, 10))

// レインボーエフェクト
kb.SetLEDEffect(led.NewRainbowEffect(5))
```

## ビルドとフラッシュ

```bash
# ファームウェアをビルド
tinygo build -target=waveshare-rp2040-zero -o firmware-waveshare-zero.uf2 .

# RP2040にフラッシュ
# 1. BOOTSELボタンを押しながらUSBに接続
# 2. ファームウェアをデバイスにコピー
cp firmware-waveshare-zero.uf2 /Volumes/RPI-RP2/
```

## カスタマイズ

`config.go` でピン配置や周辺機器の設定を変更できます。
`keymap.go` でキーレイアウトを編集できます。
