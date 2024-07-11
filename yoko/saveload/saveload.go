package saveload

import (
	"Yoko/common"
	"Yoko/enemy"
	"Yoko/player"
	"encoding/json"
)

// プレイヤー------------------------------------------
func LoadAllData() (*common.GameParam, *player.PstateStruct, *enemy.EnemyStruct) {
	//ゲーム--------------------------------------
	s, err := loadValue("GAME")
	if s == "" || err != nil {
		return nil, nil, nil
	}
	game, err := DeserializeGameStruct(s)
	if err != nil {
		return nil, nil, nil
	}

	//プレイヤー--------------------------------------
	s, err = loadValue("PLAYER")
	if s == "" || err != nil {
		return nil, nil, nil
	}
	plst, err := DeserializePlayerStruct(s)
	if err != nil {
		return nil, nil, nil
	}

	//敵--------------------------------------
	s, err = loadValue("ENEMY")
	if s == "" || err != nil {
		return nil, nil, nil
	}
	enst, err := DeserializeEnemyStruct(s)
	if err != nil {
		return nil, nil, nil
	}
	return &game, &plst, &enst
}

func SaveAllData(gm common.GameParam, pl player.PstateStruct, en enemy.EnemyStruct) {
	//ゲーム--------------------------------------------------
	s, err := SeriSaalizeGameStruct(gm)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}
	err = saveValue("GAME", s)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}

	//プレイヤー---------------------------------------------------
	s, err = SeriSaalizePlayerStruct(pl)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}
	err = saveValue("PLAYER", s)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}

	//敵---------------------------------------------------
	s, err = SeriSaalizeEnemyStruct(en)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}
	err = saveValue("ENEMY", s)
	if err != nil {
		common.DebugStr = "Save Error:" + err.Error()
	}
}

// -----------------------------------------------------
// ゲーム構造体をシリアライズ
func SeriSaalizeGameStruct(m common.GameParam) (string, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ゲーム構造体をデシリアライズ
func DeserializeGameStruct(data string) (common.GameParam, error) {
	var m common.GameParam
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return common.GameParam{}, err
	}
	return m, nil
}

// -----------------------------------------------------
// プレイヤー構造体をシリアライズ
func SeriSaalizePlayerStruct(m player.PstateStruct) (string, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// プレイヤー構造体をデシリアライズ
func DeserializePlayerStruct(data string) (player.PstateStruct, error) {
	var m player.PstateStruct
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return player.PstateStruct{}, err
	}
	return m, nil
}

// -----------------------------------------------------
// 敵構造体をシリアライズ
func SeriSaalizeEnemyStruct(m enemy.EnemyStruct) (string, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// 敵構造体をデシリアライズ
func DeserializeEnemyStruct(data string) (enemy.EnemyStruct, error) {
	var m enemy.EnemyStruct
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return enemy.EnemyStruct{}, err
	}
	return m, nil
}
