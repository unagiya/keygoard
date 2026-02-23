# リポジトリ構成

---

## フォルダ・ファイル構成

### 現在の構成

```
keygoard/
├── CLAUDE.md                    ← Claude のプロジェクトメモリ
├── README.md                    ← プロジェクト概要・セットアップ手順（人間向け）
├── ROADMAP.md                   ← フェーズ計画・完了状況
├── Makefile                     ← ビルド・テスト・フラッシュコマンド
├── go.mod                       ← Go モジュール定義（フレームワーク本体）
├── go.work                      ← Go ワークスペース定義（ローカル開発用）
│
├── engine/                      ← 全体統合層（最上位パッケージ）
│   ├── keyboard.go              ← メインループ・HID 送信（//go:build tinygo）
│   ├── keymap.go                ← キーマップ構造体
│   └── errors.go                ← エラー定義
│
├── matrix/                      ← マトリクススキャン・デバウンス
│   ├── const.go                 ← マトリクスサイズ定数
│   ├── matrix.go                ← GPIO スキャン実装（//go:build tinygo）
│   ├── debounce.go              ← デバウンスロジック（machine 非依存）
│   └── debounce_test.go         ← デバウンスユニットテスト
│
├── keycode/                     ← HID キーコード定数定義
│   ├── keycode.go               ← キーコード型・定数
│   └── keycode_test.go          ← ユニットテスト
│
├── examples/                    ← フレームワーク使用例（別モジュール）
│   └── zero-kb02/               ← zero-kb02 向けファームウェアサンプル
│       ├── go.mod               ← モジュール定義（github.com/.../examples/zero-kb02）
│       ├── main.go              ← エントリポイント
│       ├── config.go            ← ピン定義・ボード設定
│       └── keymap.go            ← キーマップ定義
│
├── docs/                        ← プロジェクトドキュメント
│   ├── product-requirements.md  ← プロダクト要求定義
│   ├── functional-design.md     ← 機能設計（システム構成図・API）
│   ├── architecture.md          ← 技術仕様（スタック・制約・パフォーマンス）
│   ├── repository-structure.md  ← 本ファイル
│   └── packages/                ← パッケージ実装レベル仕様
│       ├── keycode.md
│       └── matrix.md
│
└── .steering/                   ← 作業単位ドキュメント（スペック駆動開発）
    └── phase1/                  ← Phase 1 作業ドキュメント
        ├── requirements.md      ← 要求仕様
        ├── design.md            ← 設計
        └── tasklist.md          ← タスクリスト・進捗
```

### 将来の構成（Phase 3 以降）

```
keygoard/
├── ...（上記と同様）
│
├── engine/
│   └── peripheral/              ← 周辺機器インターフェース（Phase 3 で追加）
│       ├── encoder/             ← ロータリーエンコーダー（GP3/GP4）
│       ├── led/                 ← RGB LED WS2812（GP1、12個）
│       ├── oled/                ← OLED SSD1306（GP12/GP13）
│       └── joystick/            ← アナログジョイスティック（Phase 4）
│
├── examples/
│   └── zero-kb02/
│       └── ...
│
└── docs/packages/               ← フェーズ進行に伴い追加
    ├── hid.md                   ← Phase 1
    ├── layer.md                 ← Phase 2
    ├── encoder.md               ← Phase 3
    ├── led.md                   ← Phase 3
    ├── oled.md                  ← Phase 3
    ├── joystick.md              ← Phase 4
    ├── split.md                 ← Phase 5
    └── storage.md               ← Phase 6
```

---

## ディレクトリの役割

### ソースコード

| ディレクトリ | 役割 | machine 依存 |
|---|---|---|
| `engine/` | 全コンポーネントを統合するメインループ。パッケージ依存の最上位層 | あり |
| `engine/peripheral/` | 周辺機器ドライバ群。`engine` からはインターフェース経由で参照（Phase 3 以降） | あり |
| `matrix/` | マトリクススキャン（GPIO）とデバウンスロジック。ファイル単位で依存を分離 | 一部（`matrix.go` のみ） |
| `keycode/` | HID キーコードの型・定数定義。他パッケージに依存しない | なし |
| `examples/` | ボード別のエントリポイント・設定・キーマップ。ビルドターゲットはここを指定する | あり |

### ドキュメント

| ディレクトリ / ファイル | 役割 |
|---|---|
| `docs/` | プロジェクト全体のドキュメント置き場 |
| `docs/packages/` | パッケージ実装レベルの仕様。コード実装と同時に更新する |
| `.steering/` | 作業単位（feature ブランチ単位）のスペック駆動ドキュメント |
| `.steering/[タイトル]/` | 1 つの feature ブランチに対応する要求・設計・タスクリスト |

### 設定・ツール

| ファイル / ディレクトリ | 役割 |
|---|---|
| `Makefile` | ビルド・テスト・フラッシュ・フォーマット・Lint の統一インターフェース |
| `go.mod` | Go モジュール定義（Go 1.25）。フレームワーク本体のモジュール |
| `go.work` | Go ワークスペース定義。`examples/` を別モジュールとして扱いながらローカル参照を可能にする |
| `.vscode/` | VS Code 設定（TinyGo 拡張設定含む）。リポジトリにコミットして共有する |
| `.idea/` | IntelliJ / GoLand 設定。リポジトリにコミットして共有する |

---

## ファイル配置ルール

### ソースファイル

**`//go:build tinygo` タグの使い分け**

フレームワーク本体（`engine/`, `matrix/`, `keycode/`）のうち、`machine` パッケージに依存するファイルには必ず `//go:build tinygo` タグを付ける。
これにより、標準 `go test` でのビルドエラーを防ぐ。

`examples/` は別モジュールかつ TinyGo 専用ビルドターゲットであるため、タグなしで `machine` を import してよい。

| タグ | 対象ファイル例 |
|---|---|
| `//go:build tinygo` | `matrix/matrix.go`, `engine/keyboard.go` |
| タグなし | `matrix/debounce.go`, `keycode/keycode.go`, `matrix/const.go` |
| タグなし（`examples/` のみ） | `examples/zero-kb02/config.go`, `examples/zero-kb02/main.go` |

**テストファイル**

- `*_test.go` は対象と同じパッケージ内に配置する
- `machine` 非依存パッケージのみ標準 `go test` で実行可能
- `machine` 依存パッケージのテストは `tinygo test -target=waveshare-rp2040-zero`

**エラー定義**

- `errors.New` で定義したエラーは各パッケージの `errors.go` にまとめる
- 例: `engine/errors.go`, `matrix/errors.go`

**定数**

- パッケージ内のマジックナンバーは定数化する
- 複数ファイルで共有する定数は `const.go` に集約する（例: `matrix/const.go`）

### ドキュメントファイル

**`docs/packages/*.md`**

- パッケージごとに 1 ファイル（例: `docs/packages/matrix.md`）
- コード実装と同時に作成・更新する（実装後に書かない）
- Phase 完了時に対応ファイルが揃っていることを確認する

**`.steering/[タイトル]/`**

feature ブランチを作成したら、実装より前に以下の順で作成する。

| ファイル | 書くタイミング |
|---|---|
| `requirements.md` | feature ブランチ作成直後（実装より前） |
| `design.md` | requirements 確定後・実装前 |
| `tasklist.md` | design 確定後・実装中に更新 |

タイトルはブランチ名や Phase 名に対応させる（例: `.steering/phase1/`）。

### ビルド成果物

- `.elf` / `.uf2` などのビルド成果物は `.gitignore` で除外する
- `zero-kb02.elf` などをリポジトリにコミットしない
