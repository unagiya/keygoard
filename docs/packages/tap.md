# tap パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/tap`（非公開: 外部 import 不可）

## 概要

タップ/ホールド判定を提供するパッケージです。LT（Layer Tap）/ TT（Tap-Toggle）キーコードの押下時間に基づき、タップかホールドかを判定します。`engine` パッケージの `tick()` から呼び出されます。

## ファイル構成

| ファイル          | machine 依存 | 役割                              |
|-------------------|:------------:|-----------------------------------|
| `detector.go`     | なし         | `Detector` — タップ/ホールド判定  |
| `detector_test.go`| なし         | Detector のユニットテスト         |

## 定数

```go
const Threshold = 200 // ticks（= 200ms）
```

## API

### Phase

```go
type Phase uint8

const (
    Idle    Phase = iota // 待機状態
    Pending              // 判定待ち（キー押下直後）
    Holding              // ホールド確定
)
```

### Detector

```go
// engine の Keyboard 構造体にフィールドとして保持される
var d tap.Detector
```

全フィールドは `[RowCount][ColCount]` の固定サイズ配列でヒープ割り当てを回避しています。

| メソッド                          | 動作                                                    |
|-----------------------------------|---------------------------------------------------------|
| `Press(row, col, kc)`            | LT/TT キー押下時に Pending 状態に入る                  |
| `Release(row, col) (bool, Keycode)` | Pending ならタップ判定（`wasPending=true`）を返す    |
| `Advance()`                       | 全 Pending キーのカウンタをインクリメント（毎サイクル）|
| `CheckTimeout(row, col) (bool, Keycode)` | 閾値超過で Holding に遷移し `timedOut=true` を返す |
| `GetPhase(row, col) Phase`       | 指定位置の現在のフェーズを返す                          |
| `Reset(row, col)`                 | 指定位置の状態を Idle に初期化する                      |

### 状態遷移

```
Idle ──[Press(LT/TT)]──→ Pending ──[閾値超過(200ms)]──→ Holding
  ▲                           │                             │
  └──────[Release=タップ]─────┘                             │
  └──────────────────────[Release=ホールド解除]──────────────┘
```

- `Pending` でリリース → タップ（LT: キーコード送信 / TT: レイヤートグル）
- `Holding` でリリース → ホールド解除（レイヤー無効化）

## テスト

```bash
# Detector のユニットテスト（machine 非依存）
make test
```