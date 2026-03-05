package engine

import "github.com/unagiya/keygoard/keycode"

// MaxPeripherals はエンジンに登録できる周辺機器の最大数です。
const MaxPeripherals = 8

// Peripheral はエンジンに統合される周辺機器のインターフェースです。
// engine パッケージはこのインターフェース経由で周辺機器を呼び出すため、
// 具体的な周辺機器パッケージへの依存を持ちません。
type Peripheral interface {
	// Init はハードウェアを初期化します。
	Init()

	// Tick は 1 スキャンサイクルの処理を実行します。
	// 送信すべきキーコードがある場合はそのキーコードを返します。
	// 送信不要の場合は keycode.None を返します。
	Tick() keycode.Keycode

	// OnLayerChange はアクティブレイヤーが変化したときに呼び出されます。
	// layer は最上位のアクティブレイヤー番号です。
	OnLayerChange(layer int)
}
