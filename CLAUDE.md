# CLAUDE.md

## プロジェクト概要

**keygoard** — TinyGoベースのカスタムキーボードファームウェアフレームワーク

- ターゲットMCU: RP2040（Waveshare RP2040-Zero）
- TinyGoターゲット名: `waveshare-rp2040-zero`
- 開発方針: **最小構成から実機動作確認しながら段階的に積み上げる**
  - v0.1.0 で「USB HIDとして認識されない」問題が発生したため、ゼロベースで再設計
  - 各フェーズで実機動作確認を必須とし、確認済みのものだけ次フェーズへ進む

---

## テスト実機：zero-kb02（sago35/keyboards）

- リポジトリ: https://github.com/sago35/keyboards/tree/main/zero-kb02
- 別ファームウェアで動作確認済みのため、開発中のエラーはハード・マイコン故障ではなくファームウェアの問題とみなしてよい

### ハードウェア仕様（ピン配置）

| 機能 | GPIO |
|------|------|
| Matrix Col 0–3（出力） | GP5, GP6, GP7, GP8 |
| Matrix Row 0–2（入力） | GP9, GP10, GP11 |
| Matrix 方式 | COL2ROW（列をHighにして行をRead） |
| RGB LED（SK6812MINI-E） | GP1、12個 |
| Encoder A / B / Button | GP3 / GP4 / GP0 |
| OLED SDA / SCL | GP12 / GP13（I2C、SSD1306 128×64、addr=0x3C） |
| Joystick X / Y / Button | GP29(ADC) / GP28(ADC) / GP0 |

### テスト環境

- 主: macOS（USB Prober / System Information で認識確認）
- 副: 将来的に複数 OS で確認

---

## 開発フェーズ

詳細は [ROADMAP.md](ROADMAP.md) を参照。**現在地: Phase 1 着手前。**

---

## 開発戦略

### ブランチ戦略

フェーズ内のトピック単位でブランチを細かく分けることで、巻き戻しを容易にする。

```
main
└── phase1
    ├── feature/matrix-scan       ← マトリクススキャン実装
    ├── feature/debounce          ← デバウンス実装
    ├── feature/hid-basic         ← 基本HID実装
    └── feature/keymap-layer0     ← レイヤー0キーマップ実装
```

- `main` — 安定版（各 Phase 完了・実機確認済みのものだけマージ）
- `phase{N}` — Phase 開発ブランチ（`main` から分岐）
- `feature/xxx` — トピックブランチ（`phase{N}` から分岐・マージ）

フロー: `feature/xxx` → `phase{N}` → `main`

### バージョニング

`v0.{マイナー}.{パッチ}` の形式で管理する。

| タグ | 意味 |
|------|------|
| `v0.(N-1).1`, `v0.(N-1).2`... | Phase N 内の開発チェックポイント（`phase{N}` ブランチ上） |
| `v0.N.0` | Phase N 完了（実機確認済み・`main` マージ時） |

フェーズ完了時に初めてマイナーバージョンが `N` に上がる。開発中のチェックポイントは常に前フェーズのマイナーバージョンのパッチとして積み上がる。バージョンは常に単調増加する。

**例:**

```
v0.0.1  ← Phase 1: マトリクススキャン完成チェックポイント
v0.0.2  ← Phase 1: HID 認識チェックポイント
v0.1.0  ← Phase 1 完了（実機確認・main マージ）
v0.1.1  ← Phase 2: MO レイヤーチェックポイント
v0.1.2  ← Phase 2: TT/LT チェックポイント
v0.2.0  ← Phase 2 完了（実機確認・main マージ）
v0.2.1  ← Phase 3: エンコーダーチェックポイント
v0.3.0  ← Phase 3 完了（実機確認・main マージ）
```

### 仕様書

- `docs/packages/` フォルダにコードに対する仕様書を日本語で記載する
- パッケージごとに 1 ファイル（例: `docs/packages/matrix.md`, `docs/packages/hid.md`）
- 仕様書はコード実装と同時に更新する（実装後に書かない）
- Phase 完了時に対応するドキュメントが揃っていることを確認する

```
docs/packages/
├── matrix.md        ← マトリクススキャン・デバウンス仕様
├── hid.md           ← HID レポート仕様
├── keycode.md       ← キーコード定義仕様
├── layer.md         ← レイヤーシステム仕様（Phase 2 以降）
├── encoder.md       ← エンコーダー仕様（Phase 3 以降）
├── led.md           ← LED 仕様（Phase 3 以降）
├── oled.md          ← OLED 仕様（Phase 3 以降）
├── joystick.md      ← ジョイスティック仕様（Phase 4 以降）
├── split.md         ← Split 通信仕様（Phase 5 以降）
└── storage.md       ← Flash 永続化仕様（Phase 6 以降）
```

---

## TinyGo 制約（重要：必ず守ること）

### 使用禁止
- `reflect` パッケージ（TinyGo でランタイムエラーになる）
- `fmt.Sprintf` / `fmt.Printf`（ヒープ割り当てが大きい）
  - デバッグ出力は `println()` のみ使用する
- `net` / `os` など TinyGo 未対応の標準ライブラリ
- goroutine の多用（スケジューラオーバーヘッド。メインループは 1 つの goroutine で動かす）

### ヒープ割り当てを避ける（ホットパス内）
- スキャンループ内での `append` による再割り当ては禁止
- 固定サイズ配列を優先する（スライスの初期確保は初期化時に行う）
- 文字列連結より固定バイト列を使う

### USB HID 実装上の注意（v0.1.0 の失敗から）
- **TinyGo 標準の `machine/usb/hid` を使う**（自前 HID 記述子は書かない）
- **Composite HID は Phase 4 まで追加しない**（キーボード単体の認識を先に確実にする）
- USB エニュメレーション完了まで（起動後数百 ms）はキースキャンを開始しない

### サードパーティパッケージ
- 新規追加は **事前にユーザーへ確認すること**
- `tinygo-org/drivers` は使用許可済み

### メモリ目標
- RAM 使用量 128KB 以下

---

## アーキテクチャ方針

### パッケージ依存の方向（下方向のみ）

```
engine/          ← 全体を統合する最上位層
  ├── matrix/    ← マトリクススキャン・デバウンス
  ├── hid/       ← HID レポート生成
  ├── keycode/   ← キーコード定義（他に依存しない）
  └── peripheral/
        ├── encoder/
        ├── led/
        ├── oled/
        └── joystick/
```

- 循環依存禁止
- `engine` から `peripheral` への依存はインターフェース経由（Phase 3 以降）
- 新しい周辺機器は `peripheral/` 配下に新しいパッケージとして作成する
- Phase 1 は構造を作りすぎない。動くことを最優先にする

---

## ビルド・検証コマンド

開発コマンドは `Makefile` にまとめている。

| コマンド | 内容 |
|---|---|
| `make test` | ① 標準 Go: `machine` 非依存パッケージ（keycode, matrix/debounce）<br>② TinyGo: RP2040 ターゲットで全パッケージのビルド検証（実機接続が必要） |
| `make build` | TinyGo でのビルド確認（zero-kb02） |
| `make flash` | RP2040 への書き込み |
| `make fmt` | コードフォーマット（goimports） |
| `make lint` | Lint（golangci-lint） |

---

## コーディング規約

- フォーマット: `goimports`（TinyGo 固有パッケージの import は手動確認）
- コメント: 日本語で記述する
- マジックナンバーは定数化する
- エラーは `errors.New` で定義し、各パッケージの `errors.go` にまとめる
- グローバル変数: 組み込み上やむを得ない場合のみ許容（必ずコメントで理由を書く）

---

## コードレビュー観点

1. **TinyGo 制約違反がないか** — `reflect` / `fmt.Sprintf` / 未対応パッケージの使用
2. **ホットパス（スキャンループ内）でヒープ割り当てが発生していないか**
3. **現フェーズのスコープ外の実装が混入していないか**
4. **ハードウェア非依存ロジックにユニットテストがあるか** — keycode / debounce / layer
5. **エラーを握りつぶしていないか**
6. **グローバル変数の使用に理由があり、コメントがあるか**
7. **対応する仕様書（docs/packages/）が更新されているか**

---

## テスト方針

### テストコマンドの使い分け

`machine` パッケージをインポートするパッケージは標準 `go test` でビルドエラーになる。
パッケージの依存状況に応じてコマンドを使い分ける。

| コマンド | 対象 | 備考 |
|---|---|---|
| `go test ./keycode/... ./matrix/...` | `machine` 非依存パッケージ | 標準 Go で実行可能。`matrix.go` は `//go:build tinygo` タグで除外される |
| `tinygo test -target=waveshare-rp2040-zero ./...` | 全パッケージ（実機テスト） | 実機接続が必要 |

### ユニットテスト対象

ハードウェア非依存のロジックにはユニットテストを書く。

- Phase 1: `keycode`, `matrix/debounce`（デバウンスロジックは `machine` から切り離して設計する）
- Phase 2 以降: `engine/layer`, `engine/tap`

テストファイルは対象と同じパッケージ内に `*_test.go` で配置する。

### 実機テスト

- **各フェーズ完了時に必須。実機確認なしに次フェーズへ進まない**
- macOS で USB キーボードとして認識されることを確認してから `v0.N.0` タグを打つ

---

## コミット規約

Conventional Commits 形式、日本語で記述する。

```
<type>: <説明>
```

| type | 用途 |
|------|------|
| `feat:` | 新機能 |
| `fix:` | バグ修正 |
| `docs:` | ドキュメント |
| `refactor:` | リファクタリング |
| `test:` | テスト追加・修正 |
| `chore:` | ビルド・設定等の雑務 |