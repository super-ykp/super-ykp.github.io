package enemy

import (
	"Yoko/camera"
	"Yoko/common"
	"Yoko/player"
	"image"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type EcStruct struct {
	IsUse    bool
	EStete   uint //状態
	SizeType int

	MHP int
	HP  int
	DF  int

	X float64 //X座標。ワールド位置
	Y float64 //Y座標
	Z int     //Zオーダー

	vx float64 //ジャンプ中ベクトル
	vy float64

	Vibration     int     //ヒットストップ振動
	Magnification float64 //拡大率
	Radius        int     //半径
	rotate        int     //回転
	DispPngIndex  uint    //表示PNGインデックス
	Counter       uint    //カウンタ
	BakusanY      int     //爆散するY座標
}

const (
	ENEMY_MAX    = 20
	pNG_OFFSET_X = 64
	pNG_OFFSET_Y = 0
	GRAVITY      = 0.1

	S_STAND         = 0  //待機
	S_DMG           = 1  //ダメージ
	S_DMGBROW       = 2  //ダメージ後倒れ
	S_BROW          = 3  //吹き飛び中
	S_DOWN          = 4  //ダウン中
	S_DMGLEVING1    = 5  //ダメージ後跳ね回り爆発
	S_DMGLEVING2    = 6  //ダメージ後倒れ爆発
	S_DMGLEVING3    = 7  //ダメージ後倒れ爆発（アイテム出さない）
	S_LEVING1       = 8  //退場中
	S_LEVING2       = 9  //退場中
	S_LEVING3       = 10 //退場中(爆発時、消滅)
	S_DMGBROWREVNGE = 11 //ダメージ後、吹き飛び、リベンジ
	S_ROLL          = 12 //距離を取る
	S_JAMP          = 13 //ジャンプ
	S_TAME          = 14 //ため
	S_ATTACK        = 15 //攻撃
	S_RETURN        = 16 //攻撃後戻る
	S_KAISYUMACHI   = 17 //マネージャの回収判定待ち
	S_KAISYUMACHI2  = 18 //マネージャの回収判定待ち(アイテム出さない)
)

type EnemyStruct struct {
	enemyImage   []*ebiten.Image
	hpBar        *ebiten.Image
	EList        []EcStruct `json:"EList"`
	AttackBuffer [10]bool
}

var (
	this      *EnemyStruct
	EnemyType = []float64{0.6, 0.7, 1, 1.2, 2}
)

// --===================================================================================
// 初期化
func Init(t *EnemyStruct, e *ebiten.Image, p *player.PlayerStruct) {
	this = t
	offs := [][]int{
		{0, 0, 32, 32},   //0 待機
		{32, 0, 32, 32},  //1
		{64, 0, 32, 32},  //2 ダメージ
		{96, 0, 32, 32},  //3 ふっとび
		{128, 0, 32, 32}, //4 倒れてる
		{160, 0, 32, 32}, //5 回転
		{192, 0, 32, 32}, //6 怒り
	}
	for i := 0; i < len(offs); i++ {
		offX := pNG_OFFSET_X + offs[i][0]
		offy := pNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]

		imageRect := image.Rect(offX, offy, offW, offH)
		rect := e.SubImage(imageRect).(*ebiten.Image)
		this.enemyImage = append(this.enemyImage, rect)
	}
	//HPバー
	imageRect := image.Rect(288, 0, 288+32, 1)
	this.hpBar = e.SubImage(imageRect).(*ebiten.Image)
}

// リセット
func Reset() {
	for i := 0; i < len(this.AttackBuffer); i++ {
		this.AttackBuffer[i] = false
	}
	this.EList = []EcStruct{}
}

// 計算
func Update() {
	//敵バッファ全周
	for i := 0; i < len(this.EList); i++ {
		//未使用ならスキップ
		ec := &this.EList[i]
		if !ec.IsUse {
			continue
		}

		ec.Counter++

		//ステータスに応じた処理
		switch ec.EStete {
		case S_STAND: //立っている/////////////////////////////////////////////////////////////////
			ec.DispPngIndex = (ec.Counter/30)%2 + 0

		case S_DMG, S_DMGLEVING1, S_DMGLEVING2, S_DMGLEVING3, S_DMGBROW, S_DMGBROWREVNGE: //ダメージエフェクト中................
			//表示するPNG
			ec.DispPngIndex = 2
			//ヒットストップ振動
			v := int(1 * ec.Magnification)
			if v < 1 {
				v = 1
			}
			ec.Vibration = v * (-1 * int((ec.Counter/5)%2))

			//ダメージ後体力0の退場は、以下三種
			//S_DMGLEVING跳ね回り爆発
			//S_DMGLEVING2そのまま落下爆発
			//S_DMGBROWLEVINGそのまま落下爆発（アイテム出さない）
			if ec.EStete == S_DMGLEVING1 || ec.EStete == S_DMGLEVING2 || ec.EStete == S_DMGLEVING3 { //HP0時の最後のヒットストップなら.............
				//退場に移行
				ec.Counter = 0
				//退場ベクトル設定
				if ec.EStete == S_DMGLEVING2 || ec.EStete == S_DMGLEVING3 {
					ec.vx = 2
					ec.vy = -3
					if ec.EStete == S_DMGLEVING3 {
						ec.EStete = S_LEVING3
					} else {
						ec.EStete = S_LEVING2
					}
				} else {
					ec.vx = 10 + rand.Float64()*15
					ec.vy = -(2.5 + rand.Float64()*2.5)
					ec.EStete = S_LEVING1
				}
				ec.BakusanY = -(rand.Intn(30))

			} else if ec.EStete == S_DMGBROW { //ダメージ後倒れる..........................
				//ふっとぶ
				ec.EStete = S_BROW
				ec.vx = 2
				ec.vy = -1

				ec.Counter = 0
			} else if ec.EStete == S_DMGBROWREVNGE { //ダメージ後ころがり、反撃..........................
				ec.EStete = S_ROLL
				ec.vx = 3
				ec.vy = -2
				ec.Counter = 0

			} else { //状態遷移なしにしばらくしたら立ち状態に復帰...............................
				if ec.Counter > 30 {
					ec.EStete = S_STAND
				}
			}
		case S_ROLL: //ころがり
			ec.DispPngIndex = 5
			ec.rotate = (ec.rotate + 1) % 360
			ec.X += 2
			ec.vy += GRAVITY
			//座標にベクトル加算
			ec.Y += ec.vy
			if 0 < ec.Y {
				ec.Y = 0
			}
			if ec.Counter > 60*5 || int(ec.X)+(ec.Radius) > camera.X+common.SCREEN_WIDTH/2 {
				ec.EStete = S_JAMP
				ec.vx = 0
				ec.vy = -3 + rand.Float64()*0.5
				ec.Counter = 0
			}
		case S_JAMP:
			ec.DispPngIndex = 6
			ec.vy += GRAVITY
			//座標にベクトル加算
			ec.Y += ec.vy
			if 0 < ec.vy {
				ec.EStete = S_TAME
				ec.Counter = uint(rand.Intn(20))
				ec.rotate = 0
			}
		case S_TAME: //攻撃ため
			ec.DispPngIndex = 6
			v := 1
			ec.Vibration = v * (-1 * int((ec.Counter/2)%2))
			if ec.Counter > 60 {
				ec.Counter = 0
				ec.EStete = S_ATTACK
			}
		case S_ATTACK: //攻撃
			ec.DispPngIndex = 6
			xv := float64(camera.X) - (ec.X)
			yv := -16 - (ec.Y)
			v := float64(math.Sqrt(xv*xv + yv*yv))

			//ベクトル設定
			ec.vx = (xv / v) * 20
			ec.vy = (yv / v) * 20
			ec.X += ec.vx
			ec.Y += ec.vy
			//ヒット
			if (ec.X) < float64(camera.X) {
				ec.EStete = S_RETURN
				ec.Counter = 0
				ec.vx = 2 + rand.Float64()*2
				ec.vy = -2

				for i := 0; i < len(this.AttackBuffer); i++ {
					if !this.AttackBuffer[i] {
						this.AttackBuffer[i] = true //ダメージ
						break
					}
				}
				if player.GetThis().PState.SP > 0 {
					spdmg := (player.GetThis().PState.SPIKE - ec.DF)
					if spdmg < 1 {
						spdmg = 1
					}
					ec.HP -= spdmg
					if ec.HP <= 0 {
						ec.HP = 0
						ec.EStete = S_DMGLEVING3
						ec.Counter = 0
					}
				}
			}
		case S_RETURN:
			//PNGイメージ
			ec.DispPngIndex = 6

			ec.X += ec.vx
			ec.Y += ec.vy
			ec.vy += GRAVITY
			//着地したら
			if ec.Y >= 0 {
				ec.Y = 0
				//通常に
				ec.EStete = S_STAND
				ec.Counter = 0
			}

		case S_BROW: //吹き飛ばされ中.......................................................................
			//PNGイメージ
			ec.DispPngIndex = 3
			//ベクトルでふっとぶ
			ec.X += ec.vx
			ec.Y += ec.vy
			ec.vy += GRAVITY
			//着地したら
			if ec.Y >= 0 {
				ec.Y = 0
				//ダウン状態に
				ec.EStete = S_DOWN
				ec.Counter = 0
			}

		case S_DOWN: //ダウン中.......................................................................
			//PNGイメージ
			ec.DispPngIndex = 4
			ec.rotate = 0
			//一定時間で起き上がる
			if ec.Counter > 30 && ec.HP > 0 {
				ec.EStete = S_STAND
			}

		case S_LEVING1, S_LEVING2, S_LEVING3: //退場モード.......................................................................
			//退場モードは
			// S_LEVING　回転跳ね回り
			// S_LEVING2　ダウン
			// S_LEVING3 ダウン（アイテム出さない）
			if ec.EStete == S_LEVING2 || ec.EStete == S_LEVING3 {
				ec.DispPngIndex = 4
			} else {
				ec.DispPngIndex = 5
				ec.rotate = (ec.rotate + 5) % 360
			}
			//画面上外ではタイマー加算停止
			if ec.Y < common.CAM_Y_OFFSET+16 {
				ec.Counter--
			}
			//地面で跳ね返る
			if 0 < ec.Y {
				ec.vy = -ec.vy * 0.09 / 1
			}
			//画面左右端で跳ね返る
			if int(ec.X) < camera.X-common.SCREEN_WIDTH/2 || (camera.X+common.SCREEN_WIDTH/2 < int(ec.X) && ec.EStete != S_LEVING3) {
				v := float64(-1)
				if camera.X-common.SCREEN_WIDTH/2 > int(ec.X) {
					v = 1
				}
				ec.vx = math.Abs(ec.vx) * v
			}
			//ベクトル補正。空気抵抗と重力
			ec.vx = ec.vx * 0.95 / 1
			ec.vy += GRAVITY
			//座標にベクトル加算
			ec.X += ec.vx
			ec.Y += ec.vy

			//走馬灯を見れる時間が終わった
			if float64(ec.BakusanY) < ec.Y && ec.vy > 0 {
				if ec.EStete == S_LEVING3 {
					ec.EStete = S_KAISYUMACHI2
				} else {
					ec.EStete = S_KAISYUMACHI
				}
			}
		case S_KAISYUMACHI, S_KAISYUMACHI2:
			//上位存在による回収待ち
		}
	}

}

// 魔法によるダメージモーション
func Update2() {
	//敵バッファ全周
	for i := 0; i < len(this.EList); i++ {
		//未使用ならスキップ
		ec := &this.EList[i]
		if !ec.IsUse {
			continue
		}

		ec.Counter++
		//ステータスに応じた処理
		switch ec.EStete {
		case S_DMG, S_DMGLEVING1, S_DMGLEVING2, S_DMGLEVING3, S_DMGBROW, S_DMGBROWREVNGE: //ダメージエフェクト中................
			//表示するPNG
			ec.DispPngIndex = 2
			//ヒットストップ振動
			v := int(1 * ec.Magnification)
			if v < 1 {
				v = 1
			}
			ec.Vibration = v * (-1 * int((ec.Counter/5)%2))
		}
	}
}

// 描画。都合上プレイヤーの描画を割り込ませる
func Draw(screen *ebiten.Image) {
	//Zソートする
	sort.Slice(this.EList, func(i, j int) bool {
		return this.EList[i].Z < this.EList[j].Z
	})

	//プレイヤー表示フラグ
	isPdraw := true

	for i := 0; i < len(this.EList); i++ {
		//敵バッファが利用中でないなら次
		ec := &this.EList[i]
		if !ec.IsUse {
			continue
		}

		//次の敵の表示座標が、0を超えたらプレイヤーを割り込ませる
		if isPdraw && ec.Z >= 0 {
			isPdraw = false
			player.Draw(screen)
		}

		//敵表示
		drawImageOption := ebiten.DrawImageOptions{}

		//拡大
		//回転
		if ec.rotate > 0 {
			drawImageOption.GeoM.Translate(-32/2, -32/2)
			drawImageOption.GeoM.Rotate(float64(ec.rotate))
			drawImageOption.GeoM.Translate(32/2, 32/2)
		}
		drawImageOption.GeoM.Translate(-32/2, -32)
		drawImageOption.GeoM.Scale(ec.Magnification, ec.Magnification)
		drawImageOption.GeoM.Translate(32/2, 32)

		ofX := int(ec.X-(32)/2) - camera.CamOffsetX + ec.Vibration
		ofY := int(ec.Y-(32)) + ec.Z - camera.CamOffsetY
		drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))
		screen.DrawImage(this.enemyImage[ec.DispPngIndex], &drawImageOption)

		//----------------------------------------------------
		//HPバー
		drawImageOption = ebiten.DrawImageOptions{}
		drawImageOption.ColorScale.Scale(1, 0, 0, 1)
		ofX = int(ec.X-(32)/2) - camera.CamOffsetX + ec.Vibration
		ofY = ec.Z + int(ec.Y) - int(32*ec.Magnification) - 5 - camera.CamOffsetY
		drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))
		screen.DrawImage(this.hpBar, &drawImageOption)

		//HPバー
		drawImageOption = ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Scale(float64(ec.HP)/float64(ec.MHP), 1)
		drawImageOption.ColorScale.Scale(1, 1, 1, 1)
		ofX = int(ec.X-(32)/2) - camera.CamOffsetX + ec.Vibration
		ofY = ec.Z + int(ec.Y) - int(32*ec.Magnification) - 5 - camera.CamOffsetY
		drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))
		screen.DrawImage(this.hpBar, &drawImageOption)
	}

	//プレイヤーがここまで描画されないこともあるのでその場合はここで描く
	if isPdraw || player.IsDown() {
		isPdraw = false
		player.Draw(screen)
	}
}

// 敵を配置する
func SetEnemy(x int, etype int, mHP int, df int) {
	ec := getBuffer()
	if ec == nil {
		return
	}

	ec.IsUse = true
	ec.X = float64(x)
	ec.Y = 0
	ec.Z = rand.Intn(10) - 4
	ec.rotate = 0
	//サイズ
	ec.SizeType = etype
	ec.Magnification = EnemyType[etype]
	ec.Radius = int(16 * ec.Magnification)

	//パラメータ
	hp := mHP
	switch etype {
	case 0, 1:
		hp = int(math.Ceil(float64(hp) / 6))
	case 2:
		hp = hp * 1
	case 3:
		hp = 1 + hp*2
	case 4:
		hp = 1 + hp*4
	}
	ec.MHP = hp + rand.Intn(3) - 1
	ec.HP = hp
	ec.DF = df + rand.Intn(3) - 1
	//カウンタ系初期化
	ec.Counter = 0
	ec.EStete = S_STAND

}

// バッファの空きを探し返却、なければ追加する
func getBuffer() *EcStruct {

	for i := 0; i < len(this.EList); i++ {
		ec := &this.EList[i]
		if ec.IsUse {
			continue
		}
		return ec
	}
	//敵の上限（一応）
	if ENEMY_MAX <= len(this.EList) {
		return nil
	}
	ep := EcStruct{}
	this.EList = append(this.EList, ep)
	return &this.EList[len(this.EList)-1]
}

// 敵ステータス取得
func GetThis() *EnemyStruct {
	return this
}

// 攻撃を受け付ける状態
func IsHitOk(ec *EcStruct) bool {
	return ec.EStete == S_STAND || ec.EStete == S_DMG || ec.EStete == S_DOWN
}

// 状態強制変更
func SetState(ec *EcStruct, state uint) {
	if ec.EStete != state {
		ec.EStete = state
		ec.Counter = 0
	}
}
