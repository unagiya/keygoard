//go:build tinygo

package engine

import (
	"machine"
	"machine/usb"
	hidkb "machine/usb/hid/keyboard"
	"time"

	"github.com/unagiya/keygoard/matrix"
)

// defaultProductName はフレームワークのデフォルト USB Product Name です。
const defaultProductName = "keygoard"

// Keyboard はキーボードエンジンです。
// マトリクススキャン・キーマップ参照・HID 送信を統合します。
type Keyboard struct {
	scanner     *matrix.Scanner
	keymap      *Keymap
	productName string
	prevState   [matrix.RowCount][matrix.ColCount]bool
	resolver    *Resolver
}

// New は Config からキーボードエンジンを生成します。
func New(cfg *Config) *Keyboard {
	name := cfg.ProductName
	if name == "" {
		name = defaultProductName
	}

	return &Keyboard{
		scanner:     cfg.Scanner,
		keymap:      cfg.Keymap,
		productName: name,
		resolver:    NewResolver(cfg.Keymap),
	}
}

// Run はキーボードを起動します。
// ハードウェア初期化・USB エニュメレーション待機の後、スキャンループに入ります。
// この関数は戻りません。
func (kb *Keyboard) Run() {
	usb.Product = kb.productName

	kb.setup()

	for {
		kb.tick()
		time.Sleep(1 * time.Millisecond)
	}
}

// setup はハードウェアを初期化します。
// USB エニュメレーション完了まで待機してからスキャンを開始できる状態にします。
func (kb *Keyboard) setup() {
	kb.scanner.Init()

	// USB エニュメレーション完了まで待機
	// fixme: InitEndpointComplete ポーリングに変更することで待機時間を最小化できる
	_ = machine.USBDev
	time.Sleep(500 * time.Millisecond)
}

// tick は 1 スキャンサイクルを実行します。
// キー状態に変化があった場合のみ HID レポートを送信します。
func (kb *Keyboard) tick() {
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

			kc := kb.resolver.Resolve(row, col)
			if kc == 0 {
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
