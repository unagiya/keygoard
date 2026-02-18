package storage

// ConfigData はFlashに保存する設定データ。
type ConfigData struct {
	HIDMode       uint8 // 0=6KRO, 1=NKRO
	LEDEffectID   uint8 // 現在のLEDエフェクトID
	LEDBrightness uint8 // LED輝度（0-255）
	Reserved      [61]byte
}

// SerializeConfig は設定をバイト列にシリアライズする。
func SerializeConfig(cfg *ConfigData) []byte {
	buf := make([]byte, ConfigMaxSize)
	buf[0] = cfg.HIDMode
	buf[1] = cfg.LEDEffectID
	buf[2] = cfg.LEDBrightness
	// 残りはゼロ（reserved）
	return buf
}

// DeserializeConfig はバイト列から設定をデシリアライズする。
func DeserializeConfig(data []byte) (*ConfigData, error) {
	if len(data) < 3 {
		return nil, ErrInvalidMagic
	}
	cfg := &ConfigData{
		HIDMode:       data[0],
		LEDEffectID:   data[1],
		LEDBrightness: data[2],
	}
	return cfg, nil
}
