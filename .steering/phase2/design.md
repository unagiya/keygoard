# Phase 2 設計

## 実装アプローチ

Phase 1 の `engine.New()` + `Init()` + 手動ループ を `engine.New()` + `Run()` に統合し、
USB Product Name を含む設定を `Config` 構造体で受け取る形に変更する。

---

## API 変更

### Before（Phase 1）

```go
// examples/zero-kb02/main.go
func main() {
    scanner := matrix.New(matrixCols, matrixRows)
    kb := engine.New(scanner, defaultKeymap)
    kb.Init()

    for {
        kb.Tick()
        time.Sleep(1 * time.Millisecond)
    }
}
```

### After（Phase 2）

```go
// examples/zero-kb02/main.go
func main() {
    kb := engine.New(&engine.Config{
        Scanner:     matrix.New(matrixCols, matrixRows),
        Keymap:      defaultKeymap,
        ProductName: "zero-kb02",
    })
    kb.Run()
}
```

---

## 変更するコンポーネント

### engine パッケージ

#### config.go（新規）

`Config` 構造体を定義する。`//go:build tinygo` タグ不要（`machine` に依存しない）。

```go
package engine

import "github.com/unagiya/keygoard/matrix"

// Config はキーボードエンジンの設定です。
type Config struct {
    // Scanner はマトリクススキャナーです。
    Scanner *matrix.Scanner

    // Keymap はキーマップです。
    Keymap *Keymap

    // ProductName は USB デバイスとして OS に表示される名前です。
    // 空文字列の場合はデフォルト値 "keygoard" を使用します。
    // ASCII のみ、最大 126 文字。
    ProductName string
}
```

注意: `Config` に `Scanner *matrix.Scanner` を含めるが、`matrix.Scanner` 自体は
`matrix/matrix.go`（`//go:build tinygo`）で定義されている。
`Config` 構造体はフィールド型としてポインタを持つだけなので、
`config.go` 自体に `//go:build tinygo` タグは不要
（ただし `matrix.Scanner` 型が標準 Go ビルドで未定義の場合はタグが必要になる可能性がある。
実装時に確認する）。

#### keyboard.go（変更）

`//go:build tinygo` タグ付きのまま。

```go
// New は Config からキーボードエンジンを生成します。
func New(cfg *Config) *Keyboard

// Run はキーボードを起動します。
// ハードウェア初期化・USB エニュメレーション待機の後、スキャンループに入ります。
// この関数は戻りません。
func (kb *Keyboard) Run()
```

**`Run()` の内部処理:**

```
1. USB Product Name を設定（usb.Product に代入）
2. scanner.Init()
3. USB エニュメレーション完了まで待機（500ms）
4. 無限ループ:
   a. kb.Tick()
   b. time.Sleep(1ms)
```

**`Init()` と `Tick()` について:**
- `Init()` は `Run()` に統合するため、公開メソッドとしては削除する
- `Tick()` は内部メソッド（非公開）に変更する

#### keymap.go（変更なし）

`Keymap` 構造体はそのまま維持。

#### errors.go（変更なし）

### examples/zero-kb02/

#### main.go（変更）

`for` ループ・`time.Sleep` を削除し、`engine.New()` + `kb.Run()` に変更する。
`_ "machine/usb/hid/keyboard"` の blank import は引き続き必要（HID ハンドラ登録のため）。

#### keymap.go（変更なし）

#### config.go（変更）

`matrixCols` / `matrixRows` の定義は維持。
Product Name は `main.go` 側の `engine.Config` で設定するため、config.go への変更は不要。

---

## USB Product Name の実装詳細

TinyGo は `machine/usb` パッケージに以下のパッケージ変数を提供している:

```go
package usb

var (
    Product      string
    Manufacturer string
    // ...
)
```

- ホストが USB String Descriptor を要求した時点で遅延評価される
- `Run()` の冒頭で `usb.Product` に値を代入すれば、エニュメレーション時に反映される
- ASCII のみ対応、最大 126 文字
- 空文字列の場合はボードのデフォルト値（例: "Pico"）が使われる

**フレームワーク側のデフォルト値:** `"keygoard"`
（`Config.ProductName` が空の場合に適用）

---

## データフロー（変更なし）

Phase 1 と同じ。API の呼び出し方が変わるだけで、内部のスキャン→デバウンス→HID 送信の流れは変更なし。

```
物理キー押下
    ↓
matrix.Scan()        → rawState [3][4]bool
    ↓
debounce.Update()    → stableState [3][4]bool
    ↓
keymap.Layer0[][]    → keycodes
    ↓
hid Down/Up          → USB HID レポート
    ↓
OS（macOS）
```

---

## 影響範囲

| ファイル | 変更内容 |
|---|---|
| `engine/config.go` | 新規作成（`Config` 構造体） |
| `engine/keyboard.go` | `New` のシグネチャ変更、`Run()` 追加、`Init()` 非公開化、`Tick()` 非公開化 |
| `examples/zero-kb02/main.go` | 新 API に対応（`for` ループ削除） |
| `docs/packages/engine.md` | 新規作成（engine パッケージ仕様書） |

**変更なし:** `keycode/`、`matrix/`、`engine/keymap.go`、`engine/errors.go`
