//go:build js && wasm
// +build js,wasm

// --===================================================================================
// 　wasmにビルドされるときのみ以下が採用される
// --===================================================================================
package loader

import (
	"Rip/common"
	"syscall/js"
)

// javascriptに、自作関数SetWasm関数を登録する
func Init() {
	js.Global().Set("SetWasm", js.FuncOf(SetWasm))
}

// この関数は設定されるまで無限ループで呼ばれる。URLが設定されたらtrue
func IsUrlStandBy() bool {
	if m_URL == "" {
		return false
	}
	return true
}

// javascriptから起動情報を引き取る関数。呼び出し元はjavascrpt
func SetWasm(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 {
		m_URL = args[0].String()             //カレントURL
		common.IsSmartPhone = args[1].Bool() //スマホ起動ならtrue
		common.Builddate = args[2].String()  //ビルド日付（デバッグ用)
	}
	return 1
}
