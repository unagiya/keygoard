# keygoard

**TinyGoベースのゲーミング自作キーボード用ファームウェアフレームワーク**

RP2040/RP2350マイコン向けに構築された、モダンで高性能なキーボードファームウェアフレームワークです。ジョイスティック、OLEDディスプレイ、分割キーボード構成をサポートし、特にゲーミングキーボード向けに設計されています。

[English README](README.md)

## 特徴

- **マトリックススキャン**: COL2ROW/ROW2COL対応、設定可能なデバウンス
- **高性能**: 1msスキャン間隔、1000Hz USBレポートレート（ゲーミング最適化）
- **レイヤーシステム**: 16層、MO、TG、TT、LT操作対応
- **USB HID複合デバイス**: 6KRO/NKROキーボード + 32ボタンゲームパッド
- **アナログジョイスティック**: 2軸、オーバーサンプリング、キャリブレーション、デッドゾーン対応
- **OLEDディスプレイ**: SSD1306 128x32/64対応、両側表示
- **分割キーボード**: UART双方向通信（460800bps）
- **ロータリーエンコーダー**: 最大2個、クリック対応
- **RGB LED**: WS2812、リアクティブエフェクト（KeyPress、FadeOut、Ripple、Rainbow、Breathing）
- **マクロ**: キーシーケンスマクロ、遅延サポート（最大16マクロ）
- **コンボキー**: 複数キー同時押し検出（最大16コンボ）
- **設定永続化**: Flashベースのキーマップ・設定保存
- **簡単セットアップ**: プロジェクト生成用CLIツール

## クイックスタート

### インストール

```bash
# CLIツールをインストール
go install github.com/unagiya/keygoard/cmd/keygoard@latest

# 新規キーボードプロジェクトを作成
keygoard init my-keyboard
cd my-keyboard
```

### 設定

`config.go` を編集してハードウェアに合わせます：

```go
var boardConfig = engine.BoardConfig{
    Rows: []machine.Pin{
        machine.GP0, machine.GP1, machine.GP2, machine.GP3,
    },
    Cols: []machine.Pin{
        machine.GP4, machine.GP5, machine.GP6, machine.GP7,
    },
    MatrixType: engine.COL2ROW,

    Joystick: &joystick.Config{
        PinX: machine.ADC0,
        PinY: machine.ADC1,
    },

    OLED: &oled.Config{
        I2C:    machine.I2C0,
        Width:  128,
        Height: 32,
    },
}
```

`keymap.go` を編集してレイアウトを定義：

```go
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー 0
    {
        {keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
        {keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.MO(1)},
    },
    // レイヤー 1
    {
        {keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4},
        {keycode.KC_F1, keycode.KC_F2, keycode.KC_F3, keycode.KC_TRNS},
    },
}
```

### ビルドとフラッシュ

```bash
# ファームウェアをビルド
tinygo build -target=pico -o firmware.uf2 .

# RP2040にフラッシュ
# 1. BOOTSELボタンを押しながらUSBに接続
# 2. ファームウェアをデバイスにコピー
cp firmware.uf2 /Volumes/RPI-RP2/
```

## アーキテクチャ

```
keygoard/
├── engine/              # コアキーボードエンジン
├── matrix/              # マトリックススキャン＆デバウンス
├── keycode/             # キーコード定義
├── hid/                 # USB HIDレポート
├── peripheral/          # ジョイスティック、OLED、LED、エンコーダーサポート
├── split/               # 分割キーボード通信
├── storage/             # Flashベース永続化ストレージ
└── cmd/keygoard/        # CLIツール
```

## ハードウェアサポート

### マイコン
- RP2040（Raspberry Pi Pico）
- RP2350
- Waveshare RP2040-Zero

### マトリックス
- 片側最大16x8（128キー）
- COL2ROWとROW2COL対応
- 3msデバウンス（設定可能）

### 周辺機器
- **ジョイスティック**: 2軸（分割キーボードで4軸）
- **OLED**: SSD1306（I2C）、128x32または128x64
- **分割**: UART 460800bps（双方向）
- **ロータリーエンコーダー**: 最大2個
- **RGB LED**: WS2812、最大128個

## サンプル

`examples/` ディレクトリに完全動作するサンプルがあります：

- `simple-gamepad` - ジョイスティック・OLED付き4x6キーボード
- `encoder-test` - ロータリーエンコーダーデモ
- `led-test` - RGB LEDエフェクトショーケース
- `split-bidirectional` - 双方向分割キーボード
- `reactive-led` - リアクティブLEDエフェクト
- `split-led` - 両側LED制御
- `nkro-test` - NKRO機能
- `macro-test` - マクロ記録・再生
- `combo-test` - コンボキー
- `persistent-keymap` - Flashベースキーマップ永続化
- `waveshare-rp2040-zero` - zero-kb02ボード（3×4マトリックス、エンコーダ、ジョイスティック、OLED）

キーコードリファレンス、分割キーボード詳細設定、パフォーマンスチューニング、GPIO制約については[詳細ドキュメント](docs/README_ja.md)を参照してください。

## 開発

### 必要環境
- Go 1.21以降
- TinyGo 0.30.0以降
- RP2040開発ボード

### ソースからビルド
```bash
git clone https://github.com/unagiya/keygoard.git
cd keygoard
go mod download
```

### サンプルの実行
```bash
cd examples/simple-gamepad
tinygo build -target=pico -o firmware.uf2 .
```

## ロードマップ

詳細は [ROADMAP.md](ROADMAP.md) を参照してください。

## 設計思想

1. **ゲーミングファースト**: 低レイテンシと高性能を優先
2. **型安全**: Goの型システムを活用した信頼性の高いファームウェア
3. **メモリ効率**: 128KB以下のRAM使用を目標
4. **拡張性**: 新機能や周辺機器の追加が容易

## ライセンス

MITライセンス - 詳細は LICENSE ファイルを参照

## コントリビューション

貢献を歓迎します！ガイドラインは [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。

## 謝辞

- [TinyGo](https://tinygo.org/) で構築
- QMK、KMK、ZMKファームウェアプロジェクトからインスピレーション
- 周辺機器ドライバに [tinygo-drivers](https://github.com/tinygo-org/drivers) を使用

## 関連ドキュメント

- [サンプル実装ガイド](examples/simple-gamepad/README_ja.md)
- [詳細ドキュメント](docs/README_ja.md)
