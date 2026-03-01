package main

import (
	"machine"

	"github.com/unagiya/keygoard/matrix"
)

// zero-kb02 のピン配置
// 参照: https://github.com/sago35/keyboards/tree/main/zero-kb02
var (
	// matrixCols は列ピン（出力）です。COL2ROW 方式で 1 本ずつ High にします。
	matrixCols = [matrix.ColCount]machine.Pin{
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO8,
	}

	// matrixRows は行ピン（プルダウン入力）です。列が High のとき押下で High を返します。
	matrixRows = [matrix.RowCount]machine.Pin{
		machine.GPIO9,
		machine.GPIO10,
		machine.GPIO11,
	}

	// encoderPinA はエンコーダーの A 信号ピンです。
	encoderPinA = machine.GPIO3

	// encoderPinB はエンコーダーの B 信号ピンです。
	encoderPinB = machine.GPIO4

	// ledPin は RGB LED（SK6812MINI-E）のデータピンです。
	ledPin = machine.GPIO1
)
