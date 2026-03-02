package engine

// Pin は GPIO ピン番号です。
// GPIO 番号を整数で指定します（例: GP5 → Pin(5)）。
type Pin uint8

// I2CBus は I2C バスの識別子です。
type I2CBus uint8

const (
	// I2C0 は I2C バス 0 です。
	I2C0 I2CBus = 0
	// I2C1 は I2C バス 1 です。
	I2C1 I2CBus = 1
)

// LEDType は LED デバイスの種別です。
type LEDType uint8

const (
	// WS2812 は WS2812 / SK6812MINI-E（RGB、3 バイト/LED）です。
	WS2812 LEDType = iota
	// SK6812 は SK6812（RGBW、4 バイト/LED）です。
	SK6812
)

// Rotation はディスプレイの回転角度（90° 単位）です。
type Rotation uint8

const (
	// Rotation0 は回転なし（デフォルト）です。
	Rotation0 Rotation = 0
	// Rotation90 は 90° 時計回りです。
	Rotation90 Rotation = 1
	// Rotation180 は 180° 回転です。
	Rotation180 Rotation = 2
	// Rotation270 は 270° 時計回りです。
	Rotation270 Rotation = 3
)
