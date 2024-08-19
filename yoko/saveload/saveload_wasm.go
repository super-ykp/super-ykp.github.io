//go:build js && wasm
// +build js,wasm

// --===================================================================================
// 　wasmにビルドされるときのみ以下が採用される
// --===================================================================================

package saveload

import (
	"syscall/js"
)

// クッキーに文字列を出力
func saveValue(key string, value string) error {
	jsFunc := js.Global().Get("WriteCookie")

	jsFunc.Invoke(js.ValueOf(key), js.ValueOf(value))
	return nil
}

// クッキーから文字列を取得
func loadValue(key string) (string, error) {
	jsFunc := js.Global().Get("getCookie")
	valu := jsFunc.Invoke(js.ValueOf(key))

	return valu.String(), nil
}
