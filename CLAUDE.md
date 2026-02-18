# CLAUDE.md

## プロジェクト概要
- TinyGoベースのカスタムゲーミングキーボードファームウェアフレームワーク
- ターゲット: RP2040 / RP2350
- 現在のバージョン: v0.1.0（ベータ版開発完了）

## ビルド・検証コマンド
- `make build` - 全パッケージのビルド確認
- `make test` - テスト実行
- `make fmt` - コードフォーマット
- `make lint` - Lint実行（golangci-lint）
- `make build-example` - サンプルのTinyGoビルド（RP2040向け）

## TinyGo制約
- `reflect`パッケージは使用禁止
- `fmt.Sprintf`等のヒープ割り当てが大きい関数は極力避ける
- goroutineの使用は最小限にする（スケジューラのオーバーヘッド）
- メモリ使用量128KB以下を目標とする
- 標準ライブラリの一部は使用不可（net, os等）。TinyGoの対応状況を確認すること
- サードパーティパッケージの新規追加は事前にユーザーへ確認すること

## アーキテクチャ
- パッケージ依存方向:
  - `keycode` → 他パッケージに依存しない（最下層）
  - `matrix`, `hid`, `split`, `peripheral/` → `keycode`に依存可
  - `engine` → 全パッケージを統合する（最上層）
- 循環依存は禁止
- 新しい周辺機器は `peripheral/` 配下にパッケージを作成する
- `cmd/`配下はCLIツール専用

## コーディング規約
- `goimports`でフォーマット（TinyGo固有パッケージのimportは手動で確認）
- コメントは日本語で記述
- マジックナンバーは定数化する
- エラーは`errors.New`で定義し、各パッケージの`errors.go`にまとめる
- グローバル変数は組み込み上やむを得ない場合のみ許容

## コミット規約
- Conventional Commits形式を使用
- 日本語で記述
- 形式: `<type>: <説明>`
  - `feat:` 新機能
  - `fix:` バグ修正
  - `docs:` ドキュメント
  - `refactor:` リファクタリング
  - `test:` テスト追加・修正
  - `chore:` ビルド・設定等の雑務

## ブランチ戦略
- `main` - 安定版（Phase完了時にマージ）
- `phase{N}` - Phase開発ブランチ（例: `phase3`）
- `feature/xxx` - 機能ブランチ（phase{N}から分岐・マージ）
- フロー: `feature/xxx` → `phase{N}` → `main`

## テスト方針
- ハードウェア非依存のパッケージにはユニットテストを書く
  - 対象: `keycode`, `hid/report`, `split/protocol`, `matrix/debounce`, `engine/layer`, `engine/tap`
- テストファイルは対象と同じパッケージ内に `*_test.go` で配置
- `go test ./...` で実行可能な状態を維持する

## ドキュメント方針
- README.md（英語）とREADME_ja.md（日本語）を両方維持する
- 新機能追加時はROADMAP.mdの進捗も更新する
