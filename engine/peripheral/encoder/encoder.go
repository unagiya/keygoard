//go:build tinygo

package encoder

import (
	"machine"

	"github.com/unagiya/keygoard/keycode"
)

// Config はエンコーダーの設定です。
type Config struct {
	// PinA は A 信号ピンです。
	PinA machine.Pin
	// PinB は B 信号ピンです。
	PinB machine.Pin
	// KeyCW は時計回りに割り当てるキーコードです。
	KeyCW keycode.Keycode
	// KeyCCW は反時計回りに割り当てるキーコードです。
	KeyCCW keycode.Keycode
}

// Encoder はロータリーエンコーダーの入力を処理します。
// engine.Peripheral インターフェースを実装します。
type Encoder struct {
	pinA    machine.Pin
	pinB    machine.Pin
	keyCW   keycode.Keycode
	keyCCW  keycode.Keycode
	decoder Decoder
}

// New は Config からエンコーダーを生成します。
func New(cfg *Config) *Encoder {
	return &Encoder{
		pinA:   cfg.PinA,
		pinB:   cfg.PinB,
		keyCW:  cfg.KeyCW,
		keyCCW: cfg.KeyCCW,
	}
}

// Init は GPIO ピンをプルアップ入力に設定します。
func (e *Encoder) Init() {
	e.pinA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	e.pinB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

// Tick は A/B 信号を読み取り、回転方向に応じたキーコードを返します。
// 回転がない場合は keycode.None を返します。
func (e *Encoder) Tick() keycode.Keycode {
	dir := e.decoder.Update(e.pinA.Get(), e.pinB.Get())
	switch dir {
	case DirCW:
		return e.keyCW
	case DirCCW:
		return e.keyCCW
	default:
		return keycode.None
	}
}

// OnLayerChange はレイヤー変更通知を受け取ります。
// エンコーダーはレイヤー変更に反応しないため、何もしません。
func (e *Encoder) OnLayerChange(layer int) {}
