package engine

import "github.com/unagiya/keygoard/matrix"

// Config はキーボードエンジンの設定です。
type Config struct {
	// Scanner はマトリクススキャナーです。
	Scanner *matrix.Scanner

	// Keymap はキーマップです。
	Keymap *Keymap

	// ProductName は USB デバイスとして OS に表示される名前です。
	// 空文字列の場合はデフォルト値 "keygoard" を使用します。
	// ASCII のみ、最大 126 文字。
	ProductName string
}
