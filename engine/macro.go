package engine

import (
	"github.com/unagiya/keygoard/keycode"
)

// マクロの制限
const (
	MaxMacros     = 16 // 最大マクロ数
	MaxMacroSteps = 32 // マクロあたりの最大ステップ数
)

// MacroStepType はマクロステップの種類を表す。
type MacroStepType uint8

const (
	MacroStepEnd     MacroStepType = 0x00 // 終端
	MacroStepKeyDown MacroStepType = 0x01 // キー押下
	MacroStepKeyUp   MacroStepType = 0x02 // キー離上
	MacroStepKeyTap  MacroStepType = 0x03 // 押下+即離上
	MacroStepDelay   MacroStepType = 0x04 // 遅延（ミリ秒）
)

// MacroStep はマクロの1ステップを表す。
type MacroStep struct {
	Type  MacroStepType
	Value uint16 // キーコード値 or 遅延ミリ秒
}

// Macro は1つのマクロ定義を表す。
type Macro struct {
	steps [MaxMacroSteps]MacroStep
	count uint8 // 有効ステップ数
}

// MacroExecState はマクロ実行中の状態を保持する。
type MacroExecState struct {
	macroID    int8   // 実行中マクロID（-1=未実行）
	stepIndex  uint8  // 現在のステップインデックス
	delayTicks uint16 // 残り遅延ティック（ms単位）
}

// MacroManager はマクロの登録と実行を管理する。
type MacroManager struct {
	macros [MaxMacros]Macro
	exec   MacroExecState
}

// NewMacroManager は新しいMacroManagerを作成する。
func NewMacroManager() *MacroManager {
	return &MacroManager{
		exec: MacroExecState{macroID: -1},
	}
}

// RegisterMacro はマクロを登録する。
// idは0-15、stepsは最大32ステップ。
func (m *MacroManager) RegisterMacro(id uint8, steps []MacroStep) error {
	if id >= MaxMacros {
		return ErrInvalidConfig
	}
	if len(steps) > MaxMacroSteps {
		return ErrInvalidConfig
	}

	macro := &m.macros[id]
	macro.count = uint8(len(steps))
	for i, s := range steps {
		macro.steps[i] = s
	}
	// 残りをEndで埋める
	for i := len(steps); i < MaxMacroSteps; i++ {
		macro.steps[i] = MacroStep{Type: MacroStepEnd}
	}

	return nil
}

// Trigger はマクロの実行を開始する。
// 実行中のマクロがある場合は無視する。
func (m *MacroManager) Trigger(id uint8) {
	if id >= MaxMacros {
		return
	}
	// マクロにステップが無い場合は無視
	if m.macros[id].count == 0 {
		return
	}
	// 実行中のマクロがあれば無視
	if m.exec.macroID >= 0 {
		return
	}
	m.exec.macroID = int8(id)
	m.exec.stepIndex = 0
	m.exec.delayTicks = 0
}

// IsRunning はマクロが実行中かを返す。
func (m *MacroManager) IsRunning() bool {
	return m.exec.macroID >= 0
}

// Update はマクロの次のステップを処理する。
// 戻り値: キーコード、ステップタイプ、完了フラグ。
// 遅延中やマクロ未実行の場合はKC_NO, MacroStepEnd, falseを返す。
func (m *MacroManager) Update() (keycode.Keycode, MacroStepType, bool) {
	if m.exec.macroID < 0 {
		return keycode.KC_NO, MacroStepEnd, false
	}

	// 遅延中
	if m.exec.delayTicks > 0 {
		m.exec.delayTicks--
		return keycode.KC_NO, MacroStepDelay, false
	}

	macro := &m.macros[m.exec.macroID]

	// ステップ実行
	if m.exec.stepIndex >= macro.count {
		// マクロ完了
		m.exec.macroID = -1
		return keycode.KC_NO, MacroStepEnd, true
	}

	step := macro.steps[m.exec.stepIndex]
	m.exec.stepIndex++

	switch step.Type {
	case MacroStepEnd:
		m.exec.macroID = -1
		return keycode.KC_NO, MacroStepEnd, true

	case MacroStepDelay:
		m.exec.delayTicks = step.Value
		return keycode.KC_NO, MacroStepDelay, false

	case MacroStepKeyDown, MacroStepKeyUp, MacroStepKeyTap:
		return keycode.Keycode(step.Value), step.Type, false

	default:
		// 未知のステップタイプはスキップ
		return keycode.KC_NO, MacroStepEnd, false
	}
}
