//go:build tinygo

package engine

import (
	"machine"
	"machine/usb"
	hidkb "machine/usb/hid/keyboard"
	"time"

	"github.com/unagiya/keygoard/internal/layer"
	"github.com/unagiya/keygoard/internal/matrix"
	"github.com/unagiya/keygoard/internal/peripheral/encoder"
	"github.com/unagiya/keygoard/internal/peripheral/joystick"
	"github.com/unagiya/keygoard/internal/peripheral/led"
	"github.com/unagiya/keygoard/internal/peripheral/oled"
	"github.com/unagiya/keygoard/internal/tap"
	"github.com/unagiya/keygoard/keycode"
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
	resolver        *layer.Resolver
	tap             tap.Detector
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

	cols := pinsToCols(cfg.ColPins)
	rows := pinsToRows(cfg.RowPins)

	kb := &Keyboard{
		scanner:     matrix.New(cols, rows),
		keymap:      cfg.Keymap,
		productName: name,
		resolver:    layer.NewResolver(&cfg.Keymap.Layers),
	}

	// 標準ペリフェラルを内部生成
	if cfg.Encoder != nil {
		kb.addPeripheral(encoder.New(
			machine.Pin(cfg.Encoder.PinA),
			machine.Pin(cfg.Encoder.PinB),
			cfg.Encoder.KeyCW,
			cfg.Encoder.KeyCCW,
		))
	}
	if cfg.LED != nil {
		kb.addPeripheral(led.New(
			machine.Pin(cfg.LED.Pin),
			cfg.LED.Count,
			led.DeviceType(cfg.LED.Type),
			cfg.LED.LayerColors,
		))
	}
	if cfg.OLED != nil {
		bus := i2cBus(cfg.OLED.Bus)
		kb.addPeripheral(oled.New(
			bus,
			machine.Pin(cfg.OLED.SDA),
			machine.Pin(cfg.OLED.SCL),
			cfg.OLED.Address,
			cfg.OLED.Width,
			cfg.OLED.Height,
			oled.Rotation(cfg.OLED.Rotation),
		))
	}
	if cfg.Joystick != nil {
		sens := cfg.Joystick.Sensitivity
		if sens == 0 {
			sens = 5
		}
		dz := cfg.Joystick.DeadZone
		if dz == 0 {
			dz = 3000
		}
		kb.addPeripheral(joystick.New(
			machine.Pin(cfg.Joystick.PinX),
			machine.Pin(cfg.Joystick.PinY),
			machine.Pin(cfg.Joystick.PinButton),
			cfg.Joystick.EnableButton,
			cfg.Joystick.ButtonKey,
			sens,
			dz,
			cfg.Joystick.InvertX,
			cfg.Joystick.InvertY,
		))
	}

	// カスタムペリフェラルを追加
	for _, p := range cfg.Peripherals {
		kb.addPeripheral(p)
	}

	return kb
}

// addPeripheral は周辺機器を登録します。上限を超えた場合は無視します。
func (kb *Keyboard) addPeripheral(p Peripheral) {
	if kb.peripheralCount >= MaxPeripherals {
		return
	}
	kb.peripherals[kb.peripheralCount] = p
	kb.peripheralCount++
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
	for l := layer.MaxLayers - 1; l >= 0; l-- {
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
func (kb *Keyboard) notifyLayerChange(l int) {
	kb.prevTopLayer = l
	for i := 0; i < kb.peripheralCount; i++ {
		kb.peripherals[i].OnLayerChange(l)
	}
}

// pinsToCols は Pin スライスをマトリクス列ピン配列に変換します。
func pinsToCols(pins []Pin) [matrix.ColCount]machine.Pin {
	var cols [matrix.ColCount]machine.Pin
	for i := 0; i < len(pins) && i < matrix.ColCount; i++ {
		cols[i] = machine.Pin(pins[i])
	}
	return cols
}

// pinsToRows は Pin スライスをマトリクス行ピン配列に変換します。
func pinsToRows(pins []Pin) [matrix.RowCount]machine.Pin {
	var rows [matrix.RowCount]machine.Pin
	for i := 0; i < len(pins) && i < matrix.RowCount; i++ {
		rows[i] = machine.Pin(pins[i])
	}
	return rows
}

// i2cBus は I2CBus を machine.I2C に変換します。
func i2cBus(bus I2CBus) *machine.I2C {
	switch bus {
	case I2C1:
		return machine.I2C1
	default:
		return machine.I2C0
	}
}
