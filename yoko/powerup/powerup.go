package powerup

import (
	"Yoko/common"
	"Yoko/input"
	"Yoko/mytext"
	"Yoko/player"
	"Yoko/skill"
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	P_OFF = 0
	P_ON  = 1

	pNG_OFFSET_X  = 0
	pNG_OFFSET_Y  = 224
	sPNG_OFFSET_X = 0
	sPNG_OFFSET_Y = 192
)

type Pstruct struct {
	SelType    int    //セレクト種類　1=EXP 2=GOLD 3=SP  4=掘り出し物 5=全部
	Ptype      int    //パワーアップ種類
	imageIndex int    //イメージインデックス
	tex        string //説明テキスト
}

type PowerupStruct struct {
	player       *player.PlayerStruct
	Gparam       *common.GameParam
	powerupImage []*ebiten.Image //イメージ
	PDetail      []Pstruct       //パワーアップ一覧

	SelType int
	PLIST   [3]int //選択リストに出現しているパワーアップのID

	//描画制御-------------------------
	pmode     int
	x         float32
	y         float32
	sel       [3]bool
	press     bool
	pressBack bool
	counter   int
}

var (
	this *PowerupStruct
)

// --===================================================================================
// 初期化
func Init(t *PowerupStruct, e *ebiten.Image, p *player.PlayerStruct, gp *common.GameParam) {
	this = t
	this.player = p
	this.Gparam = gp
	setDeteil()

	//パワーアップイメージ
	for i := 0; i < len(this.PDetail); i++ {
		imageIndex := this.PDetail[i].imageIndex
		offx := sPNG_OFFSET_X + (imageIndex%20)*16
		offy := sPNG_OFFSET_Y + imageIndex/20*16
		if this.PDetail[i].Ptype >= 100 {
			offx = pNG_OFFSET_X + (imageIndex%20)*16
			offy = pNG_OFFSET_Y + imageIndex/20*16
		}

		imageRect := image.Rect(offx, offy, offx+16, offy+16)
		rect := e.SubImage(imageRect).(*ebiten.Image)
		this.powerupImage = append(this.powerupImage, rect)
	}
	this.pmode = P_OFF
	this.press = false
}

func Update() {
	//通常状態。何もしない
	if this.pmode == P_OFF {
		return
	}

	//パネル座標
	this.x = float32(5)
	this.y = -float32(common.CAM_Y_OFFSET) - 10

	//パネル表示から0.5秒は応答しない
	this.counter++
	if this.counter < 30 {
		return
	}

	//離す➝タッチ➝離す、で動作
	if !this.press && this.pressBack {
		for i := 0; i < len(this.sel); i++ {
			if this.sel[i] {
				this.pmode = P_OFF
				//ここでパワーアップ決定
				PowerUp(this.PLIST[i])
				return
			}
		}
	}
	this.pressBack = this.press

	if !this.press {
		return
	}
	//セレクト開始
	tx, ty, b := input.GetTouchPos()

	//選択箇所
	for i := 0; i < 3; i++ {
		py := this.y + 1 + (8*6+2)*float32(i)
		this.sel[i] = false
		if b && 0 < float32(tx) && tx < common.SCREEN_WIDTH && py < float32(ty) && float32(ty) < py+(8*6+2) {
			this.sel[i] = true
		}
	}

}

// 描画
func Draw(screen *ebiten.Image) {
	if this.pmode == P_OFF {
		return
	}
	c := color.RGBA{0, 0, 0, 255}
	switch this.SelType {
	case 1:
		c = color.RGBA{168, 255, 255, 255}
	case 2:
		c = color.RGBA{255, 206, 16, 255}
	case 3:
		c = color.RGBA{56, 255, 56, 255}
	}
	vector.DrawFilledRect(screen, this.x, this.y, common.SCREEN_WIDTH-10, (8*6+2)*3, c, false)
	for i := 0; i < 3; i++ {
		py := this.y + 1 + (8*6+2)*float32(i)
		drawPowerUp(screen, this.x, py, this.sel[i], i)
	}
}

// パワーアップパネル1つ
func drawPowerUp(screen *ebiten.Image, x, y float32, sel bool, index int) {
	c := color.RGBA{64, 64, 64, 255}
	if sel {
		c = color.RGBA{255, 255, 255, 255}
	}
	vector.DrawFilledRect(screen, x+1, y, common.SCREEN_WIDTH-10-2, (8 * 6), c, false)

	pIndex := this.PLIST[index]
	imageIndex := pIndex
	drawImageOption := ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(-8, -8)
	drawImageOption.GeoM.Scale(2, 2)
	drawImageOption.GeoM.Translate(float64(x)+16+8, float64(y)+16+8)
	screen.DrawImage(this.powerupImage[imageIndex], &drawImageOption)
	st := this.PDetail[pIndex].tex
	mytext.DrawN(screen, int(x)+50, int(y)+8, st)
}

// 選択開始
func StartPowerUp(stype int) {
	//パワーアップのリスト抽選
	SelectPowerUp(stype)

	this.pmode = P_ON
	this.counter = 0
	this.press = false
	this.pressBack = false
	for i := 0; i < len(this.sel); i++ {
		this.sel[i] = false
	}
}

// ユーザー選択中
func Select(press bool) bool {
	if this.pmode == P_OFF {
		return true
	}
	this.press = press
	return false
}

// -----------------------------------------------------------------------
// パワーアップ種類
const (
	P_AT    = 100
	P_MAXHP = 101

	P_HP_RECV1 = 102
	P_HP_RECV2 = 103
	P_HP_RECV3 = 104

	P_AT_T1   = 105
	P_AT_T2   = 106
	P_AT_T3   = 107
	P_MAXSP_T = 108

	P_RESHPUP = 109
	P_RESSPUP = 110

	P_RATE_EXP   = 111
	P_RATE_COIN  = 112
	P_RATE_FORCE = 113

	P_COOLDOWN = 114
	P_SPIKE    = 115

	P_CAKE     = 117
	P_ALLSKILL = 118
)

// パワーアップ種類設定
func setDeteil() {
	this.PDetail = []Pstruct{
		{1, P_AT, 0, "力アップ\n力AT+2"},
		{1, P_MAXHP, 1, "体力アップ\nMaxHP+10"},
		{1, P_RESHPUP, 10, "HP回復力アップ\nRES HP+10"},
		{1, P_RATE_EXP, 12, "EXPゲージ増加率アップ"},

		{0, P_HP_RECV1, 3, "ポテチ\n RES HPの値回復"},
		{0, P_HP_RECV2, 4, "ビスケット\n RES SPの1.5倍の値回復"},
		{0, P_HP_RECV3, 5, "ドーナッツ\n RES SPの2倍の値回復"},

		//{0, P_AT_T1, 6, "力のたね\nAT+1"},
		{0, P_AT_T2, 7, "力の葉\nAT+2"},
		{0, P_AT_T3, 8, "力の実\nAT+3"},

		{2, P_RATE_COIN, 13, "COINゲージ増加率アップ"},
		{2, P_RESHPUP, 10, "HP回復力アップ\nRES HP+5"},
		{2, P_MAXHP, 1, "体力アップ\nMaxHP+5"},

		{4, P_CAKE, 18, "掘り出し物\nケーキ\nHP,SP全回復"},
		{4, P_ALLSKILL, 19, "掘り出し物\nスキル全使用可能に!\n実行中なら最大値に"},

		{3, P_MAXSP_T, 9, "知恵の実\nMaxSP+4"},
		{3, P_RESSPUP, 11, "SP回復力アップ\nRES SP+1"},
		{3, P_RATE_FORCE, 14, "FORCEゲージ増加率アップ"},
		{3, P_SPIKE, 16, "スパイク\nSPで防御時、ダメージを返す"},
		{3, P_COOLDOWN, 15, "待ち状態スキルの有効化\n全スキル待ちカウント0に"},

		{5, skill.S_ENEMY_FEEVER, skill.S_ENEMY_FEEVER, "スキル\nねこ大発生"},
		{5, skill.S_EXPx2, skill.S_EXPx2, "スキル\n経験値取得時2倍入るよ"},
		{5, skill.S_COINx2, skill.S_COINx2, "スキル\nコイン取得時2倍入るよ"},
		{5, skill.S_EXP_TO_GOLD, skill.S_EXP_TO_GOLD, "スキル\n経験値をコインに変えるよ"},
		{5, skill.S_COIN_TO_FORCE, skill.S_COIN_TO_FORCE, "スキル\nコインをFORCEに変えるよ"},
		{5, skill.S_EXP_RECOV_HP, skill.S_EXP_RECOV_HP, "スキル\n経験値でHPが回復するよ"},
		{5, skill.S_ATx2, skill.S_ATx2, "スキル\n攻撃力が倍になるよ"},
		{5, skill.S_ELECTRICTRIGGER, skill.S_ELECTRICTRIGGER, "スキル\n攻撃体制の敵を打ち落とすよ\n威力はAT依存"},
		{5, skill.S_THUNDERBOLT, skill.S_THUNDERBOLT, "スキル\n全体攻撃\n威力はMaxSP依存"},
		{5, skill.S_LIGHTNINGSWORD, skill.S_LIGHTNINGSWORD, "スキル\nマジカル艦砲射撃\n上空の敵のみ。威力はAT+MaxSP依存"},
	}
}

// パワーアップの抽選
func SelectPowerUp(selType int) {
	this.SelType = selType
	IdList := []int{}
	ps := &player.GetThis().PState
	for i := 0; i < len(this.PDetail); i++ {
		//パワーアップ属性があっているパワーアップのみ取り出す
		if this.PDetail[i].SelType == selType {
			switch this.PDetail[i].Ptype {
			case P_RATE_EXP:
				if ps.RATE_EXP >= 100 {
					continue
				}
			case P_RATE_COIN:
				if ps.RATE_COIN >= 100 {
					continue
				}
			case P_RATE_FORCE:
				if ps.RATE_FORCE >= 100 {
					continue
				}
			}

			IdList = append(IdList, this.PDetail[i].Ptype)
		} else if this.PDetail[i].SelType == 5 && selType != 2 { //スキルは店では出ない...............
			//基本はOK
			isOk := true
			skID := this.PDetail[i].Ptype
			//スキルが全部埋まっている場合、
			if ps.Skill[0].SkillID != -1 &&
				ps.Skill[1].SkillID != -1 &&
				ps.Skill[2].SkillID != -1 &&
				ps.Skill[3].SkillID != -1 {
				//既得スキル以外、駄目

				if ps.Skill[0].SkillID != skID &&
					ps.Skill[1].SkillID != skID &&
					ps.Skill[2].SkillID != skID &&
					ps.Skill[3].SkillID != skID {
					isOk = false
				}
			}
			//上記条件を通過してきていても、アップグレード限界なら無理
			if isOk {
				for j := 0; j < 4; j++ {
					if ps.Skill[j].SkillID == skID {
						skillDeteil := skill.GetSkill(skID)
						if ps.Skill[j].CoolTimeMax <= skillDeteil.CoolTime {
							isOk = false
						}
					}
				}
			}
			if isOk {
				IdList = append(IdList, this.PDetail[i].Ptype)
			}
		}
	}
	//シャッフル
	rand.Shuffle(len(IdList), func(i, j int) { IdList[i], IdList[j] = IdList[j], IdList[i] })

	switch selType {
	case 1: //EXP
		//スキルは、最大で1つしか出ない。出る確率も1/2
		inskil := true
		/*if rand.Intn(2) == 0 {
			inskil = false
		}*/
		for i := 0; i < len(IdList); {
			//リスト中にスキルを発見
			if IdList[i] < 100 {
				//スキル採用可能なら
				if inskil {
					inskil = false //採用するが、移行は不可
				} else {
					// 要素をスライスから削除
					IdList = append(IdList[:i], IdList[i+1:]...)
					continue // 次の要素をチェック
				}
				i++
			} else {
				i++
			}
		}
		this.PLIST[0] = getPindex(IdList[0])
		this.PLIST[1] = getPindex(IdList[1])
		this.PLIST[2] = getPindex(IdList[2])

	case 2: //GOLD
		this.PLIST[0] = getPindex(P_HP_RECV1 + rand.Intn(3)) //最初の一つは回復
		this.PLIST[1] = getPindex(P_AT_T1 + rand.Intn(3))    //ATアップ
		//最後の一つはほりだしものが1/100で出てくる
		if rand.Intn(100) > 1 {
			this.PLIST[2] = getPindex(IdList[0])
		} else {
			for i := 0; i < len(IdList); i++ {
				switch rand.Intn(2) {
				case 0:
					this.PLIST[2] = getPindex(P_CAKE)
				case 1:
					this.PLIST[2] = getPindex(P_ALLSKILL)
				}
			}
		}

	case 3: //FORCE
		skillCount := 0
		//所持スキルを全滅させる
		for i := 0; i < len(IdList); {
			//リスト中にスキルを発見
			if IdList[i] < 100 {
				isDel := false
				for j := 0; j < len(this.player.PState.Skill); j++ {
					if this.player.PState.Skill[j].SkillID == IdList[i] {
						isDel = true
					}
				}
				if isDel {
					// 要素をスライスから削除
					IdList = append(IdList[:i], IdList[i+1:]...)
					continue // 次の要素をチェック

				}
			}
			i++
		}
		for i := 0; i < len(IdList); i++ {
			//リスト中にスキルを発見したら即採用
			if IdList[i] < 100 {
				this.PLIST[skillCount] = getPindex(IdList[i])
				skillCount++
				if skillCount >= 2 {
					break
				}
			}
		}
		//リストからスキルを全滅させる
		for i := 0; i < len(IdList); {
			//リスト中にスキルを発見
			if IdList[i] < 100 {

				// 要素をスライスから削除
				IdList = append(IdList[:i], IdList[i+1:]...)
				continue // 次の要素をチェック
			}
			i++
		}
		//スキルが有る限り全部スキル
		//スキルが不足したところにランダムなFORCEスキル
		for i := skillCount; i < 3; i++ {
			this.PLIST[i] = getPindex(IdList[i])
		}
	}

}

// パワーアップの効果発現
func PowerUp(pindex int) {
	ps := &player.GetThis().PState

	ptype := this.PDetail[pindex].Ptype

	switch ptype {
	case P_AT:
		ps.AT += 2
	case P_MAXHP:
		ps.MaxHP += 10
		ps.HP += 10

	case P_HP_RECV1, P_HP_RECV2, P_HP_RECV3:
		switch ptype {
		case P_HP_RECV1:
			ps.HP += ps.RES_HP
		case P_HP_RECV2:
			ps.HP += int(float64(ps.RES_HP) * 1.5)
		case P_HP_RECV3:
			ps.HP += ps.RES_HP * 2
		}
		if ps.HP > ps.MaxHP {
			ps.HP = ps.MaxHP
		}

	case P_AT_T1:
		ps.AT += 1
	case P_AT_T2:
		ps.AT += 2
	case P_AT_T3:
		ps.AT += 3
	case P_MAXSP_T:
		ps.MaxSP += 4
		ps.SP += 4

	case P_RESHPUP:
		ps.RES_HP += 10
	case P_RESSPUP:
		ps.RES_SP += 1

	case P_RATE_EXP:
		if ps.RATE_EXP < 40 {
			ps.RATE_EXP += 30
		} else {
			ps.RATE_EXP += 10
		}
	case P_RATE_COIN:
		if ps.RATE_COIN < 40 {
			ps.RATE_COIN += 30
		} else {
			ps.RATE_COIN += 10
		}
	case P_RATE_FORCE:
		if ps.RATE_FORCE < 40 {
			ps.RATE_FORCE += 30
		} else {
			ps.RATE_FORCE += 10
		}
	case P_COOLDOWN:
		for i := 0; i < 4; i++ {
			ps.Skill[i].CoolTime = 0
		}
	case P_SPIKE:
		if ps.SPIKE < ps.AT {
			ps.SPIKE = ps.AT
		}
		ps.SPIKE += 10

	case P_ALLSKILL:
		for i := 0; i < 4; i++ {
			if ps.Skill[i].SkillID != -1 {
				if ps.Skill[i].Active {
					if ps.Skill[i].ActiveTimeMax != 1 {
						ps.Skill[i].ActiveTime = ps.Skill[i].ActiveTimeMax
					}
				} else {
					ps.Skill[i].CoolTime = 0
				}

			}
		}
	case P_CAKE:
		ps.HP = ps.MaxHP
		ps.SP = ps.MaxSP

	}
	//スキルの場合、スキルをセット。
	if ptype < 100 {
		player.SkillSet(ptype)
	}

}

// パワーアップのインデックス取得
func getPindex(Ptype int) int {
	for i := 0; i < len(this.PDetail); i++ {
		if this.PDetail[i].Ptype == Ptype {
			return i
		}
	}
	return 0
}
