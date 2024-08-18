//go:build !js && !wasm
// +build !js,!wasm

// --===================================================================================
//　wasmにビルドされ　ない　ときのみ以下が採用される
// --===================================================================================

package gdebug

import (
	"fmt"
	"os"
)

func Write(textToAppend string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ファイルのオープンに失敗しました:", err)
		return
	}
	defer file.Close()

	// ファイルに文字列を追記
	if _, err := file.WriteString(textToAppend); err != nil {
		fmt.Println("ファイルへの書き込みに失敗しました:", err)
	}
}
