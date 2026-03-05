# Phase 6 要求定義: アナログジョイスティック（マウスポインティングデバイス）

## 背景

Phase 5 でフレームワーク API が整理され、利用者は `engine` + `keycode` の 2 パッケージのみで構成可能になった。Phase 6 では ROADMAP に記載のジョイスティック対応を進める。

### 方針転換

TinyGo の Composite HID（keyboard + gamepad）は GitHub Issue #3474 により動作しない。一方、keyboard + mouse の Composite HID は sago35/tinygo-keyboard で実証済み。そのため、ジョイスティックの ADC 値を **マウス移動** にマッピングする方式を採用する。

## 機能要件

### FR-1: アナログジョイスティックによるマウスカーソル移動

- ADC ピン 2 本（X 軸・Y 軸）からアナログ値を読み取る
- ADC 値をマウス移動量（dx, dy: int8）に変換する
- `mouse.Move(dx, dy, 0)` でマウスカーソルを移動させる
- 既存のキーボード HID と共存する（Composite HID: keyboard + mouse）

### FR-2: デッドゾーン

- ジョイスティック静止時の ADC ノイズによるカーソルジッターを防止する
- 中心値（32768）から閾値以内の値はゼロとして扱う
- デフォルト閾値: 3000（65535 の約 ±4.6%）
- ユーザーが閾値をカスタマイズ可能

### FR-3: 感度設定

- マウス移動速度を 1〜10 の整数値で設定可能
- デフォルト: 5
- 1 = 最低速（精密操作向け）、10 = 最高速（素早い移動向け）

### FR-4: 軸反転

- X 軸・Y 軸をそれぞれ独立に反転可能
- ジョイスティックの取り付け方向に応じて調整

### FR-5: ジョイスティックボタン

- ジョイスティックの押し込みボタンを 1 ピンで検出
- ボタン有効/無効を設定可能
- ボタン動作:
  - `ButtonKey == keycode.None` → マウス左クリック（`mouse.Click(mouse.Left)`）
  - `ButtonKey != keycode.None` → 指定キーコードを HID 送信

### FR-6: engine.Config への統合

- `JoystickConfig` 構造体を `engine` パッケージに追加
- `Config.Joystick *JoystickConfig` フィールドで有効化（nil = 無効）
- 既存の Encoder/LED/OLED と同じパターン

## 非機能要件

### NFR-1: ホットパス零アロケーション

- Tick() 内でヒープアロケーションを行わない
- 固定サイズバッファのみ使用

### NFR-2: テスタビリティ

- 軸マッピングロジック（ADC → dx/dy）は `machine` 非依存で実装
- 標準 `go test` でユニットテスト可能

### NFR-3: 既存機能への非影響

- Phase 5 の全機能（キーボード・エンコーダー・LED・OLED）が同等に動作
- ジョイスティック未設定時は一切の副作用なし

### NFR-4: MaxPeripherals の拡張

- 標準ペリフェラル 4 種（encoder/led/oled/joystick）+ カスタム余地のため MaxPeripherals を 4→8 に増加

## ユーザーストーリー

### US-1: ジョイスティックでマウス操作

> キーボードビルダーとして、キーボードに搭載したアナログジョイスティックでマウスカーソルを動かしたい。キーボードから手を離さずにマウス操作ができるようになる。

### US-2: デッドゾーン調整

> キーボードビルダーとして、ジョイスティック静止時にカーソルが勝手に動かないようにしたい。ADC ノイズによるジッターを防止できる。

### US-3: 感度調整

> キーボードビルダーとして、マウスの移動速度を調整したい。用途に応じて精密操作と素早い移動を切り替えられる。

## 受け入れ条件

- [ ] `make test`: axis_test.go のユニットテストが全件パス
- [ ] `make test`: TinyGo ビルドが成功
- [ ] `make lint`: 静的解析パス
- [ ] ジョイスティック静止時にカーソルが動かない（デッドゾーン機能）
- [ ] ジョイスティック傾倒時にカーソルが対応方向に移動
- [ ] ボタン押下でマウスクリック（またはキーコード送信）
- [ ] 既存機能（キーボード・エンコーダー・LED・OLED）が Phase 5 と同等に動作
- [ ] `_ "machine/usb/hid/mouse"` の blank import で mouse HID が有効化される

## 制約事項

- TinyGo 標準の `machine/usb/hid/mouse` パッケージを使用する
- Composite HID は keyboard + mouse のみ（gamepad は TinyGo Issue #3474 により不可）
- ADC 読み取りは `machine.ADC` を使用する
- マウス移動量は int8 範囲（-128〜127）
