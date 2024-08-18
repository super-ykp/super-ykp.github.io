package player

import (
	"Yoko/camera"
	"Yoko/common"
	"Yoko/mytext"
	"Yoko/skill"
	"Yoko/sound"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	S_STAND     = 0
	S_MOVE_LEFT = 1
	S_ATTACK    = 2
	S_DOWN      = 3
)

const (
	pNG_OFFSET_X  = 0
	pNG_OFFSET_Y  = 32
	sPNG_OFFSET_X = 0
	sPNG_OFFSET_Y = 192

	a_STAND   = 0
	a_WALK    = 1
	a_ATTACK1 = 2
	a_DOWN    = 3
	a_MAX     = 4 //アニメの最大数

	ATTACK_HIT_TYPE1 = 1
	ATTACK_HIT_TYPE2 = 2
	ATTACK_HIT_TYPE3 = 3

	SOUND1 = 10
	SOUND2 = 20

	SKILL_Y = -float32(common.CAM_Y_OFFSET) + 33
)

type SkillStruct struct {
	SkillID       int
	Level         int
	CoolTimeMax   int
	CoolTime      int
	ActiveTimeMax int
	ActiveTime    int
	Active        bool
}

type PstateStruct struct {
	X int //X座標。ワールド位置

	MaxHP int //最大HP
	HP    int
	MaxSP int //最大SP
	SP    int
	AT    int //攻撃力
	SPIKE int //スパイク

	RES_HP     int //HP回復力
	RES_SP     int //SP回復力
	RATE_EXP   int //経験値ゲージ倍率
	RATE_COIN  int //コインゲージ倍率
	RATE_FORCE int //フォースゲージ倍率

	EXP   int
	GOLD  int
	FORCE int
	Skill [4]SkillStruct

	T_EXP   int
	T_GOLD  int
	T_FORCE int

	H_T_EXP   int
	H_T_GOLD  int
	H_T_FORCE int
}

type PlayerStruct struct {
	playerImage  []*ebiten.Image
	skillImage   []*ebiten.Image
	pMode        uint //状態
	animMode     uint //アニメモード。歩くとか止まっているとか
	animIndex    uint //現アニメ中の何番目にいるか
	animCounter  uint //次PNGまでの時間カウンタ
	dispPngIndex uint //表示PNGインデックス

	PState   PstateStruct
	SelSkill int

	G_counter uint
	flg1      bool //汎用フラグ
}

var (
	this  *PlayerStruct
	animC [][][]uint
)

// --===================================================================================
// 初期化
func Init(t *PlayerStruct, p *ebiten.Image) {
	this = t
	offs := [][]int{
		{0, 0, 32, 32},    //0 A_STAND
		{32, 0, 32, 32},   //1
		{64, 0, 32, 32},   //2 A_WALK
		{96, 0, 32, 32},   //3
		{128, 0, 32, 32},  //4
		{0, 32, 32, 32},   //5 ATTACK
		{32, 32, 32, 32},  //6
		{64, 32, 32, 32},  //7
		{96, 32, 32, 32},  //8
		{128, 32, 32, 32}, //9
		{160, 32, 32, 32}, //10
		{192, 32, 48, 32}, //11
		{0, 64, 32, 32},   //12
		{32, 64, 32, 32},  //13
		{64, 64, 32, 32},  //14
		{96, 64, 32, 32},  //15
		{128, 64, 48, 32}, //16
		{160, 0, 32, 32},  //17
		{192, 0, 32, 32},  //18
		{224, 0, 48, 32},  //19
	}
	for i := 0; i < len(offs); i++ {
		offX := pNG_OFFSET_X + offs[i][0]
		offy := pNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]

		imageRect := image.Rect(offX, offy, offW, offH)
		rect := p.SubImage(imageRect).(*ebiten.Image)
		this.playerImage = append(this.playerImage, rect)
	}

	animC = make([][][]uint, a_MAX)
	animC[a_STAND] = [][]uint{{0, 20}, {1, 20}}
	animC[a_WALK] = [][]uint{{2, 10}, {4, 10}, {3, 10}, {4, 10}}
	animC[a_ATTACK1] = [][]uint{
		{5, 1}, {6, ATTACK_HIT_TYPE1*100 + 10},
		{5, 1}, {6, ATTACK_HIT_TYPE2*100 + 10},
		{5, 1}, {6, ATTACK_HIT_TYPE2*100 + 10},
		{7, SOUND1*100 + 4}, {8, 3}, {9, 3}, {10, 3}, {11, ATTACK_HIT_TYPE2*100 + 10},
		{12, 2}, {13, 2}, {14, 2}, {15, SOUND1*100 + 10}, {16, ATTACK_HIT_TYPE3*100 + 20},
	}
	animC[a_DOWN] = [][]uint{{17, 3}, {18, 3}, {17, 3}, {18, 3}, {17, 5}, {18, 5}, {17, 10}, {18, 10}, {17, 20}, {18, 60}, {19, SOUND2*100 + 0}}

	//スキル
	for i := 0; i < 20; i++ {
		offx := sPNG_OFFSET_X + (i%20)*16
		offy := sPNG_OFFSET_Y + i/20*16
		imageRect := image.Rect(offx, offy, offx+16, offy+16)
		rect := p.SubImage(imageRect).(*ebiten.Image)
		this.skillImage = append(this.skillImage, rect)
	}
}

// 計算
func Update() {
	switch this.pMode {
	case S_MOVE_LEFT: //歩き中
		this.PState.X++
	}

	//アニメーション制御---------------------------------------------------------------------
	//アニメーションフレーム。下位2桁ががフレーム
	countMax := animC[this.animMode][this.animIndex][1] % 100
	this.animCounter++
	//アニメーション表示フレームを超えたら
	if this.animCounter > countMax && countMax != 0 {
		//次のアニメパターンへ
		this.animCounter = 0
		this.animIndex++

		//アニメーションが最初まで行ったら0に
		if this.animIndex == uint(len(animC[this.animMode])) {
			this.animIndex = 0
		} else {
			act := animC[this.animMode][this.animIndex][1] / 100
			if act == SOUND1 {
				sound.Play(sound.FURI)
			}
			if act == SOUND2 {
				sound.Play(sound.PDOWN)
			}
		}
	}
	//元アニメのPng番号を指定
	this.dispPngIndex = animC[this.animMode][this.animIndex][0]
	//--------------------------------------------------------------------------------------
}

// 描画
func Draw(screen *ebiten.Image) {
	drawImageOption := ebiten.DrawImageOptions{}
	ofX := (this.PState.X - 32/2) - camera.CamOffsetX
	drawImageOption.GeoM.Translate(float64(ofX), -32-float64(camera.CamOffsetY))
	screen.DrawImage(this.playerImage[this.dispPngIndex], &drawImageOption)
}

// プレイヤーへの動作指定
func SetState(state uint) {
	//すでにそのモーションでない場合のみ
	if this.pMode != state {
		this.pMode = state
		switch state {
		case S_STAND: //停止
			this.animMode = a_STAND
		case S_MOVE_LEFT: //右に歩く
			this.animMode = a_WALK
		case S_ATTACK:
			this.animMode = a_ATTACK1
		case S_DOWN:
			this.animMode = a_DOWN
			sound.Play(sound.PGURU)
		}
		//アニメはリセット
		this.animCounter = 0
		this.animIndex = 0
	}
}

// プレイヤー状態取得
func GetThis() *PlayerStruct {
	return this
}

func GetAttackMothin() uint {
	if this.animCounter == 1 {
		aDat := animC[this.animMode][this.animIndex][1]
		atk := (aDat / 100) % 10
		return atk
	}
	return 0
}

func FeedBackAttack() {
	this.animCounter = 0
}

// ダウンアニメ中。このアニメだけは最前面に出すため
func IsDown() bool {
	return a_DOWN == this.animMode
}

// ==============================================================================================================
//スキル関係。スキルはプレイヤー付属

func Draw2(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, SKILL_Y, common.SCREEN_WIDTH, common.SCREEN_WIDTH/4, color.RGBA{128, 0, 128, 255}, false)
	for i := 0; i < 4; i++ {
		if this.PState.Skill[i].SkillID != -1 {
			drawPowerUp(screen, i)
		}
	}
}

// スキルパネル1つ
func drawPowerUp(screen *ebiten.Image, index int) {
	c := color.RGBA{64, 64, 64, 255}
	if this.PState.Skill[index].Active {
		this.G_counter++
		cc := uint8((math.Cos(float64(this.G_counter/10)) + 1) * 128 / 2)
		c = color.RGBA{cc, cc, cc, 255}

	} else if this.PState.Skill[index].CoolTime > 0 {
		c = color.RGBA{0, 0, 0, 255}
	} else if this.SelSkill == index {
		c = color.RGBA{128, 128, 128, 255}
	}
	x := (common.SCREEN_WIDTH / 4) * float32(index)

	vector.DrawFilledRect(screen, x+1, SKILL_Y+1, common.SCREEN_WIDTH/4-2, common.SCREEN_WIDTH/4-2, c, false)

	sk := this.PState.Skill[index].SkillID
	drawImageOption := ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(-8, -8)
	drawImageOption.GeoM.Scale(2, 2)
	drawImageOption.GeoM.Translate(float64(x)+common.SCREEN_WIDTH/8, float64(SKILL_Y)+common.SCREEN_WIDTH/8)
	screen.DrawImage(this.skillImage[sk], &drawImageOption)
	p := message.NewPrinter(language.Japanese)
	if this.PState.Skill[index].CoolTime == 0 {
		if this.PState.Skill[index].Active {
			mytext.DrawShadowed(screen, int(x), int(SKILL_Y), p.Sprintf("%d", this.PState.Skill[index].ActiveTime))

		} else {
			mytext.DrawShadowed(screen, int(x), int(SKILL_Y), p.Sprintf("%d OK!", this.PState.Skill[index].ActiveTime))
		}

	} else {
		mytext.DrawShadowed(screen, int(x), int(SKILL_Y)+64-16, p.Sprintf("WAIT%3d", this.PState.Skill[index].CoolTime))
	}
}

// スキルの選択。trueを返している間はプレイヤーは動けない
func SkillSel(tx, ty int, b bool, touch2 bool) bool {
	//.flg1は、前回押下中true
	if this.SelSkill != -1 && (!b || touch2) && this.flg1 && !this.PState.Skill[this.SelSkill].Active {
		SkillStart(this.SelSkill)
		this.flg1 = b
		return false
	}
	this.flg1 = b

	//----------------------------------
	this.SelSkill = -1
	//選択箇所
	for i := 0; i < 4; i++ {
		x := (common.SCREEN_WIDTH / 4) * float32(i)
		y := SKILL_Y

		if b && x < float32(tx) && tx < int(x+common.SCREEN_WIDTH/4) && y < float32(ty) && float32(ty) < y+common.SCREEN_WIDTH/4 {
			this.SelSkill = i
			return true
		}
	}
	return false
}

// 選択中にパワーアップが発生した場合、強制的に選択状態を解除する
func SkillNoSelect() {
	this.SelSkill = -1
	this.flg1 = false
}

// スキルセット(獲得またはアップグレード)
func SkillSet(skillID int) {
	sDetail := skill.GetSkill(skillID)

	for i := 0; i < len(this.PState.Skill); i++ {

		if this.PState.Skill[i].SkillID == -1 {
			//新規箇所に配置
			this.PState.Skill[i].SkillID = skillID
			this.PState.Skill[i].Level = 1
			this.PState.Skill[i].CoolTime = 0
			this.PState.Skill[i].CoolTimeMax = sDetail.CoolTime
			this.PState.Skill[i].ActiveTimeMax = sDetail.ActiveTime
			this.PState.Skill[i].ActiveTime = sDetail.ActiveTime
			break
		}

	}
}

// スキルクールダウン
func CoolDownSkill() {
	for i := 0; i < len(this.PState.Skill); i++ {
		//非アクティブならクール値減少
		if !this.PState.Skill[i].Active {
			this.PState.Skill[i].CoolTime--
			if this.PState.Skill[i].CoolTime < 0 {
				this.PState.Skill[i].CoolTime = 0
				this.PState.Skill[i].ActiveTime = this.PState.Skill[i].ActiveTimeMax
			}
		} else { //利用中なら利用限界までカウントダウン
			this.PState.Skill[i].ActiveTime--
			if this.PState.Skill[i].ActiveTime <= 0 {
				this.PState.Skill[i].Active = false
				this.PState.Skill[i].CoolTime = this.PState.Skill[i].CoolTimeMax
			}
		}
	}
}

// スキルの発動
func SkillStart(index int) {
	if this.PState.Skill[index].CoolTime == 0 {
		this.PState.Skill[index].ActiveTime = this.PState.Skill[index].ActiveTimeMax
		this.PState.Skill[index].Active = true
		this.SelSkill = -1
	}
}

// 有効化しているスキルの判定
func IsActive(SkillID int) bool {
	for i := 0; i < len(this.PState.Skill); i++ {
		if this.PState.Skill[i].Active && this.PState.Skill[i].SkillID == SkillID {
			return true
		}
	}
	return false
}

// ステータスの初期化
func InitState(maxHP, maxSP, AT int) {
	this.PState.X = common.SCREEN_WIDTH
	this.animMode = a_STAND
	this.animIndex = 0
	//-------------------------
	this.PState.MaxHP = maxHP
	this.PState.HP = this.PState.MaxHP
	this.PState.MaxSP = maxSP
	this.PState.SP = this.PState.MaxSP
	this.PState.AT = AT

	this.PState.RES_HP = maxHP
	this.PState.RES_SP = 1
	this.PState.SPIKE = 0

	this.PState.RATE_EXP = 0
	this.PState.RATE_COIN = 0
	this.PState.RATE_FORCE = 0

	this.PState.EXP = 0
	this.PState.GOLD = 0
	this.PState.FORCE = 0

	this.PState.T_EXP = 0
	this.PState.T_GOLD = 0
	this.PState.T_FORCE = 0

	this.PState.Skill = [4]SkillStruct{}
	for i := 0; i < len(this.PState.Skill); i++ {
		this.PState.Skill[i].SkillID = -1
	}
	/*
		SkillSet(skill.S_BEEM)      //TEST■■■■■■
		SkillSet(skill.S_CANCELLER) //TEST■■■■■■
		SkillSet(skill.S_LIGHTNING) //TEST■■■■■■
	*/
}
