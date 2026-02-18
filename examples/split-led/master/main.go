package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/led"
)

func main() {
	// マスター側設定（Unifiedモード: 両側のLEDをマスターが制御）
	cfg := &engine.BoardConfig{
		// Matrix configuration
		Rows: []machine.Pin{
			machine.GP0, machine.GP1, machine.GP2, machine.GP3,
		},
		Cols: []machine.Pin{
			machine.GP4, machine.GP5, machine.GP6,
		},
		MatrixType: engine.COL2ROW,

		// Timing
		ScanInterval:  1 * time.Millisecond,
		DebounceTime:  3 * time.Millisecond,
		USBReportRate: 1 * time.Millisecond,

		// LED (WS2812) — マスター側12個
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: 12,
			// KeyToLED: マスター4行×3列 + スレーブ4行×3列 = 24キー分のマッピング
			KeyToLED: []int8{
				0, 1, 2, // row0
				3, 4, 5, // row1
				6, 7, 8, // row2
				9, 10, 11, // row3
				12, 13, 14, // slave row0（オフセット: マスター12個分）
				15, 16, 17, // slave row1
				18, 19, 20, // slave row2
				21, 22, 23, // slave row3
			},
			MaxCols: 3,
		},

		// Split keyboard configuration
		Split: &engine.SplitConfig{
			IsMaster:      true,
			UART:          machine.UART1,
			BaudRate:      460800,
			Timeout:       5 * time.Millisecond,
			SendInterval:  1 * time.Millisecond,
			SlaveRows:     4,
			SlaveCols:     3,
			SlaveLEDCount: 12,                    // スレーブ側のLED数
			LEDSyncMode:   engine.LEDSyncUnified, // マスターが両側を制御
		},
	}

	// Keymap
	km := engine.NewKeymap([16][][]keycode.Keycode{
		// Layer 0
		{
			// Master side (left)
			{keycode.KC_Q, keycode.KC_W, keycode.KC_E},
			{keycode.KC_A, keycode.KC_S, keycode.KC_D},
			{keycode.KC_Z, keycode.KC_X, keycode.KC_C},
			{keycode.KC_LCTL, keycode.MO(1), keycode.KC_SPC},
			// Slave side (right) - rows 4-7
			{keycode.KC_R, keycode.KC_T, keycode.KC_Y},
			{keycode.KC_F, keycode.KC_G, keycode.KC_H},
			{keycode.KC_V, keycode.KC_B, keycode.KC_N},
			{keycode.KC_ENT, keycode.MO(1), keycode.KC_LSFT},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// リアクティブエフェクト（FadeOut）を設定
	// Unifiedモードでは、スレーブのキー押下もマスターに転送され、
	// 両側のLEDでリアクティブエフェクトが動作する
	kb.SetLEDEffect(led.NewFadeOutEffect(
		color.RGBA{R: 0, G: 100, B: 255, A: 255},
		60,
	))

	// Run keyboard
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
