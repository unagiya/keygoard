//go:build tinygo

package joystick

import (
	"machine"
	"machine/usb/hid/mouse"

	"github.com/unagiya/keygoard/keycode"
)

// Joystick はアナログジョイスティックを制御します。
// engine.Peripheral インターフェースを実装します。
type Joystick struct {
	adcX      machine.ADC
	adcY      machine.ADC
	axisX     AxisMapper
	axisY     AxisMapper
	pinButton machine.Pin
	buttonKey keycode.Keycode
	hasButton bool
	prevBtn   bool
}

// New はジョイスティックを生成します。
func New(
	pinX, pinY, pinButton machine.Pin,
	enableButton bool,
	buttonKey keycode.Keycode,
	sensitivity uint8,
	deadZone uint16,
	invertX, invertY bool,
) *Joystick {
	return &Joystick{
		adcX:      machine.ADC{Pin: pinX},
		adcY:      machine.ADC{Pin: pinY},
		axisX:     NewAxisMapper(deadZone, sensitivity, invertX),
		axisY:     NewAxisMapper(deadZone, sensitivity, invertY),
		pinButton: pinButton,
		buttonKey: buttonKey,
		hasButton: enableButton,
	}
}

// Init は ADC ハードウェアとピンを初期化します。
func (j *Joystick) Init() {
	machine.InitADC()
	j.adcX.Configure(machine.ADCConfig{})
	j.adcY.Configure(machine.ADCConfig{})
	if j.hasButton {
		j.pinButton.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
}

// Tick は ADC 値を読み取りマウス移動を送信します。
// ボタンが有効でキーコードが設定されている場合はそのキーコードを返します。
// それ以外は keycode.None を返します。
func (j *Joystick) Tick() keycode.Keycode {
	// 軸処理: ADC 読み取り → マウス移動
	dx := j.axisX.MapDelta(j.adcX.Get())
	dy := j.axisY.MapDelta(j.adcY.Get())
	if dx != 0 || dy != 0 {
		mouse.Port().Move(int(dx), int(dy))
	}

	// ボタン処理: エッジ検出
	if j.hasButton {
		// プルアップ入力: Low = 押下
		pressed := !j.pinButton.Get()
		rising := pressed && !j.prevBtn
		j.prevBtn = pressed

		if rising {
			if j.buttonKey == keycode.None {
				mouse.Port().Click(mouse.Left)
			} else {
				return j.buttonKey
			}
		}
	}

	return keycode.None
}

// OnLayerChange はレイヤー変更通知を受け取ります。
// ジョイスティックはレイヤー変更に反応しないため、何もしません。
func (j *Joystick) OnLayerChange(layer int) {}
