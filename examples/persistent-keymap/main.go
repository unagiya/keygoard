package main

import (
	"github.com/unagiya/keygoard/engine"
	"github.com/unagiya/keygoard/storage"
	"machine"
)

// Flash永続化オフセット（RP2040の末尾領域を使用）
const flashStorageOffset = 0x1F0000

// Flash保存デモ
// 起動時にFlashから設定を読み込み、SaveSettings/LoadSettingsで
// キーマップやHIDモード設定を永続化する。
func main() {
	cfg := newBoardConfig()
	km := newKeymap()

	kb, err := engine.NewKeyboard(cfg, km)
	if err != nil {
		panic(err)
	}

	// Flashストレージを初期化
	store, err := storage.New(machine.Flash, flashStorageOffset)
	if err != nil {
		panic(err)
	}

	// キーボードにストレージを設定
	kb.SetStorage(store)

	// 保存済み設定があればロード
	if err := kb.LoadSettings(); err != nil {
		// 初回起動時は設定がないのでエラーは無視
		println("設定なし、デフォルトで起動")
	}

	// キーボードを実行
	if err := kb.Run(); err != nil {
		panic(err)
	}
}
