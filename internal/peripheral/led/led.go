//go:build tinygo

// Package led は WS2812 互換の RGB LED 制御を提供します。
package led

import (
	"image/color"
	"machine"

	"github.com/unagiya/keygoard/keycode"

	"tinygo.org/x/drivers/ws2812"
)

// MaxLEDs は制御可能な LED の最大数です。
const MaxLEDs = 12

// DeviceType は LED デバイスの種別です。
type DeviceType uint8

const (
	// TypeWS2812 は WS2812 / SK6812MINI-E（RGB、3 バイト/LED）です。
	TypeWS2812 DeviceType = iota
	// TypeSK6812 は SK6812（RGBW、4 バイト/LED）です。
	TypeSK6812
)

// MaxLayerColors はレイヤーごとの色設定の最大数です。
const MaxLayerColors = 4

// LED は WS2812 互換の RGB LED を制御します。
// engine.Peripheral インターフェースを実装します。
type LED struct {
	pin         machine.Pin
	deviceType  DeviceType
	ws          ws2812.Device
	count       uint8
	layerColors [MaxLayerColors]color.RGBA
	buf         [MaxLEDs]color.RGBA
}

// New は LED を生成します。
func New(pin machine.Pin, count uint8, deviceType DeviceType, layerColors [MaxLayerColors]color.RGBA) *LED {
	if count > MaxLEDs {
		count = MaxLEDs
	}
	return &LED{
		pin:         pin,
		count:       count,
		deviceType:  deviceType,
		layerColors: layerColors,
	}
}

// Init はデータピンを出力に設定し、デフォルト色で全 LED を点灯します。
func (l *LED) Init() {
	l.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	switch l.deviceType {
	case TypeSK6812:
		l.ws = ws2812.NewSK6812(l.pin)
	default:
		l.ws = ws2812.NewWS2812(l.pin)
	}
	l.OnLayerChange(0)
}

// Tick はスキャンサイクルごとに呼ばれます。
// LED は毎サイクルの処理がないため keycode.None を返します。
func (l *LED) Tick() keycode.Keycode {
	return keycode.None
}

// OnLayerChange はアクティブレイヤーが変化したときに全 LED の色を更新します。
func (l *LED) OnLayerChange(layer int) {
	if layer < 0 || layer >= MaxLayerColors {
		return
	}
	c := l.layerColors[layer]
	for i := range l.buf {
		l.buf[i] = c
	}
	l.ws.WriteColors(l.buf[:l.count])
}
