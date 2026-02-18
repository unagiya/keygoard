package engine

// DebugLogger はTinyGo向けの軽量デバッグロガー。
// enabled=false時は全メソッドがノーオペレーションとなる。
type DebugLogger struct {
	enabled bool
}

// newDebugLogger は新しいDebugLoggerを生成する。
func newDebugLogger(enabled bool) *DebugLogger {
	return &DebugLogger{enabled: enabled}
}

// Log はタグ付きメッセージを出力する。
func (d *DebugLogger) Log(tag, msg string) {
	if !d.enabled {
		return
	}
	println("["+tag+"]", msg)
}

// LogError はタグ付きエラーメッセージを出力する。
func (d *DebugLogger) LogError(tag string, err error) {
	if !d.enabled {
		return
	}
	println("["+tag+"]", "error:", err.Error())
}

// LogKeyEvent はキーイベントを出力する。
func (d *DebugLogger) LogKeyEvent(row, col int, pressed bool) {
	if !d.enabled {
		return
	}
	if pressed {
		println("[KEY]", "R:", row, "C:", col, "pressed")
	} else {
		println("[KEY]", "R:", row, "C:", col, "released")
	}
}

// LogLayerChange はレイヤー変更を出力する。
func (d *DebugLogger) LogLayerChange(from, to uint8) {
	if !d.enabled {
		return
	}
	println("[LAYER]", int(from), "->", int(to))
}

// LogSplitStatus はSplit接続状態の変化を出力する。
func (d *DebugLogger) LogSplitStatus(connected bool) {
	if !d.enabled {
		return
	}
	if connected {
		println("[SPLIT]", "connected")
	} else {
		println("[SPLIT]", "disconnected")
	}
}
