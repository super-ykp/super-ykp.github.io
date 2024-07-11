package skill

const (
	S_ENEMY_FEEVER    = 0 //ねこフィーバー
	S_EXPx2           = 1 //EX2倍増
	S_EXP_RECOV_HP    = 2 //EXP獲得でHP回復
	S_EXP_TO_GOLD     = 3 //EXPをGOLD変換
	S_COINx2          = 4 //COIN2倍増
	S_COIN_TO_SP      = 5 //COINをSPに変換
	S_ATx2            = 6 //攻撃力倍
	S_ELECTRICTRIGGER = 7 //魔法。敵全体　AT依存
	S_THUNDERBOLT     = 8 //魔法。敵全体　MaxSP依存
	S_LIGHTNINGSWORD  = 9 //魔法。攻撃体制のみ　MaxSP依存
)

type SkillDetail struct {
	CoolTime   int //クールダウンタイム
	ActiveTime int //効果維持時間開始
}

type SkillStruct struct {
	Skills []SkillDetail
}

var (
	this *SkillStruct
)

// --===================================================================================
// 初期化
func Init(t *SkillStruct) {
	this = t

	this.Skills = []SkillDetail{
		{100, 100}, //ねこフィーバー
		{100, 100}, //EX2倍増
		{100, 100}, //EXP獲得でHP回復
		{50, 150},  //EXPをGOLD変換
		{100, 100}, //COIN2倍増
		{100, 100}, //COINをFORCEに変換
		{100, 50},  //攻撃力倍
		{30, 1},    //エレクトリッガー
		{30, 1},    //ライトニングボルト
		{20, 1},    //マジカル艦砲射撃

	}
}

func GetSkill(index int) SkillDetail {
	return this.Skills[index]
}
