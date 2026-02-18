package led

import "image/color"

// Effect represents an LED effect.
type Effect interface {
	// Update updates the LED buffer based on the effect.
	// tick is the current time in milliseconds or frame count.
	Update(buffer []color.RGBA, tick uint32)
}

// KeyEvent はキーの押下・離上イベントを表す。
type KeyEvent struct {
	Row      uint8
	Col      uint8
	LEDIndex int8 // -1 = マッピングなし
	Pressed  bool
	Tick     uint32
}

// ReactiveEffect はキーイベントに反応するエフェクトのインターフェース。
type ReactiveEffect interface {
	Effect
	OnKeyEvent(event KeyEvent)
}

// StaticEffect displays a single static color.
type StaticEffect struct {
	color color.RGBA
}

// NewStaticEffect creates a new static color effect.
func NewStaticEffect(col color.RGBA) *StaticEffect {
	return &StaticEffect{color: col}
}

// Update implements Effect.
func (e *StaticEffect) Update(buffer []color.RGBA, tick uint32) {
	for i := range buffer {
		buffer[i] = e.color
	}
}

// SetColor changes the static color.
func (e *StaticEffect) SetColor(col color.RGBA) {
	e.color = col
}

// BreathingEffect displays a breathing effect (fade in/out).
type BreathingEffect struct {
	color color.RGBA
	speed uint32 // Speed factor (lower = slower)
}

// NewBreathingEffect creates a new breathing effect.
func NewBreathingEffect(col color.RGBA, speed uint32) *BreathingEffect {
	if speed == 0 {
		speed = 10 // Default speed
	}
	return &BreathingEffect{
		color: col,
		speed: speed,
	}
}

// Update implements Effect.
func (e *BreathingEffect) Update(buffer []color.RGBA, tick uint32) {
	// Calculate brightness using sine wave
	// Map tick to 0-255 range with sine function
	phase := (tick / e.speed) % 256
	brightness := sineTable[phase]

	// Apply brightness to color
	scaledColor := color.RGBA{
		R: uint8((uint32(e.color.R) * uint32(brightness)) / 255),
		G: uint8((uint32(e.color.G) * uint32(brightness)) / 255),
		B: uint8((uint32(e.color.B) * uint32(brightness)) / 255),
		A: 255,
	}

	for i := range buffer {
		buffer[i] = scaledColor
	}
}

// RainbowEffect displays a rainbow pattern.
type RainbowEffect struct {
	speed  uint32 // Speed factor
	offset uint32 // Current offset
}

// NewRainbowEffect creates a new rainbow effect.
func NewRainbowEffect(speed uint32) *RainbowEffect {
	if speed == 0 {
		speed = 5
	}
	return &RainbowEffect{
		speed: speed,
	}
}

// Update implements Effect.
func (e *RainbowEffect) Update(buffer []color.RGBA, tick uint32) {
	e.offset = (tick / e.speed) % 256

	for i := range buffer {
		// Calculate hue for each LED
		hue := uint8((e.offset + uint32(i)*256/uint32(len(buffer))) % 256)
		buffer[i] = hsvToRGB(hue, 255, 255)
	}
}

// hsvToRGB converts HSV to RGB color.
// h: 0-255, s: 0-255, v: 0-255
func hsvToRGB(h, s, v uint8) color.RGBA {
	if s == 0 {
		return color.RGBA{R: v, G: v, B: v, A: 255}
	}

	region := h / 43
	remainder := (h - (region * 43)) * 6

	p := (v * (255 - s)) >> 8
	q := (v * (255 - ((s * remainder) >> 8))) >> 8
	t := (v * (255 - ((s * (255 - remainder)) >> 8))) >> 8

	switch region {
	case 0:
		return color.RGBA{R: v, G: t, B: p, A: 255}
	case 1:
		return color.RGBA{R: q, G: v, B: p, A: 255}
	case 2:
		return color.RGBA{R: p, G: v, B: t, A: 255}
	case 3:
		return color.RGBA{R: p, G: q, B: v, A: 255}
	case 4:
		return color.RGBA{R: t, G: p, B: v, A: 255}
	default:
		return color.RGBA{R: v, G: p, B: q, A: 255}
	}
}

// sineTable is a pre-calculated sine wave for breathing effect (0-255).
var sineTable = [256]uint8{
	128, 131, 134, 137, 140, 143, 146, 149, 152, 155, 158, 162, 165, 167, 170, 173,
	176, 179, 182, 185, 188, 190, 193, 196, 198, 201, 203, 206, 208, 211, 213, 215,
	218, 220, 222, 224, 226, 228, 230, 232, 234, 235, 237, 238, 240, 241, 243, 244,
	245, 246, 248, 249, 250, 250, 251, 252, 253, 253, 254, 254, 254, 255, 255, 255,
	255, 255, 255, 255, 254, 254, 254, 253, 253, 252, 251, 250, 250, 249, 248, 246,
	245, 244, 243, 241, 240, 238, 237, 235, 234, 232, 230, 228, 226, 224, 222, 220,
	218, 215, 213, 211, 208, 206, 203, 201, 198, 196, 193, 190, 188, 185, 182, 179,
	176, 173, 170, 167, 165, 162, 158, 155, 152, 149, 146, 143, 140, 137, 134, 131,
	128, 124, 121, 118, 115, 112, 109, 106, 103, 100, 97, 93, 90, 88, 85, 82,
	79, 76, 73, 70, 67, 65, 62, 59, 57, 54, 52, 49, 47, 44, 42, 40,
	37, 35, 33, 31, 29, 27, 25, 23, 21, 20, 18, 17, 15, 14, 12, 11,
	10, 9, 7, 6, 5, 5, 4, 3, 2, 2, 1, 1, 1, 0, 0, 0,
	0, 0, 0, 0, 1, 1, 1, 2, 2, 3, 4, 5, 5, 6, 7, 9,
	10, 11, 12, 14, 15, 17, 18, 20, 21, 23, 25, 27, 29, 31, 33, 35,
	37, 40, 42, 44, 47, 49, 52, 54, 57, 59, 62, 65, 67, 70, 73, 76,
	79, 82, 85, 88, 90, 93, 97, 100, 103, 106, 109, 112, 115, 118, 121, 124,
}
