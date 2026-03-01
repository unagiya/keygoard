// Package encoder はロータリーエンコーダーの入力処理を提供します。
package encoder

// Direction は回転方向を表します。
type Direction uint8

const (
	// DirNone は回転なし（または不正遷移）を表します。
	DirNone Direction = iota
	// DirCW は時計回りを表します。
	DirCW
	// DirCCW は反時計回りを表します。
	DirCCW
)

// dirTable は 2 相エンコーダーの状態遷移テーブルです。
// インデックスは [前回状態][現在状態] で、状態は A<<1|B でエンコードします。
// CW:  00→01→11→10→00
// CCW: 00→10→11→01→00
var dirTable = [4][4]Direction{
	//        curr=00   01     10     11
	/* 00 */ {DirNone, DirCW, DirCCW, DirNone},
	/* 01 */ {DirCCW, DirNone, DirNone, DirCW},
	/* 10 */ {DirCW, DirNone, DirNone, DirCCW},
	/* 11 */ {DirNone, DirCCW, DirCW, DirNone},
}

// Decoder は 2 相エンコーダーの回転方向を判定します。
// 前回の A/B 信号状態を保持し、状態遷移テーブルで O(1) 判定を行います。
type Decoder struct {
	prev uint8
}

// Update は現在の A/B 信号を受け取り、回転方向を返します。
func (d *Decoder) Update(a, b bool) Direction {
	curr := boolToUint8(a)<<1 | boolToUint8(b)
	dir := dirTable[d.prev][curr]
	d.prev = curr
	return dir
}

// boolToUint8 は bool を 0 または 1 に変換します。
func boolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}
