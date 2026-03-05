# joystick パッケージ仕様

## パッケージパス

`github.com/unagiya/keygoard/internal/peripheral/joystick`（非公開: 外部 import 不可）

## 概要

アナログジョイスティックの入力処理を提供するパッケージです。ADC ピン 2 本（X/Y 軸）からアナログ値を読み取り、マウスカーソル移動に変換します。`engine.Peripheral` インターフェースを実装しています。

利用者は `engine.JoystickConfig` 経由で設定し、`engine.New()` が内部でこのパッケージを呼び出します。

マウス移動には TinyGo 標準の `machine/usb/hid/mouse` パッケージを使用し、keyboard + mouse の Composite HID として動作します。

## ファイル構成

| ファイル          | machine 依存 | 役割                                   |
|-------------------|:------------:|----------------------------------------|
| `axis.go`         | なし         | `AxisMapper` — ADC→移動量マッピング   |
| `axis_test.go`    | なし         | AxisMapper のユニットテスト            |
| `joystick.go`     | あり         | `Joystick` 構造体（ADC + mouse.Move） |

## API

### Joystick

```go
// engine/keyboard.go から呼び出される
js := joystick.New(
    pinX, pinY, pinButton machine.Pin,
    enableButton bool,
    buttonKey keycode.Keycode,
    sensitivity uint8,
    deadZone uint16,
    invertX, invertY bool,
)
```

`engine.Peripheral` インターフェースを実装:

| メソッド               | 動作                                                   |
|------------------------|------------------------------------------------------|
| `Init()`               | `machine.InitADC()` + ADC ピン設定、ボタンピンをプルアップ入力に設定 |
| `Tick() Keycode`       | ADC 読み取り → `mouse.Port().Move()` + ボタンエッジ検出 |
| `OnLayerChange(int)`   | no-op（レイヤー変更に反応しない）                     |

`Tick()` の動作:

1. ADC X/Y 値を読み取り
2. `AxisMapper.MapDelta()` で移動量（dx, dy）を算出
3. dx != 0 || dy != 0 なら `mouse.Port().Move(dx, dy)` を呼び出し
4. ボタン有効時: エッジ検出（立ち上がり）で `mouse.Port().Click(mouse.Left)` またはキーコードを返却

**注意:** `Init()` で `machine.InitADC()` を呼ぶ必要がある。これがないと `ADC.Get()` がブロックし USB HID が動作しなくなる。

### AxisMapper

machine 非依存の軸マッピングロジック。ユニットテスト可能。

```go
m := NewAxisMapper(deadZone uint16, sensitivity uint8, invert bool)
delta := m.MapDelta(raw uint16) // ADC 生値 → int8 移動量
```

**マッピング処理:**

```
ADC 生値 (0〜65535)
  → 中心値(32768)からの偏差
  → デッドゾーン判定（閾値内 → 0）
  → 有効偏差（デッドゾーンを差し引き）
  → 正規化（有効範囲で maxOutput にスケール）
  → クランプ（-128〜127）
  → 反転（invert が true の場合）
```

感度と maxOutput の関係: `maxOutput = sensitivity * 3`

## テスト

```bash
# AxisMapper のユニットテスト（machine 非依存）
make test
```

テストケース:
- 中心値 → 0
- デッドゾーン内 → 0
- デッドゾーン境界外 → 非ゼロ
- 最大傾倒 → maxOutput 付近
- 軸反転 → 符号反転
- 感度変更 → 出力範囲変化
- 感度クランプ（0→1、255→10）
- デッドゾーン 0 / 大きすぎる値
- 正負対称性
