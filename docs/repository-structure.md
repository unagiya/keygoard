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
├── go.sum                       ← Go 依存ハッシュ
├── go.work                      ← Go ワークスペース定義（ローカル開発用）
│
├── engine/                      ← 公開 Facade（利用者が import するパッケージ）
│   ├── keyboard.go              ← メインループ・HID 送信（//go:build tinygo）
│   ├── config.go                ← Config・ペリフェラル設定（//go:build tinygo）
│   ├── types.go                 ← Pin / I2CBus / Rotation 型定義
│   ├── peripheral.go            ← Peripheral インターフェース定義
│   ├── keymap.go                ← Keymap 構造体・RowCount/ColCount 再エクスポート
│   └── errors.go                ← エラー定義
│
├── keycode/                     ← 公開パッケージ（利用者が import するパッケージ）
│   ├── keycode.go               ← キーコード型・定数
│   └── keycode_test.go          ← ユニットテスト
│
├── internal/                    ← 非公開パッケージ（外部 import 不可）
│   ├── matrix/                  ← マトリクススキャン・デバウンス
│   │   ├── const.go             ← マトリクスサイズ定数
│   │   ├── matrix.go            ← GPIO スキャン実装（//go:build tinygo）
│   │   ├── debounce.go          ← デバウンスロジック（machine 非依存）
│   │   └── debounce_test.go     ← デバウンスユニットテスト
│   ├── layer/                   ← レイヤー解決
│   │   ├── resolver.go          ← Resolver（最上位アクティブレイヤーから走査）
│   │   └── resolver_test.go     ← ユニットテスト
│   ├── tap/                     ← タップ/ホールド判定
│   │   ├── detector.go          ← Detector（LT/TT キー用）
│   │   └── detector_test.go     ← ユニットテスト
│   └── peripheral/              ← 周辺機器ドライバ群
│       ├── encoder/             ← ロータリーエンコーダー
│       │   ├── quadrature.go    ← クワドラチャデコーダ（machine 非依存）
│       │   ├── quadrature_test.go ← ユニットテスト
│       │   └── encoder.go       ← GPIO 読み取り（//go:build tinygo）
│       ├── led/                 ← RGB LED（WS2812）
│       │   └── led.go           ← LED 制御（//go:build tinygo）
│       ├── oled/                ← OLED ディスプレイ（SSD1306）
│       │   ├── font.go          ← 8x8 ビットマップフォント定義
│       │   ├── font_test.go     ← ユニットテスト
│       │   └── oled.go          ← SSD1306 I2C 制御（//go:build tinygo）
│       └── joystick/            ← アナログジョイスティック
│           ├── axis.go          ← 軸マッピングロジック（machine 非依存）
│           ├── axis_test.go     ← ユニットテスト
│           └── joystick.go      ← ADC + mouse.Move（//go:build tinygo）
│
├── examples/                    ← フレームワーク使用例（別モジュール）
│   └── zero-kb02/               ← zero-kb02 向けファームウェアサンプル
│       ├── go.mod               ← モジュール定義（github.com/.../examples/zero-kb02）
│       ├── main.go              ← エントリポイント（engine + keycode のみ import）
│       └── keymap.go            ← キーマップ定義
│
├── docs/                        ← プロジェクトドキュメント
│   ├── product-requirements.md  ← プロダクト要求定義
│   ├── functional-design.md     ← 機能設計（システム構成図・API）
│   ├── architecture.md          ← 技術仕様（スタック・制約・パフォーマンス）
│   ├── repository-structure.md  ← 本ファイル
│   ├── glossary.md              ← ユビキタス言語定義
│   ├── development-guidelines.md ← 開発ガイドライン
│   └── packages/                ← パッケージ実装レベル仕様
│       ├── engine.md
│       ├── keycode.md
│       ├── matrix.md
│       ├── layer.md
│       ├── tap.md
│       ├── encoder.md
│       ├── led.md
│       ├── oled.md
│       └── joystick.md
│
└── .steering/                   ← 作業単位ドキュメント（スペック駆動開発）
    ├── phase1/ 〜 phase5/       ← 各 Phase の作業ドキュメント
    │   ├── requirements.md
    │   ├── design.md
    │   └── tasklist.md
    ...
```

### 将来の構成（Phase 7 以降）

```
keygoard/
├── ...（上記と同様）
│
└── docs/packages/               ← フェーズ進行に伴い追加
    ├── split.md                 ← Phase 7
    └── storage.md               ← Phase 8 以降
```

---

## ディレクトリの役割

### ソースコード

| ディレクトリ | 役割 | machine 依存 |
|---|---|---|
| `engine/` | 公開 Facade。Config・Peripheral インターフェース・ハードウェア抽象型を提供 | あり |
| `keycode/` | HID キーコードの型・定数定義。他パッケージに依存しない | なし |
| `internal/matrix/` | マトリクススキャン（GPIO）とデバウンスロジック | 一部（`matrix.go` のみ） |
| `internal/layer/` | レイヤー解決ロジック | なし |
| `internal/tap/` | タップ/ホールド判定 | なし |
| `internal/peripheral/` | 周辺機器ドライバ群 | あり |
| `examples/` | ボード別のエントリポイント・キーマップ。ビルドターゲットはここを指定する | あり |

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

フレームワーク本体（`engine/`, `internal/`, `keycode/`）のうち、`machine` パッケージに依存するファイルには必ず `//go:build tinygo` タグを付ける。
これにより、標準 `go test` でのビルドエラーを防ぐ。

`examples/` は別モジュールかつ TinyGo 専用ビルドターゲットであるため、タグなしで `machine` を import してよい。

| タグ | 対象ファイル例 |
|---|---|
| `//go:build tinygo` | `engine/keyboard.go`, `engine/config.go`, `internal/matrix/matrix.go`, `internal/peripheral/encoder/encoder.go`, `internal/peripheral/led/led.go`, `internal/peripheral/oled/oled.go`, `internal/peripheral/joystick/joystick.go` |
| タグなし | `engine/types.go`, `engine/keymap.go`, `engine/peripheral.go`, `keycode/keycode.go`, `internal/matrix/const.go`, `internal/matrix/debounce.go`, `internal/layer/resolver.go`, `internal/tap/detector.go`, `internal/peripheral/encoder/quadrature.go`, `internal/peripheral/oled/font.go`, `internal/peripheral/joystick/axis.go` |
| タグなし（`examples/` のみ） | `examples/zero-kb02/main.go` |

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
