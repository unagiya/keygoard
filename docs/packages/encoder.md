ok# encoder パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/peripheral/encoder`（非公開: 外部 import 不可）

## 概要

ロータリーエンコーダーの入力処理を提供するパッケージです。2 相エンコーダーの A/B 信号から回転方向を判定し、割り当てられたキーコードを返します。`engine.Peripheral` インターフェースを実装しています。

利用者は `engine.EncoderConfig` 経由で設定し、`engine.New()` が内部でこのパッケージを呼び出します。

## ファイル構成

| ファイル            | machine 依存 | 役割                              |
|---------------------|:------------:|-----------------------------------|
| `quadrature.go`     | なし         | `Decoder` — クアドラチャデコード  |
| `quadrature_test.go`| なし         | Decoder のユニットテスト          |
| `encoder.go`        | あり         | `Encoder` 構造体（GPIO 読み取り） |

## API

### Encoder

```go
// engine/keyboard.go から呼び出される
enc := encoder.New(pinA, pinB machine.Pin, keyCW, keyCCW keycode.Keycode)
```

`engine.Peripheral` インターフェースを実装:

| メソッド               | 動作                                              |
|------------------------|-------------------------------------------------|
| `Init()`               | GPIO ピンをプルアップ入力に設定                   |
| `Tick() Keycode`       | A/B 信号読み取り → 回転判定 → キーコード返却     |
| `OnLayerChange(int)`   | no-op（レイヤー変更に反応しない）                |

### Decoder

machine 非依存のクアドラチャデコーダー。ユニットテスト可能。

```go
type Direction uint8
const (
    DirNone Direction = iota
    DirCW   // 時計回り
    DirCCW  // 反時計回り
)

d := Decoder{}
dir := d.Update(a, b) // A/B 信号の bool 値を渡す
```

**状態遷移テーブル:**

```
CW:  00 → 01 → 11 → 10 → 00
CCW: 00 → 10 → 11 → 01 → 00
```

4×4 のルックアップテーブルで O(1) 判定。

## テスト

```bash
# Decoder のユニットテスト（machine 非依存）
make test
```
