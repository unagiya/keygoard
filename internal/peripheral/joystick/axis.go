// Package joystick はアナログジョイスティックの入力処理を提供します。
package joystick

// center は ADC 16bit の中心値です。
const center = 32768

// AxisMapper は ADC 生値をマウス移動量に変換します。
// デッドゾーンと感度を考慮して int8 の移動量を算出します。
type AxisMapper struct {
	deadZone  int32
	maxOutput int32
	invert    bool
}

// NewAxisMapper は AxisMapper を生成します。
// sensitivity は 1〜10 の範囲で、範囲外の場合はクランプされます。
func NewAxisMapper(deadZone uint16, sensitivity uint8, invert bool) AxisMapper {
	sens := max(sensitivity, 1)
	sens = min(sens, 10)
	return AxisMapper{
		deadZone:  int32(deadZone),
		maxOutput: int32(sens) * 3,
		invert:    invert,
	}
}

// MapDelta は ADC 生値（0〜65535）をマウス移動量（int8）に変換します。
// デッドゾーン内の値は 0 を返します。
func (m *AxisMapper) MapDelta(raw uint16) int8 {
	deviation := int32(raw) - center

	// デッドゾーン判定
	if deviation > -m.deadZone && deviation < m.deadZone {
		return 0
	}

	// デッドゾーンを差し引いた有効偏差を算出
	var effective int32
	if deviation > 0 {
		effective = deviation - m.deadZone
	} else {
		effective = deviation + m.deadZone
	}

	// 有効範囲で正規化
	activeRange := int32(center) - m.deadZone
	if activeRange <= 0 {
		return 0
	}
	delta := effective * m.maxOutput / activeRange

	// クランプ
	delta = min(delta, 127)
	delta = max(delta, -128)

	if m.invert {
		delta = -delta
	}

	return int8(delta)
}
