package led

import "image/color"

// maxLEDs はLEDの最大数（固定配列サイズ）。
const maxLEDs = 128

// KeyPressEffect はキー押下中に点灯、離すと消灯するエフェクト。
type KeyPressEffect struct {
	color   color.RGBA
	pressed [maxLEDs]bool
}

// NewKeyPressEffect はKeyPressEffectを作成する。
func NewKeyPressEffect(col color.RGBA) *KeyPressEffect {
	return &KeyPressEffect{color: col}
}

// Update implements Effect。
func (e *KeyPressEffect) Update(buffer []color.RGBA, tick uint32) {
	for i := range buffer {
		if e.pressed[i] {
			buffer[i] = e.color
		} else {
			buffer[i] = color.RGBA{A: 255}
		}
	}
}

// OnKeyEvent implements ReactiveEffect。
func (e *KeyPressEffect) OnKeyEvent(event KeyEvent) {
	if event.LEDIndex < 0 || int(event.LEDIndex) >= maxLEDs {
		return
	}
	e.pressed[event.LEDIndex] = event.Pressed
}

// FadeOutEffect はキー押下で点灯し、時間経過で減衰するエフェクト。
type FadeOutEffect struct {
	color     color.RGBA
	duration  uint32 // フェードアウト時間（tick数）
	active    [maxLEDs]bool
	pressTick [maxLEDs]uint32
}

// NewFadeOutEffect はFadeOutEffectを作成する。
// durationはフェードアウトに要するtick数。
func NewFadeOutEffect(col color.RGBA, duration uint32) *FadeOutEffect {
	if duration == 0 {
		duration = 60 // デフォルト: 約1秒（60fps想定）
	}
	return &FadeOutEffect{
		color:    col,
		duration: duration,
	}
}

// Update implements Effect。
func (e *FadeOutEffect) Update(buffer []color.RGBA, tick uint32) {
	for i := range buffer {
		if !e.active[i] {
			buffer[i] = color.RGBA{A: 255}
			continue
		}

		elapsed := tick - e.pressTick[i]
		if elapsed >= e.duration {
			// フェードアウト完了
			e.active[i] = false
			buffer[i] = color.RGBA{A: 255}
			continue
		}

		// 線形フェード: brightness = (duration - elapsed) / duration * 255
		remaining := e.duration - elapsed
		brightness := uint32(remaining) * 255 / uint32(e.duration)
		buffer[i] = color.RGBA{
			R: uint8(uint32(e.color.R) * brightness / 255),
			G: uint8(uint32(e.color.G) * brightness / 255),
			B: uint8(uint32(e.color.B) * brightness / 255),
			A: 255,
		}
	}
}

// OnKeyEvent implements ReactiveEffect。
func (e *FadeOutEffect) OnKeyEvent(event KeyEvent) {
	if event.LEDIndex < 0 || int(event.LEDIndex) >= maxLEDs {
		return
	}
	if event.Pressed {
		e.active[event.LEDIndex] = true
		e.pressTick[event.LEDIndex] = event.Tick
	}
}

// maxRipples は同時に存在できる波紋の最大数。
const maxRipples = 8

// rippleState は1つの波紋の状態を表す。
type rippleState struct {
	active    bool
	centerX   uint8
	centerY   uint8
	startTick uint32
}

// RippleEffect はキー押下位置から波紋が広がるエフェクト。
type RippleEffect struct {
	color    color.RGBA
	duration uint32 // 波紋の持続時間（tick数）
	speed    uint32 // 波紋の広がり速度（tick数あたりの距離）
	ripples  [maxRipples]rippleState
	ledPosX  [maxLEDs]uint8 // 各LEDのX座標
	ledPosY  [maxLEDs]uint8 // 各LEDのY座標
	ledCount int
}

// NewRippleEffect はRippleEffectを作成する。
// posX, posYは各LEDの座標配列。durationは波紋の持続tick数。
func NewRippleEffect(col color.RGBA, posX, posY []uint8, duration uint32) *RippleEffect {
	if duration == 0 {
		duration = 90
	}
	e := &RippleEffect{
		color:    col,
		duration: duration,
		speed:    2, // デフォルト速度
	}

	// LED位置を固定配列にコピー
	count := len(posX)
	if count > maxLEDs {
		count = maxLEDs
	}
	if len(posY) < count {
		count = len(posY)
	}
	e.ledCount = count
	for i := 0; i < count; i++ {
		e.ledPosX[i] = posX[i]
		e.ledPosY[i] = posY[i]
	}

	return e
}

// Update implements Effect。
func (e *RippleEffect) Update(buffer []color.RGBA, tick uint32) {
	// まず全LEDを消灯
	for i := range buffer {
		buffer[i] = color.RGBA{A: 255}
	}

	// 各アクティブな波紋を処理
	for ri := range e.ripples {
		r := &e.ripples[ri]
		if !r.active {
			continue
		}

		elapsed := tick - r.startTick
		if elapsed >= e.duration {
			r.active = false
			continue
		}

		// 波紋の現在の半径（距離）
		radius := elapsed * e.speed

		// フェードアウト係数
		fade := uint32(e.duration-elapsed) * 255 / uint32(e.duration)

		// 各LEDとの距離を計算
		limit := len(buffer)
		if e.ledCount < limit {
			limit = e.ledCount
		}
		for i := 0; i < limit; i++ {
			dx := int32(e.ledPosX[i]) - int32(r.centerX)
			dy := int32(e.ledPosY[i]) - int32(r.centerY)
			distSq := uint32(dx*dx + dy*dy)

			// 波紋の幅（リング状）
			radiusSq := radius * radius
			// リングの内側と外側の距離を判定（幅=speed*2）
			ringWidth := e.speed * 4
			ringWidthSq := ringWidth * ringWidth

			// distSqがradiusSqの近辺にあるかチェック
			var diff uint32
			if distSq > radiusSq {
				diff = distSq - radiusSq
			} else {
				diff = radiusSq - distSq
			}

			if diff > ringWidthSq {
				continue
			}

			// リング内の明るさ（中心ほど明るい）
			ringBrightness := uint32(255)
			if ringWidthSq > 0 {
				ringBrightness = (ringWidthSq - diff) * 255 / ringWidthSq
			}

			// 最終明るさ = リング明るさ * フェード
			brightness := ringBrightness * fade / 255

			// 加算ブレンド
			buffer[i] = addColor(buffer[i], color.RGBA{
				R: uint8(uint32(e.color.R) * brightness / 255),
				G: uint8(uint32(e.color.G) * brightness / 255),
				B: uint8(uint32(e.color.B) * brightness / 255),
				A: 255,
			})
		}
	}
}

// OnKeyEvent implements ReactiveEffect。
func (e *RippleEffect) OnKeyEvent(event KeyEvent) {
	if !event.Pressed || event.LEDIndex < 0 || int(event.LEDIndex) >= e.ledCount {
		return
	}

	// 空きスロットを探す
	for i := range e.ripples {
		if !e.ripples[i].active {
			e.ripples[i] = rippleState{
				active:    true,
				centerX:   e.ledPosX[event.LEDIndex],
				centerY:   e.ledPosY[event.LEDIndex],
				startTick: event.Tick,
			}
			return
		}
	}

	// 空きがない場合、最も古い波紋を上書き
	oldest := 0
	oldestTick := e.ripples[0].startTick
	for i := 1; i < maxRipples; i++ {
		if e.ripples[i].startTick < oldestTick {
			oldest = i
			oldestTick = e.ripples[i].startTick
		}
	}
	e.ripples[oldest] = rippleState{
		active:    true,
		centerX:   e.ledPosX[event.LEDIndex],
		centerY:   e.ledPosY[event.LEDIndex],
		startTick: event.Tick,
	}
}

// addColor は2つの色を加算ブレンドする（255でクランプ）。
func addColor(a, b color.RGBA) color.RGBA {
	r := uint16(a.R) + uint16(b.R)
	g := uint16(a.G) + uint16(b.G)
	bl := uint16(a.B) + uint16(b.B)
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if bl > 255 {
		bl = 255
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(bl), A: 255}
}
