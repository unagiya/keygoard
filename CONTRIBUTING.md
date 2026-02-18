# Contributing to keygoard

keygoardへの貢献を検討していただき、ありがとうございます！

[English](CONTRIBUTING.md) | [日本語](#日本語)

## 日本語

### 始める前に

1. [Issue](https://github.com/unagiya/keygoard/issues)で既存の提案や問題を確認
2. 新しい機能を追加する場合は、まずIssueで議論することをお勧めします
3. [Code of Conduct](CODE_OF_CONDUCT.md) を確認（作成予定）

### 開発環境のセットアップ

#### 必要なツール

- Go 1.21以降
- TinyGo 0.30.0以降
- Make（オプション、推奨）
- RP2040開発ボード（実機テスト用）

#### セットアップ手順

```bash
# リポジトリをクローン
git clone https://github.com/unagiya/keygoard.git
cd keygoard

# 開発環境をセットアップ（Makeを使用）
make setup

# または手動で
go mod download
go install ./cmd/keygoard
```

### 開発ワークフロー

#### 1. ブランチを作成

```bash
git checkout -b feature/your-feature-name
# または
git checkout -b fix/your-bug-fix
```

ブランチ命名規則：
- `feature/` - 新機能
- `fix/` - バグ修正
- `docs/` - ドキュメントのみの変更
- `refactor/` - リファクタリング
- `test/` - テスト追加・修正

#### 2. コードを書く

```bash
# コードを編集
# ...

# フォーマット
make fmt
# または
go fmt ./...

# ビルド確認
make build

# テスト実行
make test
```

#### 3. コミット

コミットメッセージの規約：

```
<type>: <subject>

<body>

<footer>
```

**Type:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット等）
- `refactor`: バグ修正でも機能追加でもないコード変更
- `test`: テストの追加・修正
- `chore`: ビルドプロセスやツールの変更

**例:**
```
feat: add rotary encoder support

Add basic rotary encoder functionality for Phase 2.
Supports up to 2 encoders with clockwise/counter-clockwise events.

Closes #42
```

#### 4. プッシュとPull Request

```bash
# プッシュ
git push origin feature/your-feature-name

# GitHubでPull Requestを作成
# - わかりやすいタイトルと説明を書く
# - 関連するIssueをリンク
# - スクリーンショットやテスト結果を添付（該当する場合）
```

### コーディング規約

#### Goコード

- `gofmt` でフォーマット済みであること
- `golangci-lint` のチェックをパスすること
- 公開API には適切なドキュメントコメントを付けること
- テストを書くこと（可能な限り）

#### ファイル構成

```go
// パッケージコメント: パッケージの目的を説明
package example

import (
    // 標準ライブラリ
    "time"

    // サードパーティ
    "tinygo.org/x/drivers/ssd1306"

    // 内部パッケージ
    "github.com/unagiya/keygoard/keycode"
)

// 公開型・関数にはドキュメントコメントを付ける
// Type は何かを表します。
type Type struct {
    Field string // フィールドコメント
}

// NewType は新しいTypeを作成します。
func NewType() *Type {
    return &Type{}
}
```

#### 命名規則

- 短く、明確な名前を使う
- Goの一般的な慣習に従う
- 略語は大文字（`ID`, `HTTP`, `USB`）
- プライベートなヘルパー関数は小文字で始める

### テスト

```bash
# 全テストを実行
make test

# 特定のパッケージをテスト
go test ./keycode

# カバレッジ付き
go test -cover ./...
```

テストファイルの命名：`*_test.go`

### ドキュメント

- 新機能を追加したら、対応するドキュメントも更新する
- 日本語と英語の両方を更新することが望ましい
- サンプルコードを含める

### Pull Requestのチェックリスト

Pull Requestを出す前に確認：

- [ ] コードがフォーマットされている（`make fmt`）
- [ ] ビルドが通る（`make build`）
- [ ] テストが通る（`make test`）
- [ ] 新機能にはテストが追加されている
- [ ] ドキュメントが更新されている
- [ ] コミットメッセージが規約に従っている
- [ ] 関連するIssueがリンクされている

### Phase別の開発ガイド

全Phaseが完了しています（v0.1.0ベータ版）。詳細は [ROADMAP.md](ROADMAP.md) を参照してください。

現在は安定性の向上、バグ修正、ドキュメント改善、新しい周辺機器対応などの貢献を歓迎しています。

### リリースプロセス

メンテナーのみ：

1. バージョン番号を決定（セマンティックバージョニング）
2. CHANGELOG.md を更新
3. タグを作成
4. GitHub Releaseを作成

### 質問やヘルプ

- [GitHub Discussions](https://github.com/unagiya/keygoard/discussions) - 質問、アイデア共有
- [GitHub Issues](https://github.com/unagiya/keygoard/issues) - バグ報告、機能要望

### ライセンス

貢献されたコードは、プロジェクトと同じMITライセンスの下で提供されます。

---

## English

### Before You Start

1. Check existing [Issues](https://github.com/unagiya/keygoard/issues)
2. For new features, it's recommended to discuss in an Issue first
3. Review the [Code of Conduct](CODE_OF_CONDUCT.md) (coming soon)

### Development Setup

#### Required Tools

- Go 1.21 or later
- TinyGo 0.30.0 or later
- Make (optional, recommended)
- RP2040 development board (for hardware testing)

#### Setup Steps

```bash
# Clone repository
git clone https://github.com/unagiya/keygoard.git
cd keygoard

# Setup development environment (using Make)
make setup

# Or manually
go mod download
go install ./cmd/keygoard
```

### Development Workflow

#### 1. Create Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming convention:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation only
- `refactor/` - Refactoring
- `test/` - Test additions/fixes

#### 2. Write Code

```bash
# Edit code
# ...

# Format
make fmt

# Build check
make build

# Run tests
make test
```

#### 3. Commit

Commit message convention:

```
<type>: <subject>

<body>

<footer>
```

See Japanese section for types and examples.

#### 4. Push and Pull Request

```bash
# Push
git push origin feature/your-feature-name

# Create Pull Request on GitHub
```

### Coding Standards

- Code must be formatted with `gofmt`
- Must pass `golangci-lint` checks
- Public APIs need documentation comments
- Write tests when possible

### Testing

```bash
# Run all tests
make test

# Test specific package
go test ./keycode

# With coverage
go test -cover ./...
```

### Documentation

- Update documentation when adding features
- Include code examples

### Pull Request Checklist

- [ ] Code is formatted (`make fmt`)
- [ ] Build passes (`make build`)
- [ ] Tests pass (`make test`)
- [ ] Tests added for new features
- [ ] Documentation updated
- [ ] Commit messages follow convention
- [ ] Related issues linked

### Questions or Help

- [GitHub Discussions](https://github.com/unagiya/keygoard/discussions) - Questions, ideas
- [GitHub Issues](https://github.com/unagiya/keygoard/issues) - Bug reports, feature requests

### License

Contributed code will be provided under the same MIT License as the project.
