//go:build js && wasm
// +build js,wasm

// --===================================================================================
// 　wasmにビルドされるときのみ以下が採用される
// --===================================================================================

package sendjs

import (
	"syscall/js"
)

// クッキーに文字列を出力
func Send(name, val string) {
	jsFunc := js.Global().Get("WriteCookie")

	jsFunc.Invoke(js.ValueOf(name), js.ValueOf(val))
}

// クッキーから文字列を取得
func Get(name string) string {
	jsFunc := js.Global().Get("getCookie")
	valu := jsFunc.Invoke(js.ValueOf(name))
	return valu.String()
}
