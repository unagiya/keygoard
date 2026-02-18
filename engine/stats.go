package engine

// Stats はキーボードの動作統計を保持する。
type Stats struct {
	SplitRxErrors    uint16 // split通信受信エラー
	SplitTxErrors    uint16 // split通信送信エラー
	HIDSendErrors    uint16 // HIDレポート送信エラー
	LEDWriteErrors   uint16 // LED書き込みエラー
	ScanCount        uint32 // スキャン回数
	SplitDisconnects uint16 // スレーブ切断回数
}
