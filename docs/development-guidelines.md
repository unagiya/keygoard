# 開発ガイドライン

---

## コーディング規約

### 基本方針

- フォーマット: `goimports`（import の整理も兼ねる）
- コメント: 日本語で記述する
- マジックナンバーは定数化する
- グローバル変数: 組み込み上やむを得ない場合のみ許容。必ずコメントで理由を書く

### エラー定義

- `errors.New` で定義し、各パッケージの `errors.go` に集約する
- エラーメッセージは `パッケージ名: 説明` の形式にする

```go
// errors.go
var ErrKeyOverflow = errors.New("keygoard: key overflow")
```

### コメント形式

パッケージ・型・関数には GoDoc 形式のコメントを書く。

```go
// Package matrix はマトリクススキャンとデバウンス処理を実装します。
package matrix

// Debouncer はキーごとのデバウンス状態を保持します。
// machine パッケージに依存しないため、標準 Go でテスト可能です。
type Debouncer struct { ... }

// Update は生スキャン結果を受け取り、デバウンス後の状態と変化フラグを返します。
func (d *Debouncer) Update(...) { ... }
```

### TinyGo 制約（必須）

| 禁止事項 | 理由 |
|---|---|
| `reflect` パッケージ | TinyGo でランタイムエラーになる |
| `fmt.Sprintf` / `fmt.Printf` | ヒープ割り当てが大きい |
| `net` / `os` など TinyGo 未対応パッケージ | TinyGo がサポートしない |
| goroutine の多用 | スケジューラオーバーヘッドが発生する |

- デバッグ出力は `println()` のみ使用する
- メインループは 1 goroutine で動かす

### ホットパス（スキャンループ内）の制約

| 禁止事項 | 代替手段 |
|---|---|
| `append` による再割り当て | 固定サイズ配列を初期化時に確保して再利用 |
| スライスの動的生成 | スライスの初期確保は初期化時に行う |
| 文字列連結 | 固定バイト列を使う |

---

## 命名規則

Go の標準的な命名規則に従う。以下はプロジェクト固有のルールを補足する。

### パッケージ名

- 小文字・単数形・1 単語（例: `matrix`, `keycode`, `engine`）
- アンダースコア・ハイフン不使用

### 型・構造体

- PascalCase（例: `Debouncer`, `Keycode`）
- 略語は大文字で統一（例: `HID`, `USB`, `GPIO`）

### 定数

| 種別 | 形式 | 例 |
|---|---|---|
| 公開定数 | PascalCase | `ModLeftCtrl`, `Enter`, `None` |
| 非公開定数 | camelCase | `debounceThr` |

キーコード定数はキー名をそのまま使う。`KC_` プレフィックスは付けない。

```go
// 良い
const A Keycode = 4 | 0xF000

// 悪い
const KC_A Keycode = 4 | 0xF000
```

### 変数・フィールド

| 種別 | 形式 | 例 |
|---|---|---|
| 公開変数 | PascalCase | `ErrKeyOverflow` |
| 非公開変数・フィールド | camelCase | `debounced`, `counter` |

### 関数・メソッド

- 公開: PascalCase（例: `Update`, `Init`, `Scan`）
- 非公開: camelCase

### エラー変数

- `Err` プレフィックス + PascalCase（例: `ErrKeyOverflow`）

### ファイル名

- スネークケース（例: `debounce.go`, `keycode_test.go`）
- machine 依存ファイルはモジュール名をそのまま使う（例: `matrix.go`, `hid.go`）
- テストファイルは `対象ファイル名_test.go`

---

## スタイリング規約

### Import グループ

`goimports` が自動整理するが、TinyGo 固有パッケージを含む場合は手動で確認する。
グループ順は以下の通り（空行で区切る）。

```go
import (
    // 1. 標準ライブラリ
    "errors"

    // 2. TinyGo / machine（//go:build tinygo タグ付きファイルのみ）
    "machine"
    "machine/usb/hid"

    // 3. サードパーティ
    "tinygo.org/x/drivers/..."

    // 4. 内部パッケージ
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/matrix"
)
```

### ビルドタグ

`machine` パッケージに依存するファイルには先頭に `//go:build tinygo` タグを付ける。
タグとパッケージ宣言の間には空行を 1 つ入れる。

```go
//go:build tinygo

package matrix
```

タグがないファイルは標準 `go test` でビルド・テスト可能にする。

| ファイル | タグ | テスト方法 |
|---|---|---|
| `matrix.go` | `//go:build tinygo` | TinyGo のみ |
| `debounce.go` | なし | `go test ./matrix/...` |
| `hid.go` | `//go:build tinygo` | TinyGo のみ |
| `keycode.go` | なし | `go test ./keycode/...` |

### 定数グループ

関連する定数は `const ( ... )` でグループ化し、コメントで区切る。

```go
// 修飾キー（modifier bitmap | 0xE000）
const (
    ModLeftCtrl  Keycode = 0x01 | 0xE000
    ModLeftShift Keycode = 0x02 | 0xE000
    // ...
)

// 通常キー（HID Usage Page 7 usage code | 0xF000）
const (
    A Keycode = 4 | 0xF000
    B Keycode = 5 | 0xF000
    // ...
)
```

---

## Git 規約

### コミットメッセージ

Conventional Commits 形式、日本語で記述する。

```
<type>: <説明>
```

| type | 用途 | 例 |
|---|---|---|
| `feat:` | 新機能 | `feat: マトリクススキャン実装` |
| `fix:` | バグ修正 | `fix: デバウンス閾値の境界条件を修正` |
| `docs:` | ドキュメント | `docs: matrix.md にスキャン仕様を追記` |
| `refactor:` | リファクタリング | `refactor: engine のループ変数を定数に変更` |
| `test:` | テスト追加・修正 | `test: デバウンスのリリーステストを追加` |
| `chore:` | ビルド・設定等 | `chore: Makefile に lint ターゲットを追加` |

### ブランチ戦略

```
main
└── phase{N}
    ├── feature/xxx
    └── feature/yyy
```

| ブランチ | 用途 | 分岐元 | マージ先 |
|---|---|---|---|
| `main` | 実機確認済みの安定版 | — | — |
| `phase{N}` | Phase 開発ブランチ | `main` | `main` |
| `feature/xxx` | トピックブランチ（1 機能） | `phase{N}` | `phase{N}` |

`feature/xxx` → `phase{N}` → `main` の順でマージする。
`main` へのマージは実機確認完了後のみ。

### バージョニング

`v0.{マイナー}.{パッチ}` 形式。

| タグ | 意味 |
|---|---|
| `v0.(N-1).1`, `v0.(N-1).2` … | Phase N 内の開発チェックポイント |
| `v0.N.0` | Phase N 完了（実機確認・`main` マージ） |

- フェーズ完了時にマイナーバージョンが上がる
- バージョンは常に単調増加する
- `v0.N.0` タグは `main` マージ後に打つ

**例:**

```
v0.0.1  ← Phase 1: マトリクススキャン完成チェックポイント
v0.0.2  ← Phase 1: HID 認識チェックポイント
v0.1.0  ← Phase 1 完了（実機確認・main マージ）
```
