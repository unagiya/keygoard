//go:build tinygo

package oled

import (
	"image/color"
	"machine"

	"github.com/unagiya/keygoard/keycode"

	"tinygo.org/x/drivers/ssd1306"
)

// Rotation はディスプレイの回転角度（90° 単位）です。
type Rotation uint8

const (
	Rotation0   Rotation = 0 // 回転なし（デフォルト）
	Rotation90  Rotation = 1 // 90° 時計回り
	Rotation180 Rotation = 2 // 180° 回転
	Rotation270 Rotation = 3 // 270° 時計回り
)

// Config は OLED の設定です。
type Config struct {
	// Bus は I2C バスです。
	Bus *machine.I2C
	// SDA は I2C データピンです。
	SDA machine.Pin
	// SCL は I2C クロックピンです。
	SCL machine.Pin
	// Address は I2C アドレスです（通常 0x3C）。
	Address uint16
	// Width は画面幅（ピクセル）です。
	Width int16
	// Height は画面高（ピクセル）です。
	Height int16
	// Rotation はソフトウェアによる画面回転（90° 単位）です。
	Rotation Rotation
}

// OLED は SSD1306 OLED ディスプレイを制御します。
// engine.Peripheral インターフェースを実装します。
type OLED struct {
	dev      *ssd1306.Device
	width    int16
	height   int16
	rotation Rotation
}

// New は Config から OLED を生成します。
func New(cfg *Config) *OLED {
	cfg.Bus.Configure(machine.I2CConfig{
		SDA:       cfg.SDA,
		SCL:       cfg.SCL,
		Frequency: 400_000,
	})

	dev := ssd1306.NewI2C(cfg.Bus)
	dev.Configure(ssd1306.Config{
		Width:   cfg.Width,
		Height:  cfg.Height,
		Address: cfg.Address,
	})

	return &OLED{
		dev:      dev,
		width:    cfg.Width,
		height:   cfg.Height,
		rotation: cfg.Rotation,
	}
}

// Init はディスプレイを初期化し、初期表示（"L0"）を描画します。
func (o *OLED) Init() {
	o.dev.ClearDisplay()
	o.OnLayerChange(0)
}

// Tick はスキャンサイクルごとに呼ばれます。
// OLED は毎サイクルの処理がないため keycode.None を返します。
func (o *OLED) Tick() keycode.Keycode {
	return keycode.None
}

// logicalSize は回転を考慮した論理的な画面サイズを返します。
// 0°/180° は物理サイズそのまま、90°/270° は幅と高さが入れ替わります。
func (o *OLED) logicalSize() (w, h int16) {
	if o.rotation == Rotation90 || o.rotation == Rotation270 {
		return o.height, o.width
	}
	return o.width, o.height
}

// transformPixel は論理座標を回転に応じた物理座標に変換します。
func (o *OLED) transformPixel(lx, ly int16) (px, py int16) {
	switch o.rotation {
	case Rotation90:
		return o.width - 1 - ly, lx
	case Rotation180:
		return o.width - 1 - lx, o.height - 1 - ly
	case Rotation270:
		return ly, o.height - 1 - lx
	default:
		return lx, ly
	}
}

// OnLayerChange はアクティブレイヤーが変化したときに表示を更新します。
// "L0"〜"L3" を画面中央に描画します。
func (o *OLED) OnLayerChange(layer int) {
	if layer < 0 || layer > 3 {
		return
	}

	o.dev.ClearBuffer()

	digit := byte('0' + layer)
	chars := [2]byte{'L', digit}

	logW, logH := o.logicalSize()

	// 2 文字分のスケール後の幅と高さ
	totalW := int16(2 * charWidth * fontScale)
	totalH := int16(charHeight * fontScale)

	// 論理画面の中央に配置
	startX := (logW - totalW) / 2
	startY := (logH - totalH) / 2

	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for ci, c := range chars {
		idx := glyphIndex(c)
		if idx < 0 {
			continue
		}
		g := glyphs[idx]
		offsetX := startX + int16(ci*charWidth*fontScale)

		for row := range charHeight {
			for col := range charWidth {
				if g[row]&(0x80>>col) != 0 {
					// fontScale 倍に拡大して描画
					for sy := range fontScale {
						for sx := range fontScale {
							lx := offsetX + int16(col*fontScale+sx)
							ly := startY + int16(row*fontScale+sy)
							px, py := o.transformPixel(lx, ly)
							o.dev.SetPixel(px, py, white)
						}
					}
				}
			}
		}
	}

	o.dev.Display()
}
