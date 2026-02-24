//go:build tinygo

package engine

import (
	"machine"
	"machine/usb"
	hidkb "machine/usb/hid/keyboard"
	"time"

	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
)

// defaultProductName はフレームワークのデフォルト USB Product Name です。
const defaultProductName = "keygoard"

// Keyboard はキーボードエンジンです。
// マトリクススキャン・レイヤー解決・HID 送信を統合します。
type Keyboard struct {
	scanner     *matrix.Scanner
	keymap      *Keymap
	productName string
	prevState   [matrix.RowCount][matrix.ColCount]bool
	activeKeys  [matrix.RowCount][matrix.ColCount]keycode.Keycode
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

			if curr {
				kc := kb.resolver.Resolve(row, col)
				kb.handlePress(row, col, kc)
			} else {
				kb.handleRelease(row, col)
			}
		}
	}

	kb.prevState = state
}

// handlePress はキー押下時の処理を行います。
// レイヤーアクションキーはレイヤー状態を変更し、通常キーは HID レポートを送信します。
func (kb *Keyboard) handlePress(row, col int, kc keycode.Keycode) {
	switch {
	case kc.IsMO():
		kb.resolver.Activate(kc.Layer())
		kb.activeKeys[row][col] = kc

	case kc.IsTG():
		kb.resolver.Toggle(kc.Layer())
		kb.activeKeys[row][col] = kc

	case kc == keycode.None:
		// 何もしない

	default:
		// 通常キー・修飾キー
		if err := hidkb.Keyboard.Down(hidkb.Keycode(kc)); err != nil {
			println("keygoard: Down error:", err.Error())
		}
		kb.activeKeys[row][col] = kc
	}
}

// handleRelease はキーリリース時の処理を行います。
// 押下時に記録した activeKeys に基づいて適切な解除処理を行います。
func (kb *Keyboard) handleRelease(row, col int) {
	active := kb.activeKeys[row][col]
	if active == keycode.None {
		return
	}

	switch {
	case active.IsMO():
		kb.resolver.Deactivate(active.Layer())

	case active.IsTG():
		// TG はトグル済み。リリース時は何もしない

	default:
		// 通常キー・修飾キー
		if err := hidkb.Keyboard.Up(hidkb.Keycode(active)); err != nil {
			println("keygoard: Up error:", err.Error())
		}
	}

	kb.activeKeys[row][col] = keycode.None
}
