package common

// --===================================================================================
// 　全体定義
// --===================================================================================

const (
	SCREEN_WIDTH  = 240
	SCREEN_HEIGHT = 300

	CAM_Y_OFFSET = -80

	DEBUGSPEED = 1
)

const (
	G_MAX = 10000
)

var (
	IsFastTouch  bool
	IsSmartPhone bool
	Builddate    string
	DebugStr     string
)

type GameParam struct {
	Rank   float64
	H_Rank float64

	//敵出現位置制御
	CamLastX  int
	NextEnemy int

	EN_BASE_AT int
	EN_BASE_HP int
	EN_BASE_DF int

	DebugRecRank int

	FirstPowerUp bool
	Happy        bool
}
