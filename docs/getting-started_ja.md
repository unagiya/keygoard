# クイックスタートガイド

このガイドでは、keygoardを使って最初のキーボードファームウェアを作成する手順を説明します。

## 前提条件

### 必要なソフトウェア

1. **Go 1.21以降**
   ```bash
   # バージョン確認
   go version
   ```

2. **TinyGo 0.30.0以降**
   ```bash
   # インストール（macOS）
   brew install tinygo

   # インストール（Linux）
   wget https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_amd64.deb
   sudo dpkg -i tinygo_0.30.0_amd64.deb

   # バージョン確認
   tinygo version
   ```

### 必要なハードウェア

- RP2040開発ボード（Raspberry Pi Pico推奨）
- キースイッチとダイオード
- ブレッドボードまたはPCB
- USBケーブル

## ステップ1: CLIツールのインストール

```bash
go install github.com/unagiya/keygoard/cmd/keygoard@latest
```

インストールが完了したら、CLIツールが使えることを確認：

```bash
keygoard version
# 出力: keygoard v0.1.0
```

## ステップ2: プロジェクト作成

```bash
# 新しいキーボードプロジェクトを作成
keygoard init my-first-keyboard

# プロジェクトディレクトリに移動
cd my-first-keyboard
```

生成されたファイルを確認：

```bash
ls -la
# 出力:
# main.go    - エントリーポイント
# config.go  - ハードウェア設定
# keymap.go  - キーマップ定義
# go.mod     - Go依存関係
```

## ステップ3: ハードウェア設定

`config.go` を開いて、あなたのハードウェアに合わせて編集します。

### 基本的なマトリックス設定

```go
var boardConfig = engine.BoardConfig{
    // 行ピン（出力）
    Rows: []machine.Pin{
        machine.GP0,  // 行0
        machine.GP1,  // 行1
        machine.GP2,  // 行2
        machine.GP3,  // 行3
    },

    // 列ピン（入力・プルアップ）
    Cols: []machine.Pin{
        machine.GP4,  // 列0
        machine.GP5,  // 列1
        machine.GP6,  // 列2
        machine.GP7,  // 列3
    },

    // スキャン方式
    MatrixType: engine.COL2ROW,
}
```

### マトリックススキャン方式の選択

**COL2ROW（推奨）**:
- 行ピンを出力、列ピンを入力として使用
- 最も一般的な構成

**ROW2COL**:
- 列ピンを出力、行ピンを入力として使用
- 一部の市販PCBで使用

どちらを使うか分からない場合は、まずCOL2ROWを試してください。

## ステップ4: キーマップ設定

`keymap.go` を開いて、キー配置を定義します。

### シンプルな4x4マクロパッド

```go
var keymapLayers = [16][][]keycode.Keycode{
    // レイヤー 0: デフォルト
    {
        {keycode.KC_1, keycode.KC_2, keycode.KC_3, keycode.KC_4},
        {keycode.KC_Q, keycode.KC_W, keycode.KC_E, keycode.KC_R},
        {keycode.KC_A, keycode.KC_S, keycode.KC_D, keycode.KC_F},
        {keycode.KC_Z, keycode.KC_X, keycode.KC_C, keycode.MO(1)},
    },

    // レイヤー 1: ファンクション
    {
        {keycode.KC_F1, keycode.KC_F2, keycode.KC_F3, keycode.KC_F4},
        {keycode.KC_F5, keycode.KC_F6, keycode.KC_F7, keycode.KC_F8},
        {keycode.KC_F9, keycode.KC_F10, keycode.KC_F11, keycode.KC_F12},
        {keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS, keycode.KC_TRNS},
    },

    // 残りのレイヤーは未使用
    nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
}
```

### キーコードの説明

- `KC_1`, `KC_A` など: 基本キー
- `MO(1)`: レイヤー1への一時切り替え（押している間）
- `KC_TRNS`: 透過（下のレイヤーのキーを使用）

詳細は [キーコードリファレンス](keycodes_ja.md) を参照。

## ステップ5: ビルド

```bash
# ファームウェアをビルド
tinygo build -target=pico -o firmware.uf2 .
```

ビルドが成功すると `firmware.uf2` ファイルが生成されます。

### ビルドエラーが出た場合

```bash
# 依存関係を更新
go mod tidy

# もう一度ビルド
tinygo build -target=pico -o firmware.uf2 .
```

## ステップ6: フラッシュ

### RP2040をUF2モードにする

1. **BOOTSELボタンを押しながら**USBケーブルをPCに接続
2. `RPI-RP2` という名前のドライブが表示される
3. BOOTSELボタンを離す

### ファームウェアをコピー

```bash
# macOS/Linux
cp firmware.uf2 /Volumes/RPI-RP2/

# Windows
# エクスプローラーでRPI-RP2ドライブを開き、firmware.uf2をドラッグ&ドロップ
```

コピーが完了すると、RP2040が自動的に再起動し、キーボードとして動作を開始します。

## ステップ7: テスト

### キーボード動作確認

1. テキストエディタを開く
2. キーボードのキーを押す
3. 文字が入力されることを確認

### レイヤー切り替え確認

1. MO(1)キーを押しながら他のキーを押す
2. レイヤー1のキーが入力されることを確認

## トラブルシューティング

### キーが反応しない

**原因1: 配線ミス**
- マルチメーターで導通確認
- ダイオードの向きを確認

**原因2: ピン設定が間違っている**
- `config.go` のピン番号を確認
- 実際の配線と一致しているか確認

**原因3: MatrixTypeが逆**
```go
// COL2ROW ⇔ ROW2COL を試す
MatrixType: engine.ROW2COL,
```

### 一部のキーが反応しない

**デバウンス時間を増やす:**
```go
DebounceTime: 5 * time.Millisecond,  // 3ms → 5ms
```

### キーが二重に入力される（チャタリング）

**デバウンス時間を増やす:**
```go
DebounceTime: 10 * time.Millisecond,  // さらに増やす
```

### ビルドできない

**TinyGoのバージョン確認:**
```bash
tinygo version
# 0.30.0以降であることを確認
```

**依存関係の更新:**
```bash
go clean -modcache
go mod tidy
go get -u tinygo.org/x/drivers
```

## 次のステップ

### 機能を追加する

- [ジョイスティックを追加](joystick_ja.md)
- [OLEDディスプレイを追加](oled_ja.md)
- [分割キーボードを作る](split_ja.md)

### キーマップをカスタマイズする

- [レイヤーシステムを理解する](layers_ja.md)
- [キーコード一覧](keycodes_ja.md)

### サンプルを参考にする

- [simple-gamepad サンプル](../examples/simple-gamepad/README_ja.md)

## よくある質問

### Q: 何キーまで対応していますか？

A: 片側最大128キー（16行×8列）まで対応しています。分割キーボードの場合、両側合わせて256キーまで可能です。

### Q: 6KROとは何ですか？

A: 6 Key Rollover の略で、同時に6つのキーまで認識できることを意味します。Phase 4でNKRO（全キー同時押し）に対応予定です。

### Q: 他のマイコンに対応していますか？

A: Phase 1ではRP2040/RP2350のみ対応しています。将来的に他のマイコンへの対応も検討中です。

### Q: QMKから移行できますか？

A: 基本的な機能は同等ですが、キーコード体系が異なります。QMK互換性は提供していませんが、同様の機能を実現できます。

### Q: VIA/Remapに対応していますか？

A: Phase 4で対応予定です。現在はコードで設定する必要があります。

## サポート

問題が解決しない場合は：

- [GitHub Issues](https://github.com/unagiya/keygoard/issues) でバグ報告
- [GitHub Discussions](https://github.com/unagiya/keygoard/discussions) で質問

## ライセンス

MITライセンス
