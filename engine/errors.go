package engine

import "errors"

// ErrKeyOverflow はキー同時押し上限（6KRO）超過を表します。
var ErrKeyOverflow = errors.New("keygoard: key overflow")
