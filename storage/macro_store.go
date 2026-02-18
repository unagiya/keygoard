package storage

// マクロシリアライズ定数
const (
	MacroStepsPerMacro = 32
	MacroCount         = 16
	MacroStepSize      = 3 // Type(1) + Value(2)
)

// MacroStepData はシリアライズ用マクロステップ。
type MacroStepData struct {
	Type  uint8
	Value uint16
}

// SerializeMacros はマクロ定義をバイト列にシリアライズする。
// フォーマット: [16マクロ × 32ステップ × 3バイト] = 1536バイト
func SerializeMacros(macros [MacroCount][MacroStepsPerMacro]MacroStepData) []byte {
	buf := make([]byte, MacroMaxSize)
	idx := 0
	for m := 0; m < MacroCount; m++ {
		for s := 0; s < MacroStepsPerMacro; s++ {
			step := &macros[m][s]
			buf[idx] = step.Type
			buf[idx+1] = byte(step.Value)
			buf[idx+2] = byte(step.Value >> 8)
			idx += MacroStepSize
		}
	}
	return buf
}

// DeserializeMacros はバイト列からマクロ定義をデシリアライズする。
func DeserializeMacros(data []byte) ([MacroCount][MacroStepsPerMacro]MacroStepData, error) {
	var macros [MacroCount][MacroStepsPerMacro]MacroStepData
	if len(data) < MacroMaxSize {
		return macros, ErrDataTooLarge
	}

	idx := 0
	for m := 0; m < MacroCount; m++ {
		for s := 0; s < MacroStepsPerMacro; s++ {
			macros[m][s].Type = data[idx]
			macros[m][s].Value = uint16(data[idx+1]) | uint16(data[idx+2])<<8
			idx += MacroStepSize
		}
	}
	return macros, nil
}
