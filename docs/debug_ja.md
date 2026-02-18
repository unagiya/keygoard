# デバッグ・パフォーマンス

keygoardはデバッグログとエラー統計機能を提供し、開発・トラブルシューティングを支援します。

[English](debug.md) | 日本語

## 概要

- DebugLoggerによる軽量デバッグログ出力（USB CDCシリアル）
- Stats構造体によるエラーカウンタと統計情報
- パフォーマンス最適化のガイドライン

## デバッグモード

### 有効化

```go
cfg := &engine.BoardConfig{
    Debug: true,  // デバッグ出力を有効にする
    // ...
}
```

`Debug: true`を設定すると、USB CDCシリアル（`println()`）経由で以下のログが出力されます。

### ログ出力内容

#### 起動時ログ
- マトリクスサイズ（行×列）
- 有効な周辺機器（ジョイスティック、OLED、エンコーダー、LED）
- Split設定（マスター/スレーブ、ボーレート）

#### 実行時ログ
- キーイベント（行、列、押下/離上）
- レイヤー変更（変更元→変更先）
- Split接続状態（接続/切断）

### DebugLogger

エンジン内部で使用される軽量ロガーです。`enabled=false`時は全メソッドがノーオペレーションとなり、パフォーマンスへの影響はありません。

```go
// エンジン内部での使用例（直接のインスタンス化は不要）
logger.Log("engine", "キーボード初期化完了")
logger.LogError("hid", err)
logger.LogKeyEvent(row, col, pressed)
logger.LogLayerChange(fromLayer, toLayer)
logger.LogSplitStatus(connected)
```

## エラー統計

### Stats構造体

```go
type Stats struct {
    SplitRxErrors    uint16  // Split受信エラー回数
    SplitTxErrors    uint16  // Split送信エラー回数
    HIDSendErrors    uint16  // HIDレポート送信エラー回数
    LEDWriteErrors   uint16  // LED書き込みエラー回数
    ScanCount        uint32  // スキャン回数（累積）
    SplitDisconnects uint16  // スレーブ切断回数
}
```

### 統計情報の取得

```go
kb, _ := engine.NewKeyboard(cfg, keymap)

// キーボード実行中に取得
stats := kb.GetStats()

// 各カウンタを確認
println("スキャン回数:", stats.ScanCount)
println("Split受信エラー:", stats.SplitRxErrors)
println("HID送信エラー:", stats.HIDSendErrors)
```

## パフォーマンス最適化tips

keygoardでは以下の最適化が実装されています。開発時の参考にしてください。

### ゼロアロケーション設計

ホットパス（毎スキャンサイクルで実行されるコード）ではヒープ割り当てを95%削減しています。

- **事前割り当てバッファ**: マトリクススキャン、デバウンス、LED、Splitプロトコルの各バッファを起動時に確保
- **ゼロアロケーションAPI**: `EncodeInto()` / `DecodeMatrixStateInto()` は既存バッファに書き込み
- **固定配列バッファ**: TapDetectorは固定長配列を使用

### OLED最適化

- 変更検出により、内容が変わった場合のみI2C転送
- 固定バイト列化でシリアライズのヒープ割り当てを回避

### TinyGoでのパフォーマンスガイドライン

| 避けるべき操作 | 推奨される代替 |
|---|---|
| `fmt.Sprintf()` | `println()` / 固定文字列結合 |
| `append()` による動的拡張 | 事前割り当てスライス |
| `map` のホットパスでの使用 | 配列 / 固定サイズ構造体 |
| goroutineの過度な使用 | メインループでの逐次処理 |
| インターフェースの過度なboxing | 具体型の直接使用 |

### メモリ目標

- 全体のメモリ使用量: 128KB以下
- ホットパスのヒープ割り当て: ゼロを目標

## 使用例

```go
package main

import (
    "github.com/unagiya/keygoard/engine"
)

func main() {
    cfg := &engine.BoardConfig{
        Debug: true,  // デバッグ出力を有効化
        // ... その他の設定 ...
    }

    kb, _ := engine.NewKeyboard(cfg, keymap)

    // OLEDに統計情報を表示するコールバック等で使用
    // stats := kb.GetStats()

    kb.Run()
}
```

## 制約・制限

- デバッグログはUSB CDCシリアル経由で出力されるため、シリアルモニタが必要
- デバッグモード有効時はわずかにパフォーマンスが低下する（`println()`のオーバーヘッド）
- Stats構造体のカウンタは`uint16`のため、65535を超えるとオーバーフローする
- ScanCountは`uint32`のため、約49日間の連続動作でオーバーフローする（1000Hz時）
