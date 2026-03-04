# Phase 6 タスクリスト: アナログジョイスティック（マウスポインティングデバイス）

## 実装タスク

### 1. コア実装

- [x] **T-1**: `internal/peripheral/joystick/axis.go` — AxisMapper 構造体と MapDelta() 実装
- [x] **T-2**: `internal/peripheral/joystick/axis_test.go` — ユニットテスト
  - デッドゾーン内 → 0
  - 中心値 → 0
  - 最大傾倒 → maxOutput
  - 軸反転
  - 感度変更による出力変化
- [x] **T-3**: `internal/peripheral/joystick/joystick.go` — Joystick 構造体（ADC + mouse.Move）

### 2. engine 統合

- [x] **T-4**: `engine/peripheral.go` — MaxPeripherals を 4→8 に変更
- [x] **T-5**: `engine/config.go` — JoystickConfig 構造体追加、Config.Joystick フィールド追加
- [x] **T-6**: `engine/keyboard.go` — joystick import 追加、New() にジョイスティック生成ロジック追加

### 3. examples 更新

- [x] **T-7**: `examples/zero-kb02/main.go` — `_ "machine/usb/hid/mouse"` import 追加、Joystick 設定追加

### 4. 検証

- [x] **T-8**: `make test` — ユニットテストパス
- [x] **T-9**: `make lint` — 静的解析パス

### 5. ドキュメント更新

- [x] **T-10**: `docs/packages/joystick.md` — パッケージ仕様書作成
- [x] **T-11**: `docs/packages/engine.md` — JoystickConfig ドキュメント追加
- [x] **T-12**: `docs/functional-design.md` — システム構成図・API 更新
- [x] **T-13**: `docs/repository-structure.md` — joystick ディレクトリ追加
- [x] **T-14**: `docs/architecture.md` — mouse HID をスタックに追加
- [x] **T-15**: `docs/glossary.md` — joystick/ADC 用語追加
- [x] **T-16**: `ROADMAP.md` — Phase 6 説明更新（gamepad → mouse）
- [x] **T-17**: `README.md` — ジョイスティック設定セクション追加

## 実装順序

```
T-1 → T-2 → T-3 → T-4 → T-5 → T-6 → T-7 → T-8 → T-9 → T-10〜T-17
```

## 完了条件

- 全ユニットテストがパス（`make test`）
- 静的解析がパス（`make lint`）
- TinyGo ビルドが成功（`make build`）
- 全ドキュメントが最新の状態
