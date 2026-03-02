# layer パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/layer`（非公開: 外部 import 不可）

## 概要

レイヤー解決ロジックを提供するパッケージです。最上位のアクティブレイヤーから下方向に走査し、最初の非 TRNS キーコードを返します。`engine` パッケージの `tick()` から呼び出されます。

## ファイル構成

| ファイル          | machine 依存 | 役割                              |
|-------------------|:------------:|-----------------------------------|
| `resolver.go`     | なし         | `Resolver` — レイヤー解決ロジック |
| `resolver_test.go`| なし         | Resolver のユニットテスト         |

## 定数

```go
const MaxLayers = 4
```

`engine` パッケージが `engine.MaxLayers` として再エクスポートしています。

## API

### Resolver

```go
// engine/keyboard.go から呼び出される
r := layer.NewResolver(layers *[MaxLayers][matrix.RowCount][matrix.ColCount]keycode.Keycode)
```

| メソッド                   | 動作                                                    |
|----------------------------|---------------------------------------------------------|
| `Resolve(row, col)`       | 最上位アクティブレイヤーから走査し、非 TRNS キーコードを返す |
| `Activate(layer)`         | レイヤーを有効にする（レイヤー 0 は変更不可）           |
| `Deactivate(layer)`       | レイヤーを無効にする（レイヤー 0 は変更不可）           |
| `Toggle(layer)`           | レイヤーの有効/無効を反転する（レイヤー 0 は変更不可）  |
| `IsActive(layer) bool`    | レイヤーが有効かどうかを返す                            |

### 解決アルゴリズム

```
Resolve(row, col):
  for layer = MaxLayers-1 downto 0:
    if !active[layer]: continue
    kc = layers[layer][row][col]
    if kc != TRNS: return kc
  return None
```

- レイヤー 0 は常に有効（Deactivate / Toggle 不可）
- `TRNS`（透過キー）は下位レイヤーにフォールバック
- すべて TRNS の場合は `None` を返す

## テスト

```bash
# Resolver のユニットテスト（machine 非依存）
make test
```
