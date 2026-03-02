# matrix パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/matrix`（非公開: 外部 import 不可）

## 概要

COL2ROW 方式のマトリクススキャンとデバウンス処理を実装します。`engine.New()` が Config の `ColPins` / `RowPins` から内部で Scanner を生成します。

## ファイル構成

| ファイル          | machine 依存 | 役割                         |
|-------------------|:------------:|------------------------------|
| `const.go`        | なし         | マトリクスサイズ定数         |
| `debounce.go`     | なし         | デバウンスロジック           |
| `debounce_test.go`| なし         | デバウンスユニットテスト     |
| `matrix.go`       | あり         | GPIO スキャン実装            |

## 定数

```go
const (
    RowCount = 3
    ColCount = 4
)
```

`engine` パッケージが `engine.RowCount` / `engine.ColCount` として再エクスポートしています。

## スキャン方式

COL2ROW: 列ピンを 1 本ずつ High にして全行ピンを読み取る。

1. 対象列ピンを High に設定
2. GPIO 安定化待機（10 µs）
3. 全行ピンを読み取る
4. 列ピンを Low に戻す
5. 全列に対して繰り返す

行ピンはプルダウン入力のため、列が High のとき押下キーは High を返す。

## デバウンス

カウンタ方式を採用。

- 閾値 `debounceThr = 5` 回連続して同じ状態を読み取ったとき確定
- スキャン間隔 1 ms × 閾値 5 = 5 ms のデバウンス時間
- 状態変化があった場合のみ `changed = true` を返す

## API

### Scanner

```go
// engine/keyboard.go から呼び出される
s := matrix.New(cols [ColCount]machine.Pin, rows [RowCount]machine.Pin)
s.Init()

// スキャン（メインループから呼び出す）
state, changed := s.Scan()
// state:   [RowCount][ColCount]bool - デバウンス後のキー状態
// changed: 状態変化があった場合 true
```

### Debouncer

```go
d := Debouncer{}
state, changed := d.Update(raw [RowCount][ColCount]bool)
```

- machine 非依存のため標準 Go でテスト可能
- `Scanner.Scan()` 内部で使用される

## テスト

```bash
# デバウンスロジックのユニットテスト（machine 非依存）
make test
```
