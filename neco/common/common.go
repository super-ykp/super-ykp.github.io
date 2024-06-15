package common

const (
	SCREEN_WIDTH       = 240
	SCREEN_HEIGHT      = 320
	GameMOdeUrlWaiting = 0
	GameModeLoading    = 1
	GameModeLogo       = 2
	GameModeGame       = 3

	GameStateStartStage = 1
	GameStateNomal      = 2
	GameStateGameOver   = 3
	GameStateBossDefeat = 4
	GameStateAllClear   = 5

	STAGE_MAX = 10
)

var (
	GameMode  int
	GameState int
	Counter   int
	KeyRelese bool
	Stage     int
)

// ステージ開始
func StartStage() {
	if Stage >= STAGE_MAX {
		GameState = GameStateAllClear
		Stage = 1
		return
	}
	GameState = GameStateStartStage
	Counter = 0
	Stage++
}

// ゲームオーバー
func StartGameOver() {
	GameState = GameStateGameOver
	Counter = 0
}

// ボス爆発中モード
func StartGameStateBossDefeat() {
	GameState = GameStateBossDefeat
	Counter = 0
}
