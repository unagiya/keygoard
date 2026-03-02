# Phase 5: パッケージ構造リファクタリング — 要求定義

## 目的

フレームワーク利用者が `engine` と `keycode` の 2 パッケージだけで
キーボードファームウェアを構築できるようにする。
内部実装を `internal/` に移動し、公開 API の表面積を最小化する。

## 背景

Phase 4 完了時点で、利用者の `main.go` は 6 パッケージを import する必要がある：

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
    "github.com/unagiya/keygoard/matrix"                    // ← 内部詳細
    "github.com/unagiya/keygoard/engine/peripheral/encoder"  // ← 内部詳細
    "github.com/unagiya/keygoard/engine/peripheral/led"      // ← 内部詳細
    "github.com/unagiya/keygoard/engine/peripheral/oled"     // ← 内部詳細
)
```

また `machine` パッケージを直接 import してピンや I2C バスを指定する必要があり、
フレームワークの抽象化が漏れている。

## 変更後のユーザー体験

```go
import (
    "github.com/unagiya/keygoard/engine"
    "github.com/unagiya/keygoard/keycode"
)

func main() {
    kb := engine.New(&engine.Config{
        ProductName: "zero-kb02",
        ColPins:     []engine.Pin{5, 6, 7, 8},
        RowPins:     []engine.Pin{9, 10, 11},
        Keymap:      defaultKeymap,
        Encoder: &engine.EncoderConfig{
            PinA: 3, PinB: 4,
            KeyCW: keycode.UpArrow, KeyCCW: keycode.DownArrow,
        },
        LED: &engine.LEDConfig{
            Pin:   1,
            Count: 12,
            LayerColors: [engine.MaxLayerColors]color.RGBA{ ... },
        },
        OLED: &engine.OLEDConfig{
            Bus:      engine.I2C0,
            SDA:      12,
            SCL:      13,
            Address:  0x3C,
            Width:    128,
            Height:   64,
            Rotation: engine.Rotation180,
        },
    })
    kb.Run()
}
```

## 要求事項

### R-1: 公開パッケージの限定

フレームワーク利用者が import するパッケージを以下の 2 つに限定する：

- `github.com/unagiya/keygoard/engine` — Facade（Config / New / Run）
- `github.com/unagiya/keygoard/keycode` — キーコード定数・レイヤーアクション関数

### R-2: 内部パッケージの隠蔽

以下のパッケージを `internal/` 配下に移動し、外部からの import を禁止する：

- `layer`（レイヤー解決）→ `internal/layer/`
- `tap`（タップ/ホールド判定）→ `internal/tap/`
- `matrix`（スキャン + デバウンス）→ `internal/matrix/`
- `encoder`（ロータリーエンコーダー）→ `internal/peripheral/encoder/`
- `led`（RGB LED）→ `internal/peripheral/led/`
- `oled`（OLED ディスプレイ）→ `internal/peripheral/oled/`

### R-3: engine を Facade 化

`engine.Config` にペリフェラルおよびマトリクスの設定を統合する：

- マトリクスのピン設定を `Config` に直接持たせる（利用者が `matrix.New()` を呼ばない）
- 各ペリフェラルの設定を Config のオプショナルフィールドとして持つ
- カスタムペリフェラル用に `Peripherals []Peripheral` フィールドを維持する

### R-4: ハードウェア抽象型の導入

`machine` パッケージへの依存を利用者から隠蔽する：

- `engine.Pin` 型（`uint8` ベース）を定義し、GPIO 番号を整数で指定可能にする
- `engine.I2CBus` 型を定義し、I2C バスを `engine.I2C0` / `engine.I2C1` で指定可能にする
- `engine.Rotation` 型を定義し、OLED 回転を `engine.Rotation0` 〜 `engine.Rotation270` で指定可能にする
- engine 内部で `machine.Pin()` 等に変換する

### R-5: Peripheral インターフェースの公開維持

`engine.Peripheral` インターフェースは引き続き公開し、
利用者が自作ペリフェラルを作成・登録できるようにする。

### R-6: 既存テストの維持

全ての既存ユニットテストが移動後も `make test` で通ること。
（`make test` は tinygo build + go test を実行する）

### R-7: 実機動作の維持

`make flash` で書き込み後、Phase 4 完了時点と同一の動作をすること。

## 受け入れ条件

- [ ] 利用者の `main.go` が `engine` と `keycode` の 2 パッケージのみで記述できる
- [ ] 利用者の `main.go` が `machine` パッケージを直接 import しない
- [ ] `internal/` 配下のパッケージが外部から import できない
- [ ] `make test` が通る
- [ ] `examples/zero-kb02/` が新しい API で書き直されている
- [ ] 実機（zero-kb02）でキー入力・レイヤー切替・エンコーダー・LED・OLED が Phase 4 と同等に動作する

## 制約事項

- 機能追加は行わない（純粋なリファクタリング）
- キーボードとしての動作が Phase 4 完了時点と同一であること
- RAM 使用量が増加しないこと（ホットパスでのヒープ割り当てなし）
