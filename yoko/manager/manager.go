package manager

import (
	"Yoko/bg"
	"Yoko/camera"
	"Yoko/common"
	"Yoko/effect"
	"Yoko/enemy"
	"Yoko/input"
	"Yoko/item"
	"Yoko/magic"
	"Yoko/mytext"
	"Yoko/myui"
	"Yoko/player"
	"Yoko/powerup"
	"Yoko/saveload"
	"Yoko/skill"
	"Yoko/sound"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	G_TITLE    = 0
	G_GAME     = 1
	G_SEL      = 2
	G_GAMEOVER = 3
	G_RETRY    = 4
	G_HAPPY    = 5

	S_ENEMY_HP = 0
	S_ENEMY_AT = 1
	S_ENEMY_DF = 2
)

type ManagerStruct struct {
	input   *input.InputStruct
	sound   *sound.SoundStruct
	bg      *bg.BgStruct
	player  *player.PlayerStruct
	enemy   *enemy.EnemyStruct
	effect  *effect.EffectStruct
	item    *item.ItemStruct
	powerup *powerup.PowerupStruct
	skill   *skill.SkillStruct
	magic   *magic.MagicStruct
	//ゲームパラメータ
	Gparam common.GameParam

	gMode int
	//汎用カウンタ
	Counter int
	//インターフェイス押下フラグ
	keyPress       bool //キープレス
	press          bool
	pressEnable    bool
	pressKeyEnable bool
}

var (
	this *ManagerStruct
)

// --===================================================================================
// 初期化
func Init(t *ManagerStruct) {
	this = t
	resetGame()
}

func SetClass(in *input.InputStruct, so *sound.SoundStruct, bg *bg.BgStruct, pl *player.PlayerStruct, en *enemy.EnemyStruct, ef *effect.EffectStruct, it *item.ItemStruct, po *powerup.PowerupStruct, sk *skill.SkillStruct, mg *magic.MagicStruct) {
	this.input = in
	this.sound = so
	this.bg = bg
	this.item = it
	this.player = pl
	this.enemy = en
	this.effect = ef
	this.powerup = po
	this.skill = sk
	this.magic = mg
}

// 計算
func Update() bool {
	pThis := player.GetThis()
	ethis := enemy.GetThis()
	eL := ethis.EList
	//プレイヤー座標
	pX := pThis.PState.X

	//キー押下------------------------------------
	if this.gMode != G_GAMEOVER {
		this.keyPress = input.GetKeyPress()
	} else {
		if this.pressKeyEnable {
			this.keyPress = input.GetKeyPress()
		} else {
			this.keyPress = false
			if !input.GetKeyPress() {
				this.pressKeyEnable = true
			}
		}
	}

	if !this.keyPress || this.gMode != G_GAMEOVER {
		//マウス、タッチパッド
		if this.pressEnable {
			this.press = input.GetPress()
		} else {
			this.press = false
			if !input.GetPress() {
				this.pressEnable = true
			}
		}
	}
	touchX, touchY, isTouch := input.GetTouchPos()

	//システム--------------------------------
	switch this.gMode {
	case G_TITLE: //タイトルは再押下されるまで
		if this.press || this.keyPress {
			this.gMode = G_GAME

			myui.SetUIMode(myui.UI_NOMAL)
		}
	case G_GAME:
		if myui.GetUIselect(touchX, touchY, this.press) == -1 {
			this.gMode = G_RETRY
			this.pressEnable = false
			myui.SetUIMode(myui.UI_RETRY)
			return false
		}

		if !this.Gparam.Happy && this.Gparam.Rank >= 100 {
			this.Gparam.Happy = true
			this.gMode = G_HAPPY
			myui.SetUIMode(myui.UI_HAPPY)
			sound.Play(sound.HAPPY)
			return false
		}
	case G_HAPPY:
		if myui.GetUIselect(touchX, touchY, this.press) == 1 {
			this.gMode = G_GAME
			myui.SetUIMode(myui.UI_NOMAL)
			return false
		}
		return false
	case G_RETRY:
		switch myui.GetUIselect(touchX, touchY, this.press) {
		case -1:
			this.gMode = G_GAME
			myui.SetUIMode(myui.UI_NOMAL)
			ForceGameOver()
		case 1:
			myui.SetUIMode(myui.UI_NOMAL)
			this.gMode = G_GAME
			this.pressEnable = false
		}

		return false
	}

	//魔法---------------------------------------------------------
	//即発動スキル感知
	mtype := []int{skill.S_ELECTRICTRIGGER, skill.S_THUNDERBOLT, skill.S_LIGHTNINGSWORD}
	for i := 0; i < len(mtype); i++ {

		if player.IsActive(mtype[i]) {
			sound.Play(sound.HIT2)

			magic.MagicGo(mtype[i])

			if player.IsActive(skill.S_ELECTRICTRIGGER) {
				dmg := pThis.PState.AT * 3
				MagicDamage(dmg, skill.S_ELECTRICTRIGGER)
			} else if player.IsActive(skill.S_THUNDERBOLT) {
				dmg := pThis.PState.MaxSP * 3
				MagicDamage(dmg, skill.S_THUNDERBOLT)
			} else if player.IsActive(skill.S_LIGHTNINGSWORD) {
				dmg := (pThis.PState.MaxSP + pThis.PState.AT) * 3
				MagicDamage(dmg, skill.S_LIGHTNINGSWORD)
			}
			player.CoolDownSkill()
		}
	}

	//魔法エフェクト発動中
	if magic.GetMagicActive() {
		return true
	}

	//パワーアップセレクト----------------------
	if this.gMode == G_GAME {
		if pThis.PState.EXP >= common.G_MAX {
			pThis.PState.EXP -= common.G_MAX
			powerup.StartPowerUp(1)
			this.pressEnable = false
			this.gMode = G_SEL
		} else if pThis.PState.GOLD >= common.G_MAX {
			pThis.PState.GOLD -= common.G_MAX
			powerup.StartPowerUp(2)
			this.pressEnable = false
			this.gMode = G_SEL
		} else if pThis.PState.FORCE >= common.G_MAX {
			pThis.PState.FORCE -= common.G_MAX
			powerup.StartPowerUp(3)
			this.pressEnable = false
			this.gMode = G_SEL
		}
	}

	if this.gMode == G_SEL {
		if powerup.Select(this.press) {
			player.SkillNoSelect() //プレイヤーがスキル押下中にパワーアップ発動したら、スキル選択却下
			this.gMode = G_GAME
			//プレイヤーパワーアップのたびに、セーブを実施
			saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)
		}

		return false
	}

	//敵の攻撃------------------------------------------------------------
	for i := 0; i < len(ethis.AttackBuffer); i++ {
		//敵の攻撃バッファを見る。値があったら
		if ethis.AttackBuffer[i] {
			//プレイヤー体力減少
			EnemyAT := this.Gparam.EN_BASE_AT + rand.Intn(3) - 1

			pThis.PState.SP -= EnemyAT
			if pThis.PState.SP < 0 {
				pThis.PState.HP += pThis.PState.SP
				pThis.PState.SP = 0
			}

			//バッファはクリア
			ethis.AttackBuffer[i] = false
			for j := 0; j < 10; j++ {
				effect.SetEffect(int(pX), -20, 0, effect.E_DAMAGE, 0)
			}
			//ヒット音
			sound.Play(sound.PHIT)
			//プレイヤー体力0
			if pThis.PState.HP <= 0 {
				pThis.PState.HP = 0

				//ゲームオーバー時、セーブを実施
				saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)

				if this.gMode != G_GAMEOVER {
					player.SetState(player.S_DOWN)
					StartGameOver()
					break
				}
			}
		}
	}
	//プレイヤーにカメラ追従
	if camera.X < pX {
		camera.X = pX
	}

	//プレイヤー生存していない
	if pThis.PState.HP <= 0 {
		KaisyuEnemy()
		UpdateGameOver()
		return true
	}

	//アイテムの取得----------------------------------------------------------
	it := item.GetThis()
	for i := 0; i < len(it.ItemList); i++ {
		ic := &it.ItemList[i]
		if !ic.IsUse || ic.MoveMode != item.M_KETTEIMACH {
			continue
		}
		switch ic.ItemType {
		case item.I_EXP:
			plus := int(float64(100) * bairitu(pThis.PState.RATE_EXP) / 100)
			if player.IsActive(skill.S_EXPx2) { //EXP倍スキル発動中
				plus *= 2
			}
			pThis.PState.EXP += plus
			pThis.PState.T_EXP += plus
			//記録
			if pThis.PState.H_T_EXP < pThis.PState.T_EXP {
				pThis.PState.H_T_EXP = pThis.PState.T_EXP
			}
			if player.IsActive(skill.S_EXP_RECOV_HP) { //EXPで回復発動中
				pThis.PState.HP += pThis.PState.RES_HP / 20
				if pThis.PState.HP > pThis.PState.MaxHP {
					pThis.PState.HP = pThis.PState.MaxHP
				}
			}
		case item.I_COIN:
			plus := int(float64(100) * bairitu(pThis.PState.RATE_COIN) / 100)
			if player.IsActive(skill.S_EXP_TO_GOLD) { //コイン倍スキル発動中
				plus *= 2
			}
			pThis.PState.GOLD += plus
			pThis.PState.T_GOLD += plus
			//記録
			if pThis.PState.H_T_GOLD < pThis.PState.T_GOLD {
				pThis.PState.H_T_GOLD = pThis.PState.T_GOLD
			}
		case item.I_FORCE:
			plus := int(float64(200) * bairitu(pThis.PState.RATE_FORCE) / 100)
			pThis.PState.FORCE += plus
			pThis.PState.T_FORCE += plus
			//記録
			if pThis.PState.H_T_FORCE < pThis.PState.T_FORCE {
				pThis.PState.H_T_FORCE = pThis.PState.T_FORCE
			}

			pThis.PState.SP += this.player.PState.RES_SP
			if pThis.PState.SP > pThis.PState.MaxSP {
				pThis.PState.SP = pThis.PState.MaxSP
			}

		}

		sound.Play(sound.CATCH)
		ic.IsUse = false
	}

	//パワーアップ発生時は、即アップグレードへ
	if pThis.PState.EXP >= common.G_MAX || pThis.PState.GOLD >= common.G_MAX || pThis.PState.FORCE >= common.G_MAX {
		return false
	}

	//プレイヤーの操作------------------------------------------------------
	isAttack := false
	//すべての敵を見て、打撃判定内にいたら
	for i := 0; i < len(eL); i++ {
		ec := &eL[i]
		if !ec.IsUse {
			continue
		}
		//攻撃開始
		if int(ec.X)-ec.Radius < pX+10 && enemy.IsHitOk(ec) {
			//プレイヤー攻撃開始
			isAttack = true
		}
	}

	//スキル選択----------------------------
	touchX2, touchY2, isTouch2 := input.GetTouchPos2()
	skillSelecttyu := false
	if isTouch2 {
		skillSelecttyu = player.SkillSel(touchX2, touchY2, isTouch2, true)
	} else {
		skillSelecttyu = player.SkillSel(touchX, touchY, isTouch, false)
	}
	if isAttack { //攻撃可能なら.........
		//攻撃開始
		player.SetState(player.S_ATTACK)
	} else {
		//プレイヤー移動..............
		if (this.press || this.keyPress) && !skillSelecttyu {
			player.SetState(player.S_MOVE_LEFT)
			this.Gparam.Rank += 0.004
			if this.Gparam.H_Rank < this.Gparam.Rank {
				this.Gparam.H_Rank = this.Gparam.Rank
			}
			//敵ステータス更新
			SetEnemyState()
		} else {
			player.SetState(player.S_STAND)
		}
	}

	//プレイヤー攻撃---------------------------------------------------------------
	//敵爆破及びアイテム放出

	//プレイヤーモーションが攻撃判定を持っていたら
	atType := player.GetAttackMothin()
	isAttackEffect := true
	fAttackCount := 0
	if atType > 0 {
		player.CoolDownSkill()

		goldCount := 0
		revengeEnemy := 0
		//全敵
		for i := 0; i < len(eL); i++ {
			ec := &eL[i]
			if !ec.IsUse {
				continue
			}
			//攻撃範囲内にいるやつ
			if int(ec.X)-ec.Radius < pX+30 && enemy.IsHitOk(ec) {
				if isAttackEffect {
					effect.SetEffect(int(ec.X)-int(ec.Radius/2), -16, 0, effect.E_HITEFFEC, 0)
					isAttackEffect = false
				}
				attackRate := 1
				if player.IsActive(skill.S_ATx2) {
					attackRate = 2
				}
				//攻撃ヒット
				atDmg := 0
				switch atType {
				case player.ATTACK_HIT_TYPE1, player.ATTACK_HIT_TYPE2:
					atDmg = pThis.PState.AT*attackRate - ec.DF
					sound.Play(sound.HIT1)
					if player.IsActive(skill.S_COIN_TO_SP) {
						if rand.Intn(2) == 0 {
							ItemSatter(ec.X, ec.Y, item.I_FORCE, 1, true)
						}
					} else {
						if goldCount < 2 {
							ItemSatter(ec.X, ec.Y, item.I_COIN, rand.Intn(2)+1, true)
						}
						goldCount++
					}
					if !player.IsActive(skill.S_COIN_TO_SP) {
						ec.HIT++
					} else {
						ec.HIT = 0
					}
				case player.ATTACK_HIT_TYPE3:
					atDmg = pThis.PState.AT*2*attackRate - ec.DF*2
					sound.Play(sound.HIT2)
					if fAttackCount < 1 {
						ItemSatter(ec.X, ec.Y, item.I_FORCE, rand.Intn(3)+1, true)
					}
					if !player.IsActive(skill.S_COIN_TO_SP) {
						ec.HIT += 2
					}
					fAttackCount++
				}

				//最低でも1は保証
				if atDmg < 1 {
					atDmg = 1
				}
				ec.HP -= atDmg

				if ec.HP <= 0 { //体力がなくなったらダメージモーション後退場
					ec.HP = 0
					enemy.SetState(ec, enemy.S_DMGLEVING1)
				} else if atType == player.ATTACK_HIT_TYPE3 { //フィニッシュブロー
					revengeEnemy++
					revengeMax := 1
					if rand.Intn(3) == 0 {
						revengeMax = 2
					}

					if revengeEnemy <= revengeMax {
						enemy.SetState(ec, enemy.S_DMGBROWREVNGE)
					} else {
						enemy.SetState(ec, enemy.S_DMGBROW)
					}
				} else { //通常攻撃
					enemy.SetState(ec, enemy.S_DMG)
				}
			}
		}
	}

	//一定位置毎に敵を配置する.............................................
	if camera.X > this.Gparam.CamLastX+this.Gparam.NextEnemy {
		//大きさは0-4まである。上ほど出現率は低い
		enemyP := []int{0, 10, 50, 80, 99}
		ran := rand.Intn(100) + 1
		eType := 0
		for i := len(enemyP) - 1; 0 <= i; i-- {
			if enemyP[i] <= ran {
				eType = i
				break
			}
		}
		enemy.SetEnemy(camera.X+common.SCREEN_WIDTH/2+64, eType, this.Gparam.EN_BASE_HP, this.Gparam.EN_BASE_DF)

		this.Gparam.CamLastX = camera.X
		//次の敵が配置される距離
		if player.IsActive(skill.S_ENEMY_FEEVER) {
			this.Gparam.NextEnemy = rand.Intn(20)
		} else {
			this.Gparam.NextEnemy = rand.Intn(20) + 40
		}
	}
	//..............................................
	KaisyuEnemy()
	return true
}

// 描画
func Draw(screen *ebiten.Image) {
	switch this.gMode {
	case G_TITLE:
		//		c := color.RGBA{0, 0, 255, 255}
		//		mytext.DrawG(screen, 50, "NECONECO RAID", 0, c)
	case G_GAMEOVER:
		y := -150 + this.Counter
		if y > 50 {
			y = 50
		}
		c := color.RGBA{255, 0, 0, 255}
		mytext.DrawG(screen, y, "GAME OVER", 0, c)
	}
}

// 回収待ちの敵を回収し、爆発に変える
func KaisyuEnemy() {
	ethis := enemy.GetThis()
	eL := ethis.EList

	//回収待ちの敵を爆発に変換する
	for i := 0; i < len(eL); i++ {
		ec := &eL[i]
		if !ec.IsUse {
			continue
		}
		//敵回収待ちなら
		if ec.EStete == enemy.S_KAISYUMACHI || ec.EStete == enemy.S_KAISYUMACHI2 {
			//その敵は消滅
			ec.IsUse = false
			if ec.SizeType >= 4 {
				camera.SetViblation(10)
				sound.Play(sound.EXPLOSION2)
			} else {
				sound.Play(sound.EXPLOSION1)
			}

			//EXPばらまき
			oct := 1
			switch ec.SizeType {
			case 0, 1:
				oct = 1
			case 2, 3:
				oct = 10
			case 4:
				oct = 50
			}

			//ゴールド変換中でないなら補正
			if !player.IsActive(skill.S_EXP_TO_GOLD) {
				octB := 0 //ec.HIT / 2
				if octB > oct {
					oct = octB
				}
			} else {
				oct /= 2
			}

			//スパイクやられか、コインをフォース変換中
			if ec.EStete == enemy.S_KAISYUMACHI2 || player.IsActive(skill.S_COIN_TO_SP) {
				if oct < 1 {
					oct = 1
				}
				if oct > 10 {
					oct = 10
				}
			}

			if player.IsActive(skill.S_EXP_TO_GOLD) {
				ItemSatter(ec.X, ec.Y, item.I_COIN, oct, false)

			} else {
				ItemSatter(ec.X, ec.Y, item.I_EXP, oct, false)
			}

			//代わりに爆発
			m := ec.Magnification
			if m < 0 {
				m = 1
			}
			effect.SetEffect(int(ec.X), int(ec.Y-float64(ec.Radius)), int(ec.Z), effect.E_EXPLOTION, m)
		}
	}
}

// アイテムばらまき
func ItemSatter(x, y float64, itemtype int, pieces int, rt bool) {
	for i := 0; i < pieces; i++ {
		vx := float64(-8)
		vy := float64(0)
		r := float64(0)
		if rt {
			r = rand.Float64()*60 + 30
		} else if pieces <= 5 {
			r = rand.Float64() * 180
		} else {
			r = float64((180 / (pieces - 1)) * (i))
		}

		vx, vy = rotateVector(vx, vy, r)
		item.SetItem(itemtype, x, y, vx, vy)
	}
}

// ベクトルを回転させる----------------------------
func rotateVector(x, y, angle float64) (float64, float64) {
	rad := angle * (math.Pi / 180.0)
	newX := x*math.Cos(rad) - y*math.Sin(rad)
	newY := x*math.Sin(rad) + y*math.Cos(rad)
	return newX, newY
}

// ゲームオーバー開始
func StartGameOver() {
	this.Counter = 0
	this.pressEnable = false
	this.gMode = G_GAMEOVER
}

// ゲームオーバー
func UpdateGameOver() {
	this.Counter++
	if this.Counter == 150 {
		sound.Play(sound.GAMEOVER)
	}
	if this.Counter > 150 {
		if (this.press && this.pressEnable) ||
			(this.keyPress && this.pressKeyEnable) {
			resetGame()
			//セーブを実施
			saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)
			this.pressEnable = false
			this.pressKeyEnable = false
		}
	}
}

// ロード直後ゲームオーバー
func ForceGameOver() {
	player.GetThis().PState.HP = 0
	player.SetState(player.S_DOWN)
	StartGameOver()
	saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)
}

// 魔法で敵にダメージ
func MagicDamage(damage int, magicID int) {
	//全敵
	eL := enemy.GetThis().EList
	for i := 0; i < len(eL); i++ {
		ec := &eL[i]
		if !ec.IsUse {
			continue
		}
		switch ec.EStete {
		case enemy.S_LEVING1, enemy.S_LEVING2, enemy.S_LEVING3, enemy.S_KAISYUMACHI, enemy.S_KAISYUMACHI2:
			continue
		}
		//空中敵のみ
		if magicID == skill.S_LIGHTNINGSWORD {
			if ec.EStete != enemy.S_JAMP && ec.EStete != enemy.S_TAME && ec.EStete != enemy.S_ATTACK {
				continue
			}

		}
		//攻撃ヒット
		dmgA := damage
		if magicID == skill.S_ELECTRICTRIGGER {
			dmgA -= ec.DF
		}
		if dmgA < 1 {
			dmgA = 1
		}
		ec.HP -= dmgA

		enemy.SetState(ec, enemy.S_DMG)

		if ec.HP <= 0 { //体力がなくなったらダメージモーション後退場
			ec.HP = 0
			enemy.SetState(ec, enemy.S_DMGLEVING3)
		} else {
			enemy.SetState(ec, enemy.S_DMGBROW)
		}
	}
}

// ----------------------------------------------------------------------------------------------
// ゲームリセット
func resetGame() {
	this.Gparam.Rank = 1
	this.Gparam.CamLastX = 0
	this.Gparam.DebugRecRank = 0
	this.Gparam.Happy = false
	this.gMode = G_TITLE
	this.pressEnable = false
	this.pressKeyEnable = false
	this.Counter = 20

	camera.X = 0
	player.InitState(50, 5, 10)
	enemy.Reset()
	item.Reset()
	magic.Reset()
	myui.SetUIMode(myui.UI_TITLE)

	if this.powerup != nil {
		this.Gparam.FirstPowerUp = true
	}
	SetEnemyState()
}

// レートを内部倍率に変換
func bairitu(rate int) float64 {
	return float64(rate*60/100) + 60
}

// ランクに対する敵のステータス
func SetEnemyState() {
	rankBase := this.Gparam.Rank

	//ある程度を超えたら内部倍率倍
	rankUpBorder1 := float64(90)
	if rankUpBorder1 < rankBase {
		rankBase += (rankBase - rankUpBorder1)
	}
	rankUpBorder2 := float64(rankUpBorder1 + 20*2)
	if rankUpBorder2 < rankBase {
		rankBase += (rankBase - rankUpBorder2)
	}

	this.Gparam.EN_BASE_DF = int((rankBase) * ((rankBase) + 10) / 120)
	//--------------------------------------------
	//基本
	this.Gparam.EN_BASE_AT = int((rankBase)*((rankBase)+10)/300) + 2

	if 20 < rankBase {
		x := int((rankBase-20)*(rankBase-20)/200) + 3
		this.Gparam.EN_BASE_AT += x
	}

	//--------------------------------------------
	this.Gparam.EN_BASE_HP = int(rankBase)*10 + 20
	if this.Gparam.EN_BASE_HP > 100 {
		this.Gparam.EN_BASE_HP = 100
	}

	//--------------------------------------------

}
