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
// マトリクススキャン・レイヤー解決・HID 送信・周辺機器を統合します。
type Keyboard struct {
	scanner         *matrix.Scanner
	keymap          *Keymap
	productName     string
	prevState       [matrix.RowCount][matrix.ColCount]bool
	activeKeys      [matrix.RowCount][matrix.ColCount]keycode.Keycode
	resolver        *Resolver
	tap             TapDetector
	peripherals     [MaxPeripherals]Peripheral
	peripheralCount int
	prevTopLayer    int
}

// New は Config からキーボードエンジンを生成します。
func New(cfg *Config) *Keyboard {
	name := cfg.ProductName
	if name == "" {
		name = defaultProductName
	}

	kb := &Keyboard{
		scanner:     cfg.Scanner,
		keymap:      cfg.Keymap,
		productName: name,
		resolver:    NewResolver(cfg.Keymap),
	}

	for i := 0; i < len(cfg.Peripherals) && i < MaxPeripherals; i++ {
		kb.peripherals[i] = cfg.Peripherals[i]
		kb.peripheralCount++
	}

	return kb
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

	// 周辺機器を初期化
	for i := 0; i < kb.peripheralCount; i++ {
		kb.peripherals[i].Init()
	}

	// USB エニュメレーション完了まで待機
	// fixme: InitEndpointComplete ポーリングに変更することで待機時間を最小化できる
	_ = machine.USBDev
	time.Sleep(500 * time.Millisecond)

	// 初期レイヤーを通知
	kb.notifyLayerChange(0)
}

// tick は 1 スキャンサイクルを実行します。
// 周辺機器のポーリング・タップ検出・キー状態変化の検出を行い、
// 必要に応じて HID レポートを送信します。
func (kb *Keyboard) tick() {
	prevTop := kb.topLayer()

	// 周辺機器の Tick（エンコーダー等のポーリング）
	for i := 0; i < kb.peripheralCount; i++ {
		if kc := kb.peripherals[i].Tick(); kc != keycode.None {
			// 周辺機器からのキーコードはタップとして送信（Down → Up）
			if err := hidkb.Keyboard.Down(hidkb.Keycode(kc)); err != nil {
				println("keygoard: Down error:", err.Error())
			}
			if err := hidkb.Keyboard.Up(hidkb.Keycode(kc)); err != nil {
				println("keygoard: Up error:", err.Error())
			}
		}
	}

	// タップカウンタをインクリメント
	kb.tap.Advance()

	// タイムアウトチェック: pending → holding に遷移したらレイヤーを有効化
	for row := 0; row < matrix.RowCount; row++ {
		for col := 0; col < matrix.ColCount; col++ {
			if timedOut, kc := kb.tap.CheckTimeout(row, col); timedOut {
				kb.resolver.Activate(kc.Layer())
				kb.activeKeys[row][col] = kc
			}
		}
	}

	state, changed := kb.scanner.Scan()
	if !changed {
		kb.checkLayerChange(prevTop)
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
	kb.checkLayerChange(prevTop)
}

// handlePress はキー押下時の処理を行います。
// レイヤーアクションキーはレイヤー状態を変更し、通常キーは HID レポートを送信します。
// LT/TT キーはタップ検出に委ねます。
func (kb *Keyboard) handlePress(row, col int, kc keycode.Keycode) {
	switch {
	case kc.IsMO():
		kb.resolver.Activate(kc.Layer())
		kb.activeKeys[row][col] = kc

	case kc.IsTG():
		kb.resolver.Toggle(kc.Layer())
		kb.activeKeys[row][col] = kc

	case kc.IsLT(), kc.IsTT():
		// タップ/ホールド判定待ち（activeKeys はタイムアウト時に設定）
		kb.tap.Press(row, col, kc)

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
// タップ検出中のキーはタップ判定を行い、それ以外は activeKeys に基づいて解除処理を行います。
func (kb *Keyboard) handleRelease(row, col int) {
	// タップ判定チェック
	if wasPending, kc := kb.tap.Release(row, col); wasPending {
		switch {
		case kc.IsLT():
			// LT タップ: タップキーコードを Down → Up（即時送信）
			tapKc := kc.TapKeycode()
			if err := hidkb.Keyboard.Down(hidkb.Keycode(tapKc)); err != nil {
				println("keygoard: Down error:", err.Error())
			}
			if err := hidkb.Keyboard.Up(hidkb.Keycode(tapKc)); err != nil {
				println("keygoard: Up error:", err.Error())
			}
		case kc.IsTT():
			// TT タップ: レイヤートグル
			kb.resolver.Toggle(kc.Layer())
		}
		kb.activeKeys[row][col] = keycode.None
		return
	}

	active := kb.activeKeys[row][col]
	if active == keycode.None {
		return
	}

	switch {
	case active.IsMO():
		kb.resolver.Deactivate(active.Layer())

	case active.IsLT(), active.IsTT():
		// ホールド中だった LT/TT のリリース: レイヤーを無効化
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

// topLayer は最上位のアクティブレイヤー番号を返します。
func (kb *Keyboard) topLayer() int {
	for l := MaxLayers - 1; l >= 0; l-- {
		if kb.resolver.IsActive(l) {
			return l
		}
	}
	return 0
}

// checkLayerChange はレイヤー変更を検出し、変化があれば周辺機器に通知します。
func (kb *Keyboard) checkLayerChange(prevTop int) {
	currTop := kb.topLayer()
	if currTop != prevTop {
		kb.notifyLayerChange(currTop)
	}
}

// notifyLayerChange は全周辺機器にレイヤー変更を通知します。
func (kb *Keyboard) notifyLayerChange(layer int) {
	kb.prevTopLayer = layer
	for i := 0; i < kb.peripheralCount; i++ {
		kb.peripherals[i].OnLayerChange(layer)
	}
}
