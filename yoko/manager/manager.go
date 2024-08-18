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
	G_INFO     = 0 //ランク100を目指しましょう！
	G_GAME     = 1 //ゲーム中
	G_SEL      = 2 //パワーアップセレクト中
	G_GAMEOVER = 3 //ゲームオーバー演出中
	G_MENU     = 4 //リトライまち中
	G_HAPPY    = 5 //ランク100到達
)

type ManagerStruct struct {
	//外部から引き取ったクラス参照
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

	//ゲームモード
	gMode int
	//汎用カウンタ
	Counter int
	//インターフェイス押下フラグ
	IskeyPress        bool //キープレス
	IsTouch           bool //タッチプレス
	EnabledIskeyPress bool //キープレスを有効に扱うか
	EnabledIsTouch    bool //タッチプレスを有効に扱うか
}

// 自インスタンス
var (
	this *ManagerStruct
)

// --===================================================================================
// 初期化
func Init(t *ManagerStruct) {
	this = t //go使いはこんな事しちゃいけないよ！
	//ゲームリセット
	resetGame()
}

// 外部クラスの取り込み
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
	//プレイヤーと敵のインスタンスを引き出す
	pThis := player.GetThis()
	ethis := enemy.GetThis()
	//敵リスト
	eL := ethis.EList
	//プレイヤー座標
	pX := pThis.PState.X

	//キー押下---------------------------------------------------------------------------
	//基本、押しっぱなしではなく押し直しに意味があるため、一度離さないと無意味、という処理を積んでおく

	//キー押下....
	if this.gMode != G_GAMEOVER { //ゲームオーバー中でないなら
		//キープレスを更新する
		this.IskeyPress = input.GetKeyPress()
	} else {
		//キープレス取り扱い有効状態
		if this.EnabledIskeyPress {
			//今の状態を取り込む
			this.IskeyPress = input.GetKeyPress()
		} else { //無効状態なら
			//キーは押されていないことにする
			this.IskeyPress = false
			//離したとき
			if !input.GetKeyPress() {
				//取り扱い有効に
				this.EnabledIskeyPress = true
			}
		}
	}

	//画面タッチ.....
	if !this.IskeyPress || this.gMode != G_GAMEOVER {
		//マウス、タッチパッド
		if this.EnabledIsTouch {
			this.IsTouch = input.GetPress()
		} else {
			this.IsTouch = false
			if !input.GetPress() {
				this.EnabledIsTouch = true
			}
		}
	}

	//生情報も持っておく
	touchX1, touchY1, isTouch1 := input.GetTouchPos()

	//システム--------------------------------
	switch this.gMode {
	case G_INFO: //説明は再押下されるまで
		//キーまたはタッチで
		if this.IsTouch || this.IskeyPress {
			//ゲームモードに
			this.gMode = G_GAME
			//UIをゲームモードに
			myui.SetUIMode(myui.UI_NOMAL)
		}
	case G_GAME: //ゲームモード。ほぼメイン
		//UIからメニューを開く指示が来たら
		if myui.GetUIselect(touchX1, touchY1, this.IsTouch) == -1 {
			//メニューを開く
			this.gMode = G_MENU
			this.EnabledIsTouch = false
			//メニュー画面へ
			myui.SetUIMode(myui.UI_MENU)
			//この先処理しない
			return false
		}

		//ランク100到達！
		if !this.Gparam.Happy && this.Gparam.Rank >= 100 {
			//到達フラグ
			this.Gparam.Happy = true
			//ハッピーモード
			this.gMode = G_HAPPY
			//ハッピーUI
			myui.SetUIMode(myui.UI_HAPPY)
			//キュピーン
			sound.Play(sound.HAPPY)
			return false
		}
	case G_HAPPY: //おめでとう
		//続くが押されるまで
		if myui.GetUIselect(touchX1, touchY1, this.IsTouch) == 1 {
			this.gMode = G_GAME
			myui.SetUIMode(myui.UI_NOMAL)
			return false
		}
		return false
	case G_MENU: //メニュー
		switch myui.GetUIselect(touchX1, touchY1, this.IsTouch) {
		case -1: //自爆
			this.gMode = G_GAME
			myui.SetUIMode(myui.UI_NOMAL)
			ForceGameOver()
		case 1: //通常に戻り
			myui.SetUIMode(myui.UI_NOMAL)
			this.gMode = G_GAME
			this.EnabledIsTouch = false
		}

		return false
	}

	//魔法---------------------------------------------------------
	//即発動スキル感知(前回フレームでアクティブ化される)
	mtype := []int{skill.S_ELECTRICTRIGGER, skill.S_THUNDERBOLT, skill.S_LIGHTNINGSWORD}
	//上記どの魔法か
	for i := 0; i < len(mtype); i++ {
		//一致していたら
		if player.IsActive(mtype[i]) {
			sound.Play(sound.HIT2)
			//エフェクト
			magic.MagicGo(mtype[i])

			if player.IsActive(skill.S_ELECTRICTRIGGER) { //エレクトリッガー
				dmg := pThis.PState.AT * 3
				MagicDamage(dmg, skill.S_ELECTRICTRIGGER)
			} else if player.IsActive(skill.S_THUNDERBOLT) { //サンダーボルト
				dmg := pThis.PState.MaxSP * 3
				MagicDamage(dmg, skill.S_THUNDERBOLT)
			} else if player.IsActive(skill.S_LIGHTNINGSWORD) { //マジカル艦砲射撃
				dmg := (pThis.PState.MaxSP + pThis.PState.AT) * 3
				MagicDamage(dmg, skill.S_LIGHTNINGSWORD)
			}
			//即クールダウン
			player.CoolDownSkill()
		}
	}

	//魔法エフェクト発動中は敵味方ストップ
	if magic.GetMagicActive() {
		//ここから下は処理打ち切り
		return true
	}

	//パワーアップセレクト----------------------
	//ゲームモードである事
	if this.gMode == G_GAME {
		//EXPが規定値を超えたら
		if pThis.PState.EXP >= common.G_MAX {
			//既定値だけマイナスし
			pThis.PState.EXP -= common.G_MAX
			//パワーアップセレクトを開始する
			powerup.StartPowerUp(1)
			//タッチ状態は無効化する
			this.EnabledIsTouch = false
			//セレクトモード
			this.gMode = G_SEL
		} else if pThis.PState.GOLD >= common.G_MAX {
			pThis.PState.GOLD -= common.G_MAX
			powerup.StartPowerUp(2)
			this.EnabledIsTouch = false
			this.gMode = G_SEL
		} else if pThis.PState.FORCE >= common.G_MAX {
			pThis.PState.FORCE -= common.G_MAX
			powerup.StartPowerUp(3)
			this.EnabledIsTouch = false
			this.gMode = G_SEL
		}
	}

	//パワーアップのセレクト中
	if this.gMode == G_SEL {
		//セレクト実施されたら
		if powerup.Select(this.IsTouch) {
			player.SkillNoSelect() //プレイヤーがスキル押下中にパワーアップ発動したら、スキル選択却下
			//ゲームモードに戻る
			this.gMode = G_GAME
			//プレイヤーパワーアップのたびに、セーブを実施
			saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)
		}

		return false
	}

	//敵の攻撃------------------------------------------------------------
	//敵の攻撃は攻撃バッファに接触したというフラグだけが上がってくる
	//大きさ情報はなく、個体別のダメージ設定も存在しないため基本一定

	for i := 0; i < len(ethis.AttackBuffer); i++ {
		//敵の攻撃バッファを見る。値があったら
		if ethis.AttackBuffer[i] {
			//処理した敵攻撃バッファクリア
			ethis.AttackBuffer[i] = false

			//威力算出。ベース値にランダム値加算。
			EnemyAT := this.Gparam.EN_BASE_AT + rand.Intn(3) - 1
			//威力分SPを減らす
			pThis.PState.SP -= EnemyAT
			//SPがマイナスになったらマイナス分体力減少。SPは0に戻る
			if pThis.PState.SP < 0 {
				pThis.PState.HP += pThis.PState.SP
				pThis.PState.SP = 0
			}
			//ダメージエフェクト「＃」が飛ぶ
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
					//プレイヤーモーションをダウンに設定
					player.SetState(player.S_DOWN)
					//ゲームオーバー開始
					StartGameOver()
					break
				}
			}
		}
	}
	//プレイヤーにカメラ追従-------------------------------------------------------------------
	if camera.X < pX {
		camera.X = pX
	}

	//プレイヤー生存していない
	if pThis.PState.HP <= 0 {
		//KaisyuEnemyはこの関数の最後にあるが、ここで止めてしまうとそこに行けないのでここでやってしまう
		KaisyuEnemy()
		//ゲームオーバーを進行させる
		UpdateGameOver()
		//この先の処理は打ち切り
		return true
	}

	//アイテムの取得----------------------------------------------------------
	it := item.GetThis()
	//アイテムバッファ一周
	for i := 0; i < len(it.ItemList); i++ {
		ic := &it.ItemList[i]
		//要素未使用か、回収まちでなければ、次
		if !ic.IsUse || ic.MoveMode != item.M_KETTEIMACH {
			continue
		}
		//回収。アイテム種類による
		switch ic.ItemType {
		case item.I_EXP: //EXP
			//倍率補正
			plus := int(float64(100) * bairitu(pThis.PState.RATE_EXP) / 100)
			if player.IsActive(skill.S_EXPx2) { //EXP倍スキル発動中
				plus *= 2 //倍に
			}
			//取得数加算
			pThis.PState.EXP += plus
			//トータルにも加算
			pThis.PState.T_EXP += plus
			//最大値記録を更新したら記録
			if pThis.PState.H_T_EXP < pThis.PState.T_EXP {
				pThis.PState.H_T_EXP = pThis.PState.T_EXP
			}
			//EXPで回復発動中
			if player.IsActive(skill.S_EXP_RECOV_HP) {
				//RES_HPの1/20回復
				pThis.PState.HP += pThis.PState.RES_HP / 20
				//最大値は超えない
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

	//上記回収をおえてパワーアップ発生時は、即パワーアップへ
	if pThis.PState.EXP >= common.G_MAX || pThis.PState.GOLD >= common.G_MAX || pThis.PState.FORCE >= common.G_MAX {
		return false
	}

	//プレイヤーの行動------------------------------------------------------
	isAttack := false
	//すべての敵を見て、打撃判定内に一体でもいたら
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
	//スキルの選択は全身のためのタッチと同時に行える（ダブルタッチ）必要がある
	touchX2, touchY2, isTouch2 := input.GetTouchPos2()

	//スキルセレクト中。スキルの発動は押したときではなく離したとき。枠外にずらしてから離したら当然無効化させる
	skillSelecttyu := false
	if isTouch2 {
		skillSelecttyu = player.SkillSel(touchX2, touchY2, isTouch2, true)
	} else {
		skillSelecttyu = player.SkillSel(touchX1, touchY1, isTouch1, false)
	}

	//攻撃発動していたら---------------------------------------------------------------------
	if isAttack {
		//プレイヤー攻撃モーション
		player.SetState(player.S_ATTACK)
	} else {
		//そうでないなら、キー入力があり、スキル選択でなければ
		if (this.IsTouch || this.IskeyPress) && !skillSelecttyu {
			//プレイヤー歩くモーション
			player.SetState(player.S_MOVE_LEFT)
			//ランク上昇
			this.Gparam.Rank += 0.004
			//最大ランク更新
			if this.Gparam.H_Rank < this.Gparam.Rank {
				this.Gparam.H_Rank = this.Gparam.Rank
			}
			//敵ステータス更新
			SetEnemyState()
		} else {
			//プレイヤー立ち状態
			player.SetState(player.S_STAND)
		}
	}

	//プレイヤー攻撃のヒット---------------------------------------------------------------

	//プレイヤーモーション情報取得
	atType := player.GetAttackMothin()
	//打撃エフェクト。何体いようが1つしか出ないようにする
	isAttackEffect := true

	//ゴールドカウント。通常打撃時、巻き込みで増えるコインには上限がある
	goldCount := 0
	//反撃に出る敵の数。いくら巻き込んでも基本1、まれに2
	revengeCounter := 0
	revengeMax := 1
	//通常1体、1/4で2体が反撃
	if rand.Intn(4) == 0 {
		revengeMax = 2
	}

	//プレイヤーモーションに攻撃判定が付いていたら
	if atType > 0 {
		//スキルは打撃してクールダウンする
		player.CoolDownSkill()

		//全敵バッファループ
		for i := 0; i < len(eL); i++ {
			ec := &eL[i]
			if !ec.IsUse {
				continue
			}
			//攻撃範囲内にいるやつ
			if int(ec.X)-ec.Radius < pX+30 && enemy.IsHitOk(ec) {
				//打撃エフェクト。重なっていても一回だけ
				if isAttackEffect {
					effect.SetEffect(int(ec.X)-int(ec.Radius/2), -16, 0, effect.E_HITEFFEC, 0)
					isAttackEffect = false
				}
				//攻撃倍率
				attackRate := 1
				//攻撃力2倍スキル使用時
				if player.IsActive(skill.S_ATx2) {
					attackRate = 2
				}
				//攻撃ヒット
				atDmg := 0
				switch atType {
				case player.ATTACK_HIT_TYPE1, player.ATTACK_HIT_TYPE2: //通常打撃。TYPE1,2は結局同じに..........
					//与えるダメージの計算式。王道のアルテリオス式
					atDmg = pThis.PState.AT*attackRate - ec.DF
					//打撃音
					sound.Play(sound.HIT1)
					//コインをフォース変換
					if player.IsActive(skill.S_COIN_TO_FORCE) {
						//通常より減らす
						if rand.Intn(2) == 0 {
							//フォース放出
							ItemSatter(ec.X, ec.Y, item.I_FORCE, 1, true)
						}
					} else {
						//重なった敵は2体まで
						if goldCount < 2 {
							//コイン放出
							ItemSatter(ec.X, ec.Y, item.I_COIN, rand.Intn(2)+1, true)
						}
						goldCount++
					}

				case player.ATTACK_HIT_TYPE3: //ドロップキック..........................................
					//ダメージは2回ぶん
					atDmg = pThis.PState.AT*2*attackRate - ec.DF*2
					//カーン
					sound.Play(sound.HIT2)
					//フォースを放出
					ItemSatter(ec.X, ec.Y, item.I_FORCE, rand.Intn(3)+1, true)
				}
				//敵に与えるダメージ............................
				//最低でも1は保証
				if atDmg < 1 {
					atDmg = 1
				}
				//減算
				ec.HP -= atDmg

				if ec.HP <= 0 { //体力がなくなったらダメージモーション後退場
					ec.HP = 0
					enemy.SetState(ec, enemy.S_DMGLEVING1)
				} else if atType == player.ATTACK_HIT_TYPE3 { //まだ生きていて食らったのがドロップキック
					//反撃可能範囲内なら
					if revengeCounter < revengeMax {
						//反撃モーションへ
						enemy.SetState(ec, enemy.S_DMGBROWREVNGE)
						revengeCounter++

					} else {
						//倒れモーションへ
						enemy.SetState(ec, enemy.S_DMGBROW)
					}
				} else {
					//打撃を食らっている
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
		//抽選
		for i := len(enemyP) - 1; 0 <= i; i-- {
			if enemyP[i] <= ran {
				eType = i
				break
			}
		}
		//その敵を配置
		enemy.SetEnemy(camera.X+common.SCREEN_WIDTH/2+64, eType, this.Gparam.EN_BASE_HP, this.Gparam.EN_BASE_DF)

		//現在のカメラ位置を保存
		this.Gparam.CamLastX = camera.X
		//次の敵が配置される距離を計算しておく
		if player.IsActive(skill.S_ENEMY_FEEVER) {
			this.Gparam.NextEnemy = rand.Intn(20)
		} else {
			this.Gparam.NextEnemy = rand.Intn(20) + 40
		}
	}
	//..............................................
	//爆発待ちの敵の回収
	KaisyuEnemy()
	return true
}

// 描画
func Draw(screen *ebiten.Image) {
	switch this.gMode {
	case G_INFO:
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
			if ec.SizeType >= 4 { //一番大きい敵は派手めなエフェクト
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

			//ゴールド変換中なら補正
			if player.IsActive(skill.S_EXP_TO_GOLD) {
				oct /= 2
			}

			//スパイクやられか、コインをフォース変換中
			if ec.EStete == enemy.S_KAISYUMACHI2 || player.IsActive(skill.S_COIN_TO_FORCE) {
				//1以下にはしない
				if oct < 1 {
					oct = 1
				}
				//10以上にもしない
				if oct > 10 {
					oct = 10
				}
			}

			if player.IsActive(skill.S_EXP_TO_GOLD) {
				//コインを出す
				ItemSatter(ec.X, ec.Y, item.I_COIN, oct, false)

			} else {
				//EXPを出す
				ItemSatter(ec.X, ec.Y, item.I_EXP, oct, false)
			}

			//爆発。元のスケールを引き継ぐ
			m := ec.Magnification
			effect.SetEffect(int(ec.X), int(ec.Y-float64(ec.Radius)), int(ec.Z), effect.E_EXPLOTION, m)
		}
	}
}

// アイテムばらまき rt=true 手前に飛び散る
func ItemSatter(x, y float64, itemtype int, pieces int, rt bool) {
	for i := 0; i < pieces; i++ {
		vx := float64(-8)
		vy := float64(0)
		r := float64(0)
		if rt { //手前に飛び散る(打撃用)
			r = rand.Float64()*60 + 30
		} else if pieces <= 5 { //5個以下ならランダムに
			r = rand.Float64() * 180
		} else { //5個超えなら整列する
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
	this.EnabledIsTouch = false
	this.gMode = G_GAMEOVER
}

// ゲームオーバー
func UpdateGameOver() {
	this.Counter++
	//音楽がなり始める時間
	if this.Counter == 150 {
		sound.Play(sound.GAMEOVER)
	}
	//誤タッチを防ぐため、しばらく入力は受け付けない
	if this.Counter > 150 {
		if (this.IsTouch && this.EnabledIsTouch) ||
			(this.IskeyPress && this.EnabledIskeyPress) {
			resetGame()
			//セーブを実施
			saveload.SaveAllData(this.Gparam, this.player.PState, *this.enemy)
			this.EnabledIsTouch = false
			this.EnabledIskeyPress = false
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
	this.gMode = G_INFO
	this.EnabledIsTouch = false
	this.EnabledIskeyPress = false
	this.Counter = 20

	camera.X = 0
	player.InitState(50, 5, 10)
	enemy.Reset()
	item.Reset()
	magic.Reset()
	myui.SetUIMode(myui.UI_TITLE)

	SetEnemyState()
}

// レートを内部倍率に変換
func bairitu(rate int) float64 {
	return float64(rate*60/100) + 60
}

// ランクに対する敵のステータス
func SetEnemyState() {
	//表示されるランクと内部で計算に使われるランクは異なる
	rankBase := this.Gparam.Rank

	//ある程度を超えたら内部倍率倍
	rankUpBorder1 := float64(100) //100を超えたら上昇度2倍例110→120
	if rankUpBorder1 < this.Gparam.Rank {
		rankBase += (rankBase - rankUpBorder1)
	}
	rankUpBorder2 := float64(120) //120を超えたら3倍 130→ 160+10 →170
	if rankUpBorder2 < this.Gparam.Rank {
		rankBase += (rankBase - rankUpBorder2)
	}

	//--------------------------------------------
	//防御力は以下の式。全部指数関数
	this.Gparam.EN_BASE_DF = int((rankBase) * ((rankBase) + 10) / 120)

	//--------------------------------------------
	//攻撃力はゆっくりとした第一段階とランク20から更に加算される値の2つで構成されている
	this.Gparam.EN_BASE_AT = int((rankBase)*((rankBase)+10)/600) + 2

	if 20 < rankBase {
		x := int((rankBase-20)*(rankBase-20)/400) + 3
		this.Gparam.EN_BASE_AT += x
	}

	//--------------------------------------------
	//HPは100まであがってそれ以上にはならない
	this.Gparam.EN_BASE_HP = int(rankBase)*10 + 20
	if this.Gparam.EN_BASE_HP > 100 {
		this.Gparam.EN_BASE_HP = 100
	}

	//--------------------------------------------

}
