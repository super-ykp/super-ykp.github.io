//go:build !js && !wasm
// +build !js,!wasm

// --===================================================================================
//　wasmにビルドされ　ない　ときのみ以下が採用される
// --===================================================================================

package saveload

import (
	"encoding/json"
	"errors"
	"os"
)

const fileName = "data.json"

// -----------------------------------------------------
func SaveData(key, value string) error {
	if err := saveValue(key, value); err != nil {
		return err
	}
	return nil
}

func LoadData(key string) string {
	// データを読み込む
	value, err := loadValue(key)
	if err != nil {
		return ""
	}
	return value
}

// --------------------------------------------------------
// SaveData はファイルにkeyとvalueを保存します。
func saveValue(key string, value string) error {
	data := make(map[string]string)

	// ファイルが存在する場合は読み込む
	if _, err := os.Stat(fileName); err == nil {
		file, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(file, &data); err != nil {
			return err
		}
	}

	// 新しいkeyであれば追加、既存のkeyであればvalueを更新
	data[key] = value

	// マップをJSONに変換
	updatedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// ファイルに書き込む
	if err := os.WriteFile(fileName, updatedData, 0644); err != nil {
		return err
	}

	return nil
}

// LoadData はファイルからkeyに対応するvalueを読み込みます。
func loadValue(key string) (string, error) {
	// ファイルが存在しない場合はエラーを返す
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	data := make(map[string]string)
	if err := json.Unmarshal(file, &data); err != nil {
		return "", err
	}

	value, exists := data[key]
	if !exists {
		return "", errors.New("key does not exist")
	}

	return value, nil
}
