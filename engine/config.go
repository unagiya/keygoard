//go:build tinygo

package engine

import (
	"image/color"

	"github.com/unagiya/keygoard/keycode"
)

// Config はキーボードエンジンの設定です。
type Config struct {
	// ProductName は USB デバイスとして OS に表示される名前です。
	// 空文字列の場合はデフォルト値 "keygoard" を使用します。
	// ASCII のみ、最大 126 文字。
	ProductName string

	// ColPins はマトリクスの列ピン（出力）の GPIO 番号です。
	ColPins []Pin

	// RowPins はマトリクスの行ピン（入力プルダウン）の GPIO 番号です。
	RowPins []Pin

	// Keymap はキーマップです。
	Keymap *Keymap

	// Encoder はロータリーエンコーダーの設定です。nil の場合はスキップします。
	Encoder *EncoderConfig

	// LED は RGB LED の設定です。nil の場合はスキップします。
	LED *LEDConfig

	// OLED は OLED ディスプレイの設定です。nil の場合はスキップします。
	OLED *OLEDConfig

	// Peripherals はカスタム周辺機器のリストです。
	// 標準ペリフェラル（Encoder/LED/OLED）と合わせて最大 MaxPeripherals 個まで登録できます。
	Peripherals []Peripheral
}

// EncoderConfig はロータリーエンコーダーの設定です。
type EncoderConfig struct {
	// PinA は A 信号ピンの GPIO 番号です。
	PinA Pin
	// PinB は B 信号ピンの GPIO 番号です。
	PinB Pin
	// KeyCW は時計回りに割り当てるキーコードです。
	KeyCW keycode.Keycode
	// KeyCCW は反時計回りに割り当てるキーコードです。
	KeyCCW keycode.Keycode
}

// MaxLayerColors はレイヤーごとの色設定の最大数です。
const MaxLayerColors = 4

// LEDConfig は RGB LED の設定です。
type LEDConfig struct {
	// Pin はデータピンの GPIO 番号です。
	Pin Pin
	// Count は LED の数です。
	Count uint8
	// Type は LED デバイス種別です。ゼロ値は WS2812（デフォルト）です。
	Type LEDType
	// LayerColors はレイヤーごとの LED 色です。
	LayerColors [MaxLayerColors]color.RGBA
}

// OLEDConfig は OLED ディスプレイの設定です。
type OLEDConfig struct {
	// Bus は I2C バスです。
	Bus I2CBus
	// SDA は I2C データピンの GPIO 番号です。
	SDA Pin
	// SCL は I2C クロックピンの GPIO 番号です。
	SCL Pin
	// Address は I2C アドレスです（通常 0x3C）。
	Address uint16
	// Width は画面幅（ピクセル）です。
	Width int16
	// Height は画面高（ピクセル）です。
	Height int16
	// Rotation はソフトウェアによる画面回転（90° 単位）です。
	Rotation Rotation
}
