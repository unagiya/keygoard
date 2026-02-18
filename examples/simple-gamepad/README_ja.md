# Simple Gamepad Keyboard サンプル

ゲームパッド機能付きシンプルなキーボードのサンプル実装です。

[English README](README.md)

## 特徴

- 4x6マトリックスキーボード
- アナログジョイスティック（X/Y軸）
- OLED表示（レイヤーとジョイスティック状態）
- ゲーミング最適化タイミング（1msスキャン、3msデバウンス、1000Hz USB）

## 必要なハードウェア

- RP2040ボード（例：Raspberry Pi Pico）
- 4x6キースイッチマトリックス
- アナログジョイスティックモジュール（ADC0/ADC1に接続）
- 128x32 OLEDディスプレイ（I2C）

## ピン配置

### マトリックス
- 行（Row）: GP0, GP1, GP2, GP3
- 列（Col）: GP4, GP5, GP6, GP7, GP8, GP9

### ジョイスティック
- X軸: GP26（ADC0）
- Y軸: GP27（ADC1）

### OLED
- I2C0（デフォルトピン: GP4=SDA, GP5=SCL）
- アドレス: 0x3C（デフォルト）

## 配線図

```
RP2040 Pico
┌─────────────────┐
│                 │
│ GP0 ─── Row 0   │
│ GP1 ─── Row 1   │
│ GP2 ─── Row 2   │
│ GP3 ─── Row 3   │
│                 │
│ GP4 ─── Col 0   │  ┌──────────┐
│ GP5 ─── Col 1   │  │          │
│ GP6 ─── Col 2   │  │  OLED    │
│ GP7 ─── Col 3   │  │ (128x32) │
│ GP8 ─── Col 4   │  │          │
│ GP9 ─── Col 5   │  └──────────┘
│                 │   SDA SCL
│ GP4 ────────────┼────┘  │
│ GP5 ────────────┼───────┘
│                 │
│ GP26(ADC0) ── Joystick X
│ GP27(ADC1) ── Joystick Y
│                 │
└─────────────────┘
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

## レイアウト

### レイヤー 0（デフォルト）
```
ESC  1    2    3    4    5
TAB  Q    W    E    R    T
CTRL A    S    D    F    GP1
SHFT Z    X    C    MO1  SPC
```

- `MO1`: レイヤー1へ一時切り替え（押している間）
- `GP1`: ゲームパッドボタン1

### レイヤー 1（ファンクション）
```
`    F1   F2   F3   F4   F5
---  ---  UP   ---  ---  ---
---  LEFT DOWN RGHT ---  GP2
---  ---  ---  ---  ---  ENT
```

- `---`: 透過（下のレイヤーのキーを使用）
- `GP2`: ゲームパッドボタン2

## 使い方

### キーボードとして
通常のキーボードとして認識されます。QWERTYキーとファンクションキーが使用可能です。

### ゲームパッドとして
- ジョイスティック：アナログ入力として認識
- ボタン：GP1、GP2キーがゲームパッドボタンとして機能
- OS側で「ゲームコントローラー」として認識されます

### OLEDディスプレイ
起動すると以下の情報が表示されます：
- 現在のレイヤー番号
- ジョイスティックのX/Y値
- CapsLock/NumLock状態

## カスタマイズ

### ピン配置を変更する

`config.go` を編集：

```go
Rows: []machine.Pin{
    machine.GP10,  // 好きなピンに変更
    machine.GP11,
    machine.GP12,
    machine.GP13,
},
```

### キーマップを変更する

`keymap.go` を編集：

```go
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー 0: 好きなキーコードに変更
    {
        {keycode.KC_ESC, keycode.KC_1, ...},
        ...
    },
}
```

### タイミングを調整する

`config.go` でパフォーマンス設定を変更：

```go
// より安定した動作（チャタリングが多い場合）
DebounceTime:  5 * time.Millisecond,

// より低遅延（高品質スイッチ向け）
DebounceTime:  2 * time.Millisecond,
```

### ジョイスティックの軸を反転する

`config.go`：

```go
Joystick: &joystick.Config{
    PinX:    machine.ADC0,
    PinY:    machine.ADC1,
    InvertX: true,   // X軸を反転
    InvertY: true,   // Y軸を反転
},
```

## トラブルシューティング

### キーが反応しない

1. マトリックスの配線を確認
2. `MatrixType` を確認（COL2ROW ⇔ ROW2COLを試す）
3. デバウンス時間を増やす（5ms等）

### ジョイスティックが動かない

1. VCC（3.3V）とGNDの接続を確認
2. X/YピンがADC0/ADC1に接続されているか確認
3. ジョイスティックモジュールが正常か確認（テスターで電圧測定）

### OLEDに何も表示されない

1. I2C接続を確認（SDA, SCL）
2. VCC（3.3V）とGNDの接続を確認
3. OLEDアドレスを確認（I2Cスキャナーツール使用）
4. 解像度設定を確認（128x32 または 128x64）

### ゲームパッドとして認識されない

Phase 1ではゲームパッドHID記述子の実装が不完全です。キーボードとしては完全に動作しますが、ゲームパッド機能は今後のアップデートで完全対応予定です。

## 応用例

### ゲーミングマクロパッド
- 左手用デバイスとして使用
- WASDキーとゲームパッドボタンの組み合わせ
- ジョイスティックで視点移動

### クリエイター向けショートカットデバイス
- 動画編集ソフトのショートカットを割り当て
- ジョイスティックでスクラブ/ズーム操作
- レイヤー切り替えでアプリケーション別設定

### ストリーマー用コントロールパネル
- OBS Studio等の配信ソフト制御
- シーン切り替え、ミュート等
- ジョイスティックで音量調整

## さらに学ぶ

- [メインREADME](../../README_ja.md) - フレームワーク全体のドキュメント
- [キーコードリファレンス](../../README_ja.md#キーコードリファレンス)
- [分割キーボードの作り方](../../README_ja.md#分割キーボード)

## ライセンス

MITライセンス - 詳細は [LICENSE](../../LICENSE) ファイルを参照
