# keygoard ロードマップ

このドキュメントでは、keygoardプロジェクトの開発計画と進捗状況を示します。

[English](#english) | [日本語](#日本語)

## 日本語

### 概要

keygoardは4つのPhaseに分けて開発を進めています。各Phaseは特定の機能群に焦点を当てており、段階的に機能を追加していきます。

---

## Phase 1: 基本機能 ✅ 完了

**目標**: ゲーミングキーボードとして完全に動作する基本機能の実装

### 実装済み機能

#### コア機能
- ✅ マトリックススキャン
  - COL2ROW / ROW2COL 両対応
  - 3ms デバウンス（設定可能）
  - 1ms スキャン間隔（1000Hz）
  - Eager デバウンスアルゴリズム
- ✅ キーコード体系
  - 基本キー（USB HID準拠）
  - モディファイア（8種類）
  - レイヤー操作（MO, TG, TT, LT）
  - ゲームパッドボタン（32個）
- ✅ レイヤーシステム
  - 最大16層
  - ビットマップベース高速切り替え
  - 透過キー対応
- ✅ USB HID
  - 6KRO キーボード
  - 32ボタンゲームパッド
  - 1000Hz レポートレート

#### 周辺機器
- ✅ アナログジョイスティック
  - 12bit ADC入力
  - オーバーサンプリング（4サンプル平均）
  - 起動時自動キャリブレーション
  - デッドゾーン処理（±5%）
  - 軸反転機能
- ✅ OLEDディスプレイ
  - SSD1306ドライバ（I2C）
  - 128x32 / 128x64 対応
  - レイヤー番号、ジョイスティック値、ロック状態表示

#### 分割キーボード
- ✅ UART通信プロトコル
  - 460800bps
  - マトリックス状態転送
  - ジョイスティック値転送
  - チェックサム検証
- ✅ マスター/スレーブ構成
  - 片側最大16x8マトリックス
  - 両側ジョイスティック対応（2軸分）
  - 5ms タイムアウト
  - 自動フォールバック

#### 開発ツール
- ✅ CLIツール（`keygoard`）
  - プロジェクト生成（`keygoard init`）
  - バージョン表示
- ✅ サンプル実装
  - simple-gamepad（ジョイスティック付き4x6キーボード）

#### ドキュメント
- ✅ 日本語・英語README
- ✅ クイックスタートガイド
- ✅ キーコードリファレンス
- ✅ サンプル実装ガイド
- ✅ 開発環境セットアップ

### 既知の制限（Phase 1）

- ⚠️ TT/LTのタップ検出が未実装（長押しのみ動作） → ✅ Phase 2 M4で解決
- ⚠️ ゲームパッドHID記述子が不完全（キーボードは完全動作）
- ⚠️ OLED表示がマスター側のみ（スレーブ側は未対応） → ✅ Phase 2 M3で解決
- ⚠️ LEDサポートなし → ✅ Phase 2 M2で解決

---

## Phase 2: 拡張機能 ✅ 完了

**目標**: 周辺機器のサポート拡充と機能強化

**期間**: 2024年Q1-Q2

### 実装済み機能

#### ロータリーエンコーダー
- ✅ 基本エンコーダーサポート
  - 最大2個対応
  - 時計回り/反時計回りイベント
  - クリック（プッシュスイッチ）対応
  - キーコードへのマッピング
  - Gray codeルックアップテーブル
  - エンジンへの統合

#### RGB LED（WS2812）
- ✅ 基本LED制御
  - 最大128個対応
  - GRB順序対応
  - WS2812ドライバ統合
- ✅ 基本エフェクト
  - Static（単色）
  - Breathing（呼吸、サイン波）
  - Rainbow（レインボー、HSV変換）
- ✅ エンジン統合
  - 自動更新ループ
  - エフェクト切り替えAPI

#### 分割キーボード双方向通信
- ✅ マスター→スレーブ通信
  - LED同期（RGB byte配列）
  - OLED同期（テキスト/バッファ）
  - 460800bps UART
- ✅ OLED両側表示
  - スレーブ側OLED対応
  - テキスト表示同期
  - コールバックベース実装

#### ジョイスティック機能拡張
- ✅ 詳細キャリブレーション
  - 手動再キャリブレーション
  - 中心位置設定API
  - 生ADC値取得
- ✅ デッドゾーン設定
  - 可変デッドゾーン
  - デフォルト5%
  - カスタマイズ可能

#### タップ検出
- ✅ TT/LT のタップ判定実装
  - タップ時間閾値（200ms）
  - 長押し判定（150ms）
  - ステートマシン実装
- ✅ 設定可能なタップ時間
  - TapConfig構造体
  - TappingTerm / HoldTerm

### Phase 2 マイルストーン

1. **M1: ロータリーエンコーダー基本実装** ✅
   - エンコーダー読み取り
   - キーコードマッピング
   - サンプル実装 (encoder-test)

2. **M2: RGB LED基本実装** ✅
   - LED制御
   - 基本エフェクト (Static, Breathing, Rainbow)
   - サンプル実装 (led-test)

3. **M3: 双方向通信** ✅
   - UART双方向通信
   - LED/OLED同期
   - サンプル実装 (split-bidirectional)

4. **M4: タップ検出とジョイスティック拡張** ✅
   - タップ判定実装
   - キャリブレーション機能
   - サンプル実装 (tap-detection)

**完了**: Phase 2の全機能実装完了

---

## Phase 3: 高度な機能 ✅ 完了

**目標**: LEDエフェクトの高度化とパフォーマンス最適化

### 実装済み機能

#### 高度なLEDエフェクト
- ✅ リアクティブエフェクト
  - KeyPressEffect（キー押下中に点灯、離すと消灯）
  - FadeOutEffect（キー押下で点灯、時間経過で減衰）
  - RippleEffect（キー押下位置から波紋が広がる）
  - ReactiveEffectインターフェース
  - KeyToLEDマッピング（マトリクス座標→LEDインデックス）
- ✅ エフェクトレジストリ
  - 最大16エフェクト登録
  - Next/Previous切り替え
  - IDベースのエフェクト選択

#### 分割キーボード両側LED制御
- ✅ スレーブキーイベント転送
  - MsgTypeKeyEvent（Slave→Master）
  - リアクティブエフェクトへの通知
- ✅ Unifiedモード（デフォルト）
  - マスターが両側のLEDバッファを統合管理
  - スレーブのキーイベントをマスターに転送
  - 両側でリアクティブエフェクトが動作
- ✅ Independentモード
  - 各側が独立してエフェクト実行
  - MsgTypeLEDMode（Master→Slave）でモード通知

#### パフォーマンス最適化
- ✅ メモリ使用量削減（ヒープ割り当て95%削減）
  - マトリクススキャンバッファの事前割り当て
  - デバウンサ結果バッファの事前割り当て
  - LEDカラー変換バッファの事前割り当て
  - splitプロトコル送信バッファの事前割り当て（txBuffer [68]byte）
  - EncodeInto() / DecodeMatrixStateInto() ゼロアロケーションAPI
  - TapDetector固定配列バッファ
  - OLED同期の変更検出と固定バイト列化

#### ユーザビリティ向上
- ✅ デバッグモード
  - println()ベースの軽量デバッグロガー（DebugLogger）
  - BoardConfig.Debugフラグによる有効/無効制御
  - キーイベント、レイヤー変更、Split接続状態のログ出力
  - 起動時ログ（マトリクスサイズ、周辺機器、Split設定）
- ✅ エラー処理改善
  - Stats構造体によるエラーカウンタ（Split Rx/Tx、HID、LED）
  - runMaster()/runSlave()の全エラー箇所にハンドリング追加
  - GetStats() APIで統計情報を取得可能

### Phase 3 マイルストーン

1. **M1: リアクティブLEDエフェクト** ✅
   - ReactiveEffectインターフェース
   - KeyPress / FadeOut / Ripple エフェクト
   - エフェクトレジストリ
   - サンプル実装 (reactive-led)

2. **M2: 分割キーボード両側LED制御** ✅
   - スレーブキーイベント転送プロトコル
   - Unified/Independentモード
   - 統合LEDバッファ
   - サンプル実装 (split-led)

3. **M3: パフォーマンス最適化** ✅
   - ホットパスのヒープ割り当て95%削減
   - 事前割り当てバッファによるゼロアロケーション化
   - OLED同期の変更検出最適化

4. **M4: ユーザビリティ向上** ✅
   - デバッグモード（DebugLogger + 起動ログ）
   - エラー処理改善（Stats + 全エラー箇所のハンドリング）

**完了**: Phase 3の全機能実装完了

---

## Phase 4: VIA/Remap対応 ✅ 完了

**目標**: 動的設定とエコシステム統合

### 実装済み機能

#### NKRO対応
- ✅ NKROキーボードHID記述子
  - 全キー同時押し対応（ビットマップベース）
  - 6KROとの動的切り替え（KC_NKRO_TOGGLE / KC_NKRO_ON / KC_NKRO_OFF）
  - CompositeHIDによる統合管理

#### 設定の永続化
- ✅ Flash保存（storageパッケージ）
  - キーマップ保存（最大8KB）
  - 設定保存（HIDモード、LEDエフェクト、輝度）
  - マクロデータ保存（最大1.5KB）
  - XORローテーションチェックサムによるデータ整合性検証
  - SaveSettings / LoadSettings / ResetSettings API

#### マクロ機能
- ✅ キーシーケンスマクロ（MacroManager）
  - 最大16マクロ、各最大32ステップ
  - KeyDown / KeyUp / KeyTap / Delay ステップタイプ
  - BoardConfig.Macrosによる宣言的定義
  - RegisterMacro() APIによる動的登録

#### コンボキー
- ✅ 複数キー同時押し検出（ComboDetector）
  - 最大16コンボ、各2〜4キー
  - 設定可能な検出ウィンドウ（デフォルト50ms）
  - BoardConfig.Combosによる宣言的定義
  - RegisterCombo() APIによる動的登録

### Phase 4 マイルストーン

1. **M1: NKRO対応** ✅
   - NKROキーボードHID記述子
   - 6KRO/NKRO動的切り替え
   - サンプル実装 (nkro-test)

2. **M2: マクロ機能** ✅
   - MacroManager実装
   - キーシーケンスマクロ + 遅延サポート
   - サンプル実装 (macro-test)

3. **M3: コンボキー** ✅
   - ComboDetector実装
   - 複数キー同時押し検出
   - サンプル実装 (combo-test)

4. **M4: 設定の永続化** ✅
   - storageパッケージ実装
   - Flash保存/読み込み/リセット
   - サンプル実装 (persistent-keymap)

**完了**: Phase 4の全機能実装完了

---

## 将来的な可能性

これらは具体的な計画ではありませんが、将来検討する可能性がある機能です：

- 🤔 Bluetooth LE対応（nRF52シリーズ）
- 🤔 ディスプレイ強化（カラーOLED、E-Ink等）
- 🤔 追加センサー（加速度、ジャイロ等）
- 🤔 オーディオフィードバック（ブザー、スピーカー）
- 🤔 ハプティックフィードバック
- 🤔 他のマイコン対応（STM32、ESP32等）
- 🤔 Web設定ツール（ブラウザベース）

---

## 貢献方法

各Phaseの機能実装に貢献したい場合：

1. [CONTRIBUTING.md](CONTRIBUTING.md) を確認
2. 関連するIssueを確認（またはIssueを作成）
3. 実装したい機能について議論
4. Pull Requestを作成

---

## バージョニング

セマンティックバージョニングを採用：

- **v0.1.0**: ベータ版（Phase 1〜4完了）✅

**現在のバージョン**: v0.1.0（ベータ版開発完了）

---

## English

### Overview

keygoard development is divided into 4 phases. Each phase focuses on specific feature sets and adds functionality incrementally.

---

## Phase 1: Core Features ✅ Completed

**Goal**: Implement basic features for a fully functional gaming keyboard

### Implemented Features

#### Core
- ✅ Matrix scanning (COL2ROW/ROW2COL, 3ms debounce, 1000Hz)
- ✅ Keycode system (basic keys, modifiers, layers, gamepad buttons)
- ✅ Layer system (16 layers, bitmap-based)
- ✅ USB HID (6KRO keyboard + 32-button gamepad, 1000Hz)

#### Peripherals
- ✅ Analog joystick (12-bit ADC, oversampling, calibration)
- ✅ OLED display (SSD1306, 128x32/64)

#### Split Keyboard
- ✅ UART communication (460800bps, master/slave)

#### Tools
- ✅ CLI tool (`keygoard init`)
- ✅ Sample implementation (simple-gamepad)

#### Documentation
- ✅ Japanese/English README
- ✅ Quick start guide
- ✅ Keycode reference

### Known Limitations (Phase 1)

- ⚠️ TT/LT tap detection not implemented → ✅ Resolved in Phase 2 M4
- ⚠️ Gamepad HID descriptor incomplete
- ⚠️ OLED on master side only → ✅ Resolved in Phase 2 M3
- ⚠️ No LED support → ✅ Resolved in Phase 2 M2

---

## Phase 2: Extended Features ✅ Completed

**Goal**: Expand peripheral support and functionality

**Period**: 2024 Q1-Q2

### Implemented Features

- ✅ Rotary encoders (2 encoders, click support, Gray code)
- ✅ RGB LED (WS2812, up to 128 LEDs, 3 basic effects)
- ✅ Bidirectional split communication (LED/OLED sync)
- ✅ OLED on both sides (callback-based sync)
- ✅ Joystick calibration mode (manual recalibration, deadzone)
- ✅ Tap detection for TT/LT (state machine, configurable timing)

### Examples Added

- encoder-test: Rotary encoder demonstration
- led-test: RGB LED effects showcase
- split-bidirectional: Full bidirectional split keyboard
- tap-detection: TT/LT tap detection and joystick calibration

**Completed**: All Phase 2 features implemented

---

## Phase 3: Advanced Features ✅ Completed

**Goal**: Advanced LED effects and performance optimization

### Implemented Features

- ✅ Reactive LED effects (KeyPress, FadeOut, Ripple)
- ✅ ReactiveEffect interface and KeyToLED mapping
- ✅ Effect registry (up to 16 effects, Next/Previous cycling)
- ✅ Split keyboard both-side LED control
  - Slave key event forwarding (MsgTypeKeyEvent)
  - Unified mode (master controls both sides)
  - Independent mode (each side runs own effects)
  - MsgTypeLEDMode protocol message
- ✅ Performance optimization (95% heap allocation reduction)
  - Pre-allocated buffers for matrix scan, debounce, LED, split protocol
  - Zero-allocation APIs: EncodeInto(), DecodeMatrixStateInto()
  - OLED sync change detection with fixed byte arrays
- ✅ Debug mode (DebugLogger with println(), BoardConfig.Debug flag)
  - Key event, layer change, split connection status logging
  - Startup log (matrix size, peripherals, split config)
- ✅ Error handling improvements (Stats struct, GetStats() API)
  - Error counters for Split Rx/Tx, HID send, LED write
  - All error sites in runMaster()/runSlave() now handled

### Examples Added

- reactive-led: Reactive LED effects demonstration
- split-led: Both-side LED control with reactive effects

**Completed**: All Phase 3 features implemented

---

## Phase 4: VIA/Remap Support ✅ Completed

**Goal**: Dynamic configuration and ecosystem integration

### Implemented Features

- [x] NKRO support (N-Key Rollover)
- [x] Flash storage for persistent configuration
- [x] Macro support (key sequences, delays)
- [x] Combo keys (simultaneous key press detection)
- [x] Dynamic keymap persistence

### Examples Added

- nkro-test: NKRO functionality demonstration
- macro-test: Macro recording and playback
- combo-test: Combo key combinations
- persistent-keymap: Flash-based keymap persistence

**Completed**: All Phase 4 features implemented

---

## Future Possibilities

Features that may be considered in the future:

- 🤔 Bluetooth LE support
- 🤔 Display enhancements
- 🤔 Additional sensors
- 🤔 Audio feedback
- 🤔 Haptic feedback
- 🤔 Other MCU support
- 🤔 Web configuration tool

---

## Contributing

To contribute to feature implementation:

1. Check [CONTRIBUTING.md](CONTRIBUTING.md)
2. Review related Issues
3. Discuss feature implementation
4. Create Pull Request

---

## Versioning

Following Semantic Versioning:

- **v0.1.0**: Beta (Phase 1-4 completed) ✅

**Current version**: v0.1.0 (Beta development completed)

---

**Last updated**: 2026-02-17
