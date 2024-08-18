package common

// --===================================================================================
// 　全体定義
// --===================================================================================

const (
	//画面縦横サイズ
	SCREEN_WIDTH  = 240
	SCREEN_HEIGHT = 300
	//地面のカメラに対する位置
	CAM_Y_OFFSET = -80

	//デバッグ用。高速化
	DEBUGSPEED = 1
)

const (
	//ゲージ一本ぶんの数値
	G_MAX = 10000
)

var (
	IsSmartPhone bool   //スマホならtrue
	Builddate    string //ビルド日時
	DebugStr     string //デバッグ用汎用文字列
)

//ゲームパラメータ。セーブ対象
type GameParam struct {
	Rank   float64 //ランク
	H_Rank float64 //最高ランク

	//ランク100到達フラグ
	Happy bool

	//敵出現位置制御
	CamLastX  int
	NextEnemy int

	//敵スターテス
	EN_BASE_AT int
	EN_BASE_HP int
	EN_BASE_DF int

	//デバッグ用、スタータスをcsvに記録
	DebugRecRank int
}
