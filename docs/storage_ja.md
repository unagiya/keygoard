# 設定永続化

keygoardはFlashメモリへの設定保存をサポートし、キーマップ・設定・マクロデータを電源オフ後も保持できます。

[English](storage.md) | 日本語

## 概要

- Flashメモリへのデータ保存・読み込み・リセット
- キーマップ保存（最大8KB）
- 設定保存（HIDモード、LEDエフェクト、輝度）
- マクロデータ保存（最大1.5KB）
- XORローテーションチェックサムによるデータ整合性検証

## 設定

### BoardConfig

```go
cfg := &engine.BoardConfig{
    Storage: true,  // Flash永続化を有効にする
    // ...
}
```

### Storageの初期化

```go
import (
    "github.com/unagiya/keygoard/storage"
    "machine"
)

// Flashデバイスとオフセットを指定して初期化
store, err := storage.New(machine.Flash, 0x100000)  // オフセットは環境に応じて調整

// Keyboardに設定
kb.SetStorage(store)

// 起動時に設定を読み込み
kb.LoadSettings()
```

## ストレージレイアウト

| 領域 | オフセット | サイズ | 説明 |
|---|---|---|---|
| ヘッダ | 0 | 16B | マジック、バージョン、チェックサム |
| キーマップ | 16 | 最大8,192B | レイヤー×行×列のキーコードデータ |
| 設定 | 8,208 | 最大64B | HIDモード、LED設定等 |
| マクロ | 8,272 | 最大1,536B | 16マクロ×32ステップ×3バイト |
| **合計** | | **9,808B** | **必要セクター数: 3（12KB）** |

### StorageHeader

```go
type StorageHeader struct {
    Magic    [4]byte  // "KBD\x00"
    Version  uint8    // 1
    Reserved [3]byte
    DataSize uint32   // ヘッダを除くデータサイズ
    Checksum uint32   // XORローテーションチェックサム
}
```

## API

### 高レベルAPI（Keyboard経由）

```go
// 設定を保存（キーマップ + HIDモード）
err := kb.SaveSettings()

// 設定を読み込み（HIDモードを復元）
err := kb.LoadSettings()

// 設定をリセット（空データで上書き）
err := kb.ResetSettings()
```

### 低レベルAPI（Storageパッケージ）

```go
// 初期化
store, err := storage.New(flash, offset)

// 個別データの読み書き
store.Read(buf, offset)
store.Write(data, offset)

// ヘッダ操作
header, err := store.ReadHeader()
store.WriteHeader(header)

// 有効性チェック（マジック + チェックサム検証）
valid := store.IsValid()

// 一括保存・読み込み
store.SaveAll(keymapData, configData, macroData)
keymapData, configData, macroData, err := store.LoadAll()

// チェックサム計算
checksum := storage.CalculateChecksum(data)
```

### ConfigData

```go
type ConfigData struct {
    HIDMode       uint8    // 0=6KRO, 1=NKRO
    LEDEffectID   uint8    // 現在のLEDエフェクトID
    LEDBrightness uint8    // LED輝度 0-255
    Reserved      [61]byte
}

// シリアライズ / デシリアライズ
data := storage.SerializeConfig(cfg)
cfg, err := storage.DeserializeConfig(data)
```

### キーマップのシリアライズ

```go
// フォーマット: [layers(1)][rows(1)][cols(1)][reserved(1)][keycode(2B) × L×R×C]
data := storage.SerializeKeymap(layers, numLayers, rows, cols)
layers, numLayers, rows, cols, err := storage.DeserializeKeymap(data)
```

### マクロのシリアライズ

```go
// 各ステップ: Type(1B) + Value(2B) = 3バイト
data := storage.SerializeMacros(macros)
macros, err := storage.DeserializeMacros(data)
```

### BlockDeviceインターフェース

```go
type BlockDevice interface {
    ReadAt(buf []byte, off int64) (int, error)
    WriteAt(p []byte, off int64) (int, error)
    EraseBlocks(start, length int64) error
}
```

`machine.Flash`がこのインターフェースを満たします。

## 使用例

```go
package main

import (
    "machine"

    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/storage"
)

func main() {
    cfg := &engine.BoardConfig{
        Storage: true,
        // ... その他の設定 ...
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // ストレージを初期化して設定
    store, _ := storage.New(machine.Flash, 0x100000)
    kb.SetStorage(store)

    // 保存済み設定を読み込み
    if err := kb.LoadSettings(); err != nil {
        // 初回起動時やデータ破損時はエラーになる（正常動作）
    }

    // キーボード実行
    kb.Run()

    // 設定変更後に保存（例: NKROモード切り替え後）
    // kb.SaveSettings()
}
```

## エラー

| エラー | 説明 |
|---|---|
| `ErrInvalidMagic` | ヘッダのマジックバイトが不正 |
| `ErrInvalidChecksum` | チェックサムが不一致 |
| `ErrDataTooLarge` | データが領域サイズを超過 |
| `ErrNotInitialized` | Storageが初期化されていない |

## 制約・制限

- Flashの書き込みはセクター単位（4KB）で行われ、書き込み前にセクターの消去が必要
- 必要セクター数は3（12KB）
- 頻繁な書き込みはFlashの寿命に影響するため、設定変更時のみ保存すること
- チェックサムはXORローテーション方式（CRCではない）
