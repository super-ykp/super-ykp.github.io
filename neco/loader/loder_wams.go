//go:build js && wasm
// +build js,wasm

package loader

import (
	"syscall/js"
)

// jsに、SetURL関数を登録する-------------------
func Init() {
	js.Global().Set("SetURL", js.FuncOf(SetURL))
}

// URLが設定されたらtrue
func IsUrlStandBy() bool {
	if m_URL == "" {
		return false
	}
	return true
}

// 外部から呼ばれる関数
func SetURL(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 {
		m_URL = args[0].String()
	}
	return 1
}
