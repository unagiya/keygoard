package main

import (
	"machine"
	"time"

	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/keycode"
	"github.com/unagiya/keygoard/peripheral/led"
)

func main() {
	// スレーブ側設定
	// Unifiedモード時: マスターからLEDカラーを受信して表示
	// Independentモード時: 自身のエフェクトで独立動作
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

		// LED (WS2812) — スレーブ側12個
		LED: &led.Config{
			Pin:   machine.GP16,
			Count: 12,
		},

		// Split keyboard configuration
		Split: &engine.SplitConfig{
			IsMaster:     false,
			UART:         machine.UART1,
			BaudRate:     460800,
			Timeout:      5 * time.Millisecond,
			SendInterval: 1 * time.Millisecond,
			SlaveRows:    4,
			SlaveCols:    3,
		},
	}

	// Keymap（スレーブでは未使用だが初期化に必要）
	km := engine.NewKeymap([16][][]keycode.Keycode{
		{
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
			{keycode.KC_NO, keycode.KC_NO, keycode.KC_NO},
		},
	})

	// Create keyboard
	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// スレーブ側ではLEDエフェクトはマスターから同期される（Unifiedモード）
	// Independentモード時は、マスターからMsgTypeLEDModeを受信後、
	// 自身のエフェクトで動作する

	// Run keyboard (slave mode)
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
