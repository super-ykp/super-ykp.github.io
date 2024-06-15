//go:build !js && !wasm
// +build !js,!wasm

package loader

//jsに、SetURL関数を登録する-------------------
func Init() {

}

//URLが設定されたらtrue
func IsUrlStandBy() bool {
	return true
}
