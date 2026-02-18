// Package led provides RGB LED (WS2812) control.
package led

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

// Config holds LED configuration.
type Config struct {
	Pin           machine.Pin // Data pin for WS2812
	Count         int         // Number of LEDs (max 128)
	SlaveLEDCount int         // スレーブ側のLED数（Phase 3 M2、マスター側のみ使用）
	KeyToLED      []int8      // キーからLEDへのマッピング（row*MaxCols+col → LED index, -1=なし）
	MaxCols       int         // マトリクスの最大列数（KeyToLEDのインデックス計算用）
}

// Controller manages RGB LEDs.
type Controller struct {
	device        ws2812.Device
	buffer        []color.RGBA
	count         int // マスター（ローカル）側のLED数
	slaveLEDCount int // スレーブ側のLED数（Phase 3 M2）
	effect        Effect
	keyToLED      []int8 // キーからLEDへのマッピング
	maxCols       int    // マトリクスの最大列数
	writeBuffer   []byte // デバイス書き込み用事前割り当てバッファ（count*3）
	colorBuffer   []byte // GetColors/GetSlaveColors用事前割り当てバッファ
}

// New creates a new LED controller.
func New(cfg *Config) (*Controller, error) {
	// Configure WS2812
	pin := cfg.Pin
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	device := ws2812.New(pin)

	// Initialize buffer（マスター + スレーブ分を確保）
	totalCount := cfg.Count + cfg.SlaveLEDCount
	buffer := make([]color.RGBA, totalCount)

	// colorBuffer: GetColors/GetSlaveColorsで使う最大サイズを事前割り当て
	colorBufSize := cfg.Count
	if cfg.SlaveLEDCount > colorBufSize {
		colorBufSize = cfg.SlaveLEDCount
	}

	return &Controller{
		device:        device,
		buffer:        buffer,
		count:         cfg.Count,
		slaveLEDCount: cfg.SlaveLEDCount,
		effect:        NewStaticEffect(color.RGBA{R: 0, G: 0, B: 0, A: 255}), // Default: off
		keyToLED:      cfg.KeyToLED,
		maxCols:       cfg.MaxCols,
		writeBuffer:   make([]byte, cfg.Count*3),
		colorBuffer:   make([]byte, colorBufSize*3),
	}, nil
}

// SetEffect sets the current effect.
func (c *Controller) SetEffect(effect Effect) {
	c.effect = effect
}

// Update updates the LED effect and writes to the LEDs.
// Should be called regularly (e.g., every frame).
// エフェクトはバッファ全体（マスター+スレーブ）に適用されるが、
// デバイスへの書き込みはマスター分のみ。
func (c *Controller) Update(tick uint32) error {
	// Update effect
	if c.effect != nil {
		c.effect.Update(c.buffer, tick)
	}

	// writeBufferに直接書き込み（マスター分のみ）
	c.fillWriteBuffer(c.buffer[:c.count])
	_, err := c.device.Write(c.writeBuffer)
	return err
}

// fillWriteBuffer はRGBAカラーをwriteBufferにGRB順で書き込む。
func (c *Controller) fillWriteBuffer(colors []color.RGBA) {
	for i, col := range colors {
		c.writeBuffer[i*3] = col.G // WS2812はGRB順
		c.writeBuffer[i*3+1] = col.R
		c.writeBuffer[i*3+2] = col.B
	}
}

// SetColor sets a specific LED color.
func (c *Controller) SetColor(index int, col color.RGBA) {
	if index >= 0 && index < c.count {
		c.buffer[index] = col
	}
}

// SetAll sets all LEDs to the same color.
func (c *Controller) SetAll(col color.RGBA) {
	for i := range c.buffer {
		c.buffer[i] = col
	}
}

// Clear turns off all LEDs.
func (c *Controller) Clear() {
	c.SetAll(color.RGBA{R: 0, G: 0, B: 0, A: 255})
}

// GetBuffer returns the internal color buffer.
func (c *Controller) GetBuffer() []color.RGBA {
	return c.buffer
}

// GetCount returns the number of LEDs.
func (c *Controller) GetCount() int {
	return c.count
}

// GetColors returns the current master LED colors as RGB byte array (Phase 2 M3).
// Format: [R1, G1, B1, R2, G2, B2, ...]
// 内部バッファを返すため、次の呼び出し前に使い切ること。
func (c *Controller) GetColors() []byte {
	buf := c.colorBuffer[:c.count*3]
	for i := 0; i < c.count; i++ {
		buf[i*3] = c.buffer[i].R
		buf[i*3+1] = c.buffer[i].G
		buf[i*3+2] = c.buffer[i].B
	}
	return buf
}

// GetSlaveColors はスレーブ側のLEDカラーをRGBバイト列で返す（Phase 3 M2）。
// Format: [R1, G1, B1, R2, G2, B2, ...]
// 内部バッファを返すため、次の呼び出し前に使い切ること。
func (c *Controller) GetSlaveColors() []byte {
	if c.slaveLEDCount == 0 {
		return nil
	}
	buf := c.colorBuffer[:c.slaveLEDCount*3]
	for i := 0; i < c.slaveLEDCount; i++ {
		idx := c.count + i
		buf[i*3] = c.buffer[idx].R
		buf[i*3+1] = c.buffer[idx].G
		buf[i*3+2] = c.buffer[idx].B
	}
	return buf
}

// GetMasterCount はマスター側のLED数を返す（Phase 3 M2）。
func (c *Controller) GetMasterCount() int {
	return c.count
}

// GetTotalCount はマスター+スレーブの合計LED数を返す（Phase 3 M2）。
func (c *Controller) GetTotalCount() int {
	return c.count + c.slaveLEDCount
}

// GetLEDIndex はキーのマトリクス座標に対応するLEDインデックスを返す。
// マッピングが未設定または対応するLEDがない場合は-1を返す。
func (c *Controller) GetLEDIndex(row, col int) int8 {
	if c.keyToLED == nil || c.maxCols == 0 {
		return -1
	}
	idx := row*c.maxCols + col
	if idx < 0 || idx >= len(c.keyToLED) {
		return -1
	}
	return c.keyToLED[idx]
}

// NotifyKeyEvent はキーイベントをReactiveEffectに通知する。
// 現在のエフェクトがReactiveEffectを実装していない場合は何もしない。
func (c *Controller) NotifyKeyEvent(event KeyEvent) {
	if re, ok := c.effect.(ReactiveEffect); ok {
		re.OnKeyEvent(event)
	}
}

// SetColorsFromBytes sets LED colors from RGB byte array (Phase 2 M3).
// Format: [R1, G1, B1, R2, G2, B2, ...]
func (c *Controller) SetColorsFromBytes(data []byte) {
	count := len(data) / 3
	if count > c.count {
		count = c.count
	}

	for i := 0; i < count; i++ {
		c.buffer[i] = color.RGBA{
			R: data[i*3],
			G: data[i*3+1],
			B: data[i*3+2],
			A: 255,
		}
	}

	// writeBufferに書き込んでデバイスに送信
	c.fillWriteBuffer(c.buffer[:c.count])
	c.device.Write(c.writeBuffer)
}
