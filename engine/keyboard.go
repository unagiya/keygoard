//go:build tinygo

package engine

import (
	"machine"
	hidkb "machine/usb/hid/keyboard"
	"time"

	"github.com/unagiya/keygoard/matrix"
)

// Keyboard はキーボードエンジンです。
// マトリクススキャン・キーマップ参照・HID 送信を統合します。
type Keyboard struct {
	scanner   *matrix.Scanner
	keymap    *Keymap
	prevState [matrix.RowCount][matrix.ColCount]bool
}

// New は Keyboard を生成します。
func New(scanner *matrix.Scanner, km *Keymap) *Keyboard {
	return &Keyboard{
		scanner: scanner,
		keymap:  km,
	}
}

// Init はハードウェアを初期化します。
// USB エニュメレーション完了まで待機してからスキャンを開始できる状態にします。
func (kb *Keyboard) Init() {
	kb.scanner.Init()

	// USB エニュメレーション完了まで待機
	// fixme: InitEndpointComplete ポーリングに変更することで待機時間を最小化できる
	_ = machine.USBDev
	time.Sleep(500 * time.Millisecond)
}

// Tick は 1 スキャンサイクルを実行します。メインループから毎回呼び出します。
// キー状態に変化があった場合のみ HID レポートを送信します。
func (kb *Keyboard) Tick() {
	state, changed := kb.scanner.Scan()
	if !changed {
		return
	}

	for row := 0; row < matrix.RowCount; row++ {
		for col := 0; col < matrix.ColCount; col++ {
			curr := state[row][col]
			if curr == kb.prevState[row][col] {
				continue
			}

			kc := kb.keymap.Layer0[row][col]
			if kc == 0 {
				// キー未割り当て
				continue
			}

			if curr {
				if err := hidkb.Keyboard.Down(hidkb.Keycode(kc)); err != nil {
					println("keygoard: Down error:", err.Error())
				}
			} else {
				if err := hidkb.Keyboard.Up(hidkb.Keycode(kc)); err != nil {
					println("keygoard: Up error:", err.Error())
				}
			}
		}
	}

	kb.prevState = state
}
