package engine

import (
	"time"

	"github.com/unagiya/keygoard/hid"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/matrix"
	"github.com/unagiya/keygoard/peripheral/encoder"
	"github.com/unagiya/keygoard/peripheral/joystick"
	"github.com/unagiya/keygoard/peripheral/led"
	"github.com/unagiya/keygoard/peripheral/oled"
	"github.com/unagiya/keygoard/split"
	"github.com/unagiya/keygoard/storage"
)

// Keyboard represents the main keyboard engine.
type Keyboard struct {
	config *BoardConfig

	// Core components
	scanner  *matrix.Scanner
	hidDev   *hid.CompositeHID
	layerMgr *LayerState
	keymap   *Keymap

	// Peripherals
	joy        *joystick.Joystick
	display    *oled.Display
	encoders   []*encoder.Encoder
	ledControl *led.Controller

	// Split keyboard
	splitMaster *split.Master
	splitSlave  *split.Slave

	// Macro (Phase 4)
	macroMgr *MacroManager

	// Combo (Phase 4)
	comboDet *ComboDetector

	// Storage (Phase 4)
	store *storage.Storage

	// State
	running           bool
	encoderEvents     []keycode.Keycode // Buffer for encoder events
	ledTick           uint32            // Tick counter for LED effects
	tapDetector       *TapDetector      // Tap detection for TT/LT (Phase 2 M4)
	prevState         [][]bool          // Previous matrix state for detecting changes
	lastSyncLayer     uint8             // OLED同期の最後に送信したレイヤー
	layerText         [9]byte           // "Layer: X" の固定バイト列
	debug             *DebugLogger      // デバッグロガー
	stats             Stats             // エラーカウンタ
	wasSlaveConnected bool              // Split接続状態変化検出用
	prevLayer         uint8             // レイヤー変更検出用
}

// NewKeyboard creates a new keyboard instance.
func NewKeyboard(config *BoardConfig, keymap *Keymap) (*Keyboard, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Validate keymap
	if err := keymap.Validate(config); err != nil {
		return nil, err
	}

	kb := &Keyboard{
		config:        config,
		keymap:        keymap,
		layerMgr:      NewLayerState(),
		tapDetector:   NewTapDetector(DefaultTapConfig()),
		macroMgr:      NewMacroManager(),
		lastSyncLayer: 0xFF, // 初回は必ず送信されるように無効値で初期化
		debug:         newDebugLogger(config.Debug),
	}

	// コンボ検出器の初期化
	kb.comboDet = NewComboDetector(config.ComboWindow)

	// マクロ登録
	for _, md := range config.Macros {
		if err := kb.macroMgr.RegisterMacro(md.ID, md.Steps); err != nil {
			return nil, err
		}
	}

	// コンボ登録
	for _, cd := range config.Combos {
		keys := make([]keycode.Keycode, cd.Count)
		for i := uint8(0); i < cd.Count; i++ {
			keys[i] = cd.Keys[i]
		}
		if err := kb.comboDet.RegisterCombo(keys, cd.Output); err != nil {
			return nil, err
		}
	}

	// "Layer: X" の固定バイト列を事前構築
	copy(kb.layerText[:], "Layer: 0")

	// Initialize previous state tracking
	rows, cols := config.GetMatrixSize()
	kb.prevState = make([][]bool, rows)
	for i := range kb.prevState {
		kb.prevState[i] = make([]bool, cols)
	}

	// Initialize matrix scanner
	kb.scanner = matrix.NewScanner(&matrix.Config{
		Rows:         config.Rows,
		Cols:         config.Cols,
		MatrixType:   config.MatrixType,
		DebounceTime: config.DebounceTime,
	})

	// Initialize HID (only for master or non-split)
	if !config.IsSplitSlave() {
		kb.hidDev = hid.NewCompositeHIDWithMode(config.DefaultHIDMode)
	}

	// Initialize joystick if configured
	if config.HasJoystick() {
		kb.joy = joystick.New(config.Joystick)
	}

	// Initialize OLED if configured (only for master or non-split)
	if config.HasOLED() && !config.IsSplitSlave() {
		display, err := oled.New(config.OLED)
		if err != nil {
			return nil, err
		}
		kb.display = display
	}

	// Initialize encoders if configured (Phase 2)
	if config.HasEncoders() {
		kb.encoders = make([]*encoder.Encoder, len(config.Encoders))
		for i, cfg := range config.Encoders {
			kb.encoders[i] = encoder.New(cfg)
		}
		kb.encoderEvents = make([]keycode.Keycode, 0, 4)
	}

	// Initialize LED if configured (Phase 2)
	if config.HasLED() {
		// Unifiedモード時はスレーブLED分もバッファに確保（Phase 3 M2）
		if config.IsSplitMaster() && config.Split.LEDSyncMode == LEDSyncUnified {
			config.LED.SlaveLEDCount = config.Split.SlaveLEDCount
		}
		ledCtrl, err := led.New(config.LED)
		if err != nil {
			return nil, err
		}
		kb.ledControl = ledCtrl
	}

	// Initialize split keyboard
	if config.HasSplit() {
		if config.IsSplitMaster() {
			master, err := split.NewMaster(&split.MasterConfig{
				UART:      config.Split.UART,
				BaudRate:  config.Split.BaudRate,
				Timeout:   config.Split.Timeout,
				SlaveRows: config.Split.SlaveRows,
				SlaveCols: config.Split.SlaveCols,
			})
			if err != nil {
				return nil, err
			}
			kb.splitMaster = master

			// スレーブキーイベントコールバック設定（Phase 3 M2）
			// スレーブ側のrow/colにマスター行数オフセットを加算してリアクティブエフェクトに通知
			masterRows := len(config.Rows)
			master.SetKeyEventCallback(func(row, col uint8, pressed bool) {
				kb.notifyLEDKeyEvent(int(row)+masterRows, int(col), pressed)
			})

			// Independentモード時はスレーブにモードを通知（Phase 3 M2）
			if config.Split.LEDSyncMode == LEDSyncIndependent {
				master.SendLEDMode(uint8(LEDSyncIndependent))
			}
		} else {
			slave, err := split.NewSlave(&split.SlaveConfig{
				UART:     config.Split.UART,
				BaudRate: config.Split.BaudRate,
				Interval: config.Split.SendInterval,
			})
			if err != nil {
				return nil, err
			}
			kb.splitSlave = slave

			// Setup callbacks for slave-side LED/OLED sync (Phase 2 M3)
			if kb.ledControl != nil {
				slave.SetLEDCallback(func(colors []byte) {
					kb.ledControl.SetColorsFromBytes(colors)
				})
			}

			if kb.display != nil {
				slave.SetOLEDCallback(func(syncType byte, payload []byte) {
					kb.handleOLEDSync(syncType, payload)
				})
			}
		}
	}

	return kb, nil
}

// Run starts the keyboard main loop.
func (kb *Keyboard) Run() error {
	kb.running = true

	// 起動ログ
	if kb.debug.enabled {
		kb.debug.Log("INIT", "keygoard started")
		rows, cols := kb.config.GetMatrixSize()
		if kb.config.MatrixType == COL2ROW {
			println("[INIT]", "matrix", rows, "x", cols, "COL2ROW")
		} else {
			println("[INIT]", "matrix", rows, "x", cols, "ROW2COL")
		}
		if kb.config.HasSplit() {
			if kb.config.IsSplitMaster() {
				println("[INIT]", "split master (slave", kb.config.Split.SlaveRows, "x", kb.config.Split.SlaveCols, ")")
			} else {
				kb.debug.Log("INIT", "split slave")
			}
		}
		kb.logPeripherals()
	}

	// Slave mode: simplified loop
	if kb.config.IsSplitSlave() {
		return kb.runSlave()
	}

	// Master or non-split mode: full keyboard loop
	return kb.runMaster()
}

// logPeripherals は接続されている周辺機器のログを出力する。
func (kb *Keyboard) logPeripherals() {
	if !kb.debug.enabled {
		return
	}
	msg := "peripherals:"
	if kb.ledControl != nil {
		msg += " LED"
	}
	if kb.display != nil {
		msg += " OLED"
	}
	if kb.joy != nil {
		msg += " JOY"
	}
	if len(kb.encoders) > 0 {
		msg += " ENC"
	}
	kb.debug.Log("INIT", msg)
}

// runMaster runs the main loop for master or non-split keyboard.
func (kb *Keyboard) runMaster() error {
	scanTicker := time.NewTicker(kb.config.ScanInterval)
	reportTicker := time.NewTicker(kb.config.USBReportRate)

	for kb.running {
		select {
		case <-scanTicker.C:
			kb.stats.ScanCount++

			// Scan matrix
			kb.scanner.Scan()

			// Update split master if configured
			if kb.splitMaster != nil {
				if err := kb.splitMaster.Update(); err != nil {
					kb.stats.SplitRxErrors++
					kb.debug.LogError("SPLIT", err)
				}

				// Split接続状態の変化検出
				connected := kb.splitMaster.IsSlaveConnected()
				if connected != kb.wasSlaveConnected {
					kb.debug.LogSplitStatus(connected)
					if !connected {
						kb.stats.SplitDisconnects++
					}
					kb.wasSlaveConnected = connected
				}
			}

			// Process keys and update HID reports
			kb.processKeys()

		case <-reportTicker.C:
			// Send HID reports
			if kb.hidDev != nil {
				if err := kb.hidDev.SendReports(); err != nil {
					kb.stats.HIDSendErrors++
					kb.debug.LogError("HID", err)
				}
			}

			// Update OLED display
			kb.updateDisplay()

			// Update LED (Phase 2)
			if kb.ledControl != nil {
				kb.ledTick++
				if err := kb.ledControl.Update(kb.ledTick); err != nil {
					kb.stats.LEDWriteErrors++
					kb.debug.LogError("LED", err)
				}

				// Sync LED to slave if split master (Phase 2 M3, Phase 3 M2)
				if kb.splitMaster != nil && kb.config.Split.LEDSyncMode == LEDSyncUnified {
					colors := kb.ledControl.GetSlaveColors()
					if colors != nil {
						if err := kb.splitMaster.SendLEDSync(colors); err != nil {
							kb.stats.SplitTxErrors++
							kb.debug.LogError("SPLIT", err)
						}
					}
				}
			}

			// Sync OLED to slave if split master (Phase 2 M3)
			if kb.display != nil && kb.splitMaster != nil {
				// レイヤー変化時のみOLED同期を送信
				layer := kb.layerMgr.Current()
				if layer != kb.lastSyncLayer {
					kb.lastSyncLayer = layer
					// 固定バイト列の最後の文字だけ更新
					kb.layerText[7] = '0' + layer
					if err := kb.splitMaster.SendOLEDText(0, 0, string(kb.layerText[:8])); err != nil {
						kb.stats.SplitTxErrors++
						kb.debug.LogError("SPLIT", err)
					}
				}
			}
		}
	}

	return nil
}

// runSlave runs the main loop for slave keyboard.
func (kb *Keyboard) runSlave() error {
	ticker := time.NewTicker(kb.config.ScanInterval)

	// スレーブ側のprevState初期化（Phase 3 M2）
	rows, cols := kb.config.GetMatrixSize()
	slavePrevState := make([][]bool, rows)
	for i := range slavePrevState {
		slavePrevState[i] = make([]bool, cols)
	}

	for kb.running {
		<-ticker.C
		kb.stats.ScanCount++

		// Scan matrix
		state := kb.scanner.Scan()

		// キー状態変化を検出してマスターに送信（Phase 3 M2）
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				if r < len(state) && c < len(state[r]) && state[r][c] != slavePrevState[r][c] {
					slavePrevState[r][c] = state[r][c]
					kb.debug.LogKeyEvent(r, c, state[r][c])
					if err := kb.splitSlave.SendKeyEvent(uint8(r), uint8(c), state[r][c]); err != nil {
						kb.stats.SplitTxErrors++
						kb.debug.LogError("SPLIT", err)
					}

					// Independentモード時はローカルのリアクティブエフェクトに通知
					if kb.splitSlave.GetLEDSyncMode() == 1 {
						kb.notifyLEDKeyEvent(r, c, state[r][c])
					}
				}
			}
		}

		// Read joystick if available
		var joyX, joyY int8
		if kb.joy != nil {
			joyX, joyY = kb.joy.Read()
		}

		// Send to master
		if kb.joy != nil {
			if err := kb.splitSlave.SendMatrixAndJoystick(state, joyX, joyY); err != nil {
				kb.stats.SplitTxErrors++
				kb.debug.LogError("SPLIT", err)
			}
		} else {
			if err := kb.splitSlave.SendMatrixState(state); err != nil {
				kb.stats.SplitTxErrors++
				kb.debug.LogError("SPLIT", err)
			}
		}

		// Receive LED/OLED sync from master (Phase 2 M3)
		if kb.splitSlave != nil {
			if err := kb.splitSlave.Update(); err != nil {
				kb.stats.SplitRxErrors++
				kb.debug.LogError("SPLIT", err)
			}
		}

		// Update LED if configured (Phase 2 M3, Phase 3 M2)
		if kb.ledControl != nil {
			kb.ledTick++
			if err := kb.ledControl.Update(kb.ledTick); err != nil {
				kb.stats.LEDWriteErrors++
				kb.debug.LogError("LED", err)
			}
		}
	}

	return nil
}

// processKeys processes the current key state and updates HID reports.
func (kb *Keyboard) processKeys() {
	// Clear previous reports
	if kb.hidDev != nil {
		kb.hidDev.Clear()
	}

	// マクロ実行ループ
	if kb.macroMgr.IsRunning() {
		kc, stepType, _ := kb.macroMgr.Update()
		if kc != keycode.KC_NO && kb.hidDev != nil {
			switch stepType {
			case MacroStepKeyDown:
				kb.processMacroKeycode(kc, true)
			case MacroStepKeyUp:
				kb.processMacroKeycode(kc, false)
			case MacroStepKeyTap:
				kb.processMacroKeycode(kc, true)
				// Tapは次のレポートでキーアップ（1サイクル遅延）
				// 簡易実装: 即座にレポート送信してからキーアップ
			}
		}
	}

	// Clear encoder events buffer
	kb.encoderEvents = kb.encoderEvents[:0]

	// Process encoders (Phase 2)
	for _, enc := range kb.encoders {
		event := enc.Update()
		if event != encoder.EventNone {
			kc := enc.GetKeycode(event)
			if kc != keycode.KC_NO {
				kb.encoderEvents = append(kb.encoderEvents, kc)
			}
		}
	}

	// Get current layer
	currentLayer := kb.layerMgr.Current()

	// Process encoder events as key presses
	for _, kc := range kb.encoderEvents {
		kb.processKeycode(kc)
	}

	// コンボタイムアウト処理: 保留キーを通常処理に回す
	pendingKeys, pendingCount := kb.comboDet.Update()
	for i := uint8(0); i < pendingCount; i++ {
		kb.processKeycode(pendingKeys[i])
	}

	// Update tap detector (Phase 2 M4)
	tapActions := kb.tapDetector.Update()
	for _, action := range tapActions {
		if action.Type == TapActionActivateLayer {
			kb.layerMgr.ActivateMO(action.Layer)
		}
	}

	// Process local matrix with tap detection
	rows, cols := kb.config.GetMatrixSize()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			currentState := kb.scanner.GetKeyState(r, c)
			prevState := kb.prevState[r][c]

			// Detect state change
			if currentState != prevState {
				kb.debug.LogKeyEvent(r, c, currentState)
				kb.notifyLEDKeyEvent(r, c, currentState)

				// コンボ検出（押下時のみ）
				if currentState {
					kc := kb.keymap.GetKey(currentLayer, r, c)
					if kc == keycode.KC_TRNS {
						for l := int(currentLayer) - 1; l >= 0; l-- {
							kc = kb.keymap.GetKey(uint8(l), r, c)
							if kc != keycode.KC_TRNS {
								break
							}
						}
					}
					comboOutput, consumed := kb.comboDet.ProcessKeyPress(kc)
					if comboOutput != keycode.KC_NO {
						kb.processKeycode(comboOutput)
						kb.prevState[r][c] = currentState
						continue
					}
					if consumed {
						kb.prevState[r][c] = currentState
						continue
					}
				} else {
					// キー離上時はコンボに通知
					kc := kb.keymap.GetKey(currentLayer, r, c)
					if kc == keycode.KC_TRNS {
						for l := int(currentLayer) - 1; l >= 0; l-- {
							kc = kb.keymap.GetKey(uint8(l), r, c)
							if kc != keycode.KC_TRNS {
								break
							}
						}
					}
					kb.comboDet.ProcessKeyRelease(kc)
				}

				kb.processKeyWithTap(currentLayer, r, c, currentState)
				kb.prevState[r][c] = currentState
			} else if currentState {
				// Key is still pressed (no change)
				kb.processKey(currentLayer, r, c)
			}
		}
	}

	// Process slave matrix if split master
	if kb.splitMaster != nil {
		slaveState := kb.splitMaster.GetSlaveState()
		for r := 0; r < len(slaveState); r++ {
			for c := 0; c < len(slaveState[r]); c++ {
				if slaveState[r][c] {
					// Slave keys: offset row by master row count
					kb.processKey(currentLayer, r+rows, c)
				}
			}
		}
	}

	// Update joystick axes
	if kb.joy != nil && kb.hidDev != nil {
		x, y := kb.joy.Read()
		kb.hidDev.Gamepad().GetReport().X = x
		kb.hidDev.Gamepad().GetReport().Y = y
	}

	// レイヤー変更検出
	newLayer := kb.layerMgr.Current()
	if newLayer != kb.prevLayer {
		kb.debug.LogLayerChange(kb.prevLayer, newLayer)
		kb.prevLayer = newLayer
	}

	// Update slave joystick if split master
	if kb.splitMaster != nil && kb.hidDev != nil {
		slaveX, slaveY := kb.splitMaster.GetSlaveJoystick()
		// For Phase 1, we'll combine or use slave as second axis
		// This needs proper implementation for 2-axis support
		_ = slaveX
		_ = slaveY
		// TODO: Implement 2-axis gamepad support
	}
}

// processKey processes a single key press.
func (kb *Keyboard) processKey(layer uint8, row, col int) {
	kc := kb.keymap.GetKey(layer, row, col)

	// Handle transparent keys
	if kc == keycode.KC_TRNS {
		// Look up in lower layers
		for l := int(layer) - 1; l >= 0; l-- {
			kc = kb.keymap.GetKey(uint8(l), row, col)
			if kc != keycode.KC_TRNS {
				break
			}
		}
	}

	kb.processKeycode(kc)
}

// processKeyWithTap processes a key press/release with tap detection (Phase 2 M4).
func (kb *Keyboard) processKeyWithTap(layer uint8, row, col int, pressed bool) {
	kc := kb.keymap.GetKey(layer, row, col)

	// Handle transparent keys
	if kc == keycode.KC_TRNS {
		for l := int(layer) - 1; l >= 0; l-- {
			kc = kb.keymap.GetKey(uint8(l), row, col)
			if kc != keycode.KC_TRNS {
				break
			}
		}
	}

	// Process through tap detector
	effectiveKc, shouldSend := kb.tapDetector.ProcessKey(kc, row, col, pressed)

	if shouldSend && effectiveKc != keycode.KC_NO {
		kb.processKeycode(effectiveKc)
	}
}

// processKeycode processes a keycode (from matrix or encoder).
func (kb *Keyboard) processKeycode(kc keycode.Keycode) {
	// Handle different keycode types
	switch {
	case kc.IsBasicKey():
		// Basic key: add to keyboard report
		if kb.hidDev != nil {
			kb.hidDev.Keyboard().AddKey(uint8(kc))
		}

	case kc.IsModifier():
		// Modifier: set modifier bit
		if kb.hidDev != nil {
			mask := keycode.ModifierMask(kc)
			kb.hidDev.Keyboard().SetModifier(mask)
		}

	case kc.IsLayerOp():
		// Layer operation
		kb.handleLayerOp(kc)

	case kc.IsGamepadButton():
		// Gamepad button
		if kb.hidDev != nil {
			idx := keycode.GamepadButtonIndex(kc)
			kb.hidDev.Gamepad().GetReport().SetButton(idx, true)
		}

	case kc.IsSystemControl():
		// システム制御
		kb.handleSystemControl(kc)

	case kc.IsMacro():
		// マクロ実行
		idx := kc.MacroIndex()
		if idx >= 0 {
			kb.macroMgr.Trigger(uint8(idx))
		}
	}
}

// handleLayerOp handles layer operation keycodes.
func (kb *Keyboard) handleLayerOp(kc keycode.Keycode) {
	opType, layer, ok := keycode.DecodeLayerOp(kc)
	if !ok {
		return
	}

	switch opType {
	case 0x00: // MO
		kb.layerMgr.ActivateMO(layer)
	case 0x01: // TG
		kb.layerMgr.ToggleTG(layer)
	case 0x02: // TT
		// TODO: Implement tap detection
		// For Phase 1, treat as MO
		kb.layerMgr.ActivateMO(layer)
	case 0x03: // LT
		// TODO: Implement tap detection
		// For Phase 1, treat as MO
		kb.layerMgr.ActivateMO(layer)
	}
}

// processMacroKeycode はマクロステップのキーコードをHIDレポートに反映する。
func (kb *Keyboard) processMacroKeycode(kc keycode.Keycode, down bool) {
	if kb.hidDev == nil {
		return
	}
	if kc.IsBasicKey() {
		if down {
			kb.hidDev.Keyboard().AddKey(uint8(kc))
		}
		// キーアップはClear()で処理されるため、ここでは何もしない
	} else if kc.IsModifier() {
		if down {
			mask := keycode.ModifierMask(kc)
			kb.hidDev.Keyboard().SetModifier(mask)
		}
	}
}

// RegisterMacro はマクロを登録する公開API。
func (kb *Keyboard) RegisterMacro(id uint8, steps []MacroStep) error {
	return kb.macroMgr.RegisterMacro(id, steps)
}

// RegisterCombo はコンボを登録する公開API。
func (kb *Keyboard) RegisterCombo(keys []keycode.Keycode, output keycode.Keycode) error {
	return kb.comboDet.RegisterCombo(keys, output)
}

// handleSystemControl はシステム制御キーコードを処理する。
func (kb *Keyboard) handleSystemControl(kc keycode.Keycode) {
	if kb.hidDev == nil {
		return
	}
	switch kc {
	case keycode.KC_NKRO_TOGGLE:
		kb.hidDev.Keyboard().ToggleMode()
		kb.debug.Log("HID", "NKRO toggled")
	case keycode.KC_NKRO_ON:
		kb.hidDev.Keyboard().SetMode(hid.HIDModeNKRO)
		kb.debug.Log("HID", "NKRO on")
	case keycode.KC_NKRO_OFF:
		kb.hidDev.Keyboard().SetMode(hid.HIDMode6KRO)
		kb.debug.Log("HID", "NKRO off")
	}
}

// updateDisplay updates the OLED display with current status.
func (kb *Keyboard) updateDisplay() {
	if kb.display == nil {
		return
	}

	// Get current state
	layer := kb.layerMgr.Current()
	var joyX, joyY int8
	if kb.joy != nil {
		joyX, joyY = kb.joy.Read()
	}

	// Get lock states (from keyboard report)
	capsLock := false
	numLock := false
	// TODO: Track lock states from host

	// Update display
	kb.display.ShowStatus(layer, joyX, joyY, capsLock, numLock)
}

// notifyLEDKeyEvent はキー状態変化をLEDコントローラに通知する。
func (kb *Keyboard) notifyLEDKeyEvent(row, col int, pressed bool) {
	if kb.ledControl == nil {
		return
	}
	ledIdx := kb.ledControl.GetLEDIndex(row, col)
	kb.ledControl.NotifyKeyEvent(led.KeyEvent{
		Row:      uint8(row),
		Col:      uint8(col),
		LEDIndex: ledIdx,
		Pressed:  pressed,
		Tick:     kb.ledTick,
	})
}

// Stop stops the keyboard.
func (kb *Keyboard) Stop() {
	kb.running = false
}

// SetStorage はFlashストレージを設定する。
// machine.Flashを使って初期化したstorageを渡す。
func (kb *Keyboard) SetStorage(s *storage.Storage) {
	kb.store = s
}

// SaveSettings は現在の設定をFlashに保存する。
// マトリクススキャン中に呼び出さないこと（Flash書き込み中はCPU停止の可能性あり）。
func (kb *Keyboard) SaveSettings() error {
	if kb.store == nil {
		return ErrNotRunning
	}

	// キーマップのシリアライズ
	rows, cols := kb.config.GetMatrixSize()
	keymapData := storage.SerializeKeymap(kb.keymap.layers, kb.keymap.LayerCount(), rows, cols)

	// 設定のシリアライズ
	var hidMode uint8
	if kb.hidDev != nil {
		hidMode = uint8(kb.hidDev.Keyboard().GetMode())
	}
	cfg := &storage.ConfigData{
		HIDMode: hidMode,
	}
	configData := storage.SerializeConfig(cfg)

	// マクロは空データ（将来拡張）
	macroData := make([]byte, 0)

	return kb.store.SaveAll(keymapData, configData, macroData)
}

// LoadSettings はFlashから設定を読み込む。
func (kb *Keyboard) LoadSettings() error {
	if kb.store == nil {
		return ErrNotRunning
	}

	if !kb.store.IsValid() {
		return ErrNotRunning
	}

	_, configData, _, err := kb.store.LoadAll()
	if err != nil {
		return err
	}

	// 設定の復元
	cfg, err := storage.DeserializeConfig(configData)
	if err != nil {
		return err
	}

	if kb.hidDev != nil {
		kb.hidDev.Keyboard().SetMode(hid.HIDMode(cfg.HIDMode))
	}

	return nil
}

// ResetSettings はFlash上の設定をリセットする（消去）。
func (kb *Keyboard) ResetSettings() error {
	if kb.store == nil {
		return ErrNotRunning
	}

	// 空データで上書き
	return kb.store.SaveAll(make([]byte, 0), make([]byte, 0), make([]byte, 0))
}

// SetLEDEffect sets the LED effect (Phase 2).
func (kb *Keyboard) SetLEDEffect(effect led.Effect) {
	if kb.ledControl != nil {
		kb.ledControl.SetEffect(effect)
	}
}

// GetLEDController returns the LED controller (Phase 2).
func (kb *Keyboard) GetLEDController() *led.Controller {
	return kb.ledControl
}

// handleOLEDSync handles OLED sync messages from master (Phase 2 M3).
func (kb *Keyboard) handleOLEDSync(syncType byte, payload []byte) {
	if kb.display == nil {
		return
	}

	switch syncType {
	case 0x00: // Clear
		kb.display.Clear()

	case 0x01: // Text at position
		if len(payload) >= 3 {
			row := payload[0]
			col := payload[1]
			textLen := int(payload[2])
			if len(payload) >= 3+textLen {
				text := string(payload[3 : 3+textLen])
				kb.display.SetText(int(row), int(col), text)
			}
		}

	case 0x02: // Raw buffer
		// TODO: Implement raw buffer update if needed
	}
}

// CalibrateJoystick recalibrates the joystick center position (Phase 2 M4).
// The joystick should be at rest (centered) when this is called.
func (kb *Keyboard) CalibrateJoystick() {
	if kb.joy != nil {
		kb.joy.Recalibrate()
	}
}

// GetJoystickRaw returns the raw joystick ADC values (Phase 2 M4).
// Useful for calibration UI.
func (kb *Keyboard) GetJoystickRaw() (x, y uint16) {
	if kb.joy != nil {
		return kb.joy.GetRaw()
	}
	return 0, 0
}

// GetJoystickCenter returns the calibrated joystick center (Phase 2 M4).
func (kb *Keyboard) GetJoystickCenter() (x, y uint16) {
	if kb.joy != nil {
		return kb.joy.GetCenter()
	}
	return 2048, 2048
}

// GetStats はキーボードの動作統計を返す。
func (kb *Keyboard) GetStats() Stats {
	return kb.stats
}
