package enemy

import (
	"Neco/common"
	"Neco/danmaku"
	"Neco/explosion"
	"Neco/player"
	"Neco/sound"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pngOffsetX = 0
	pngOffsetY = 16

	pngOffsetBossX = 128
	pngOffsetBpssY = 16

	HpBarPngImgOffsetX = 0
	HpBarPngImgOffsetY = 48
	HpBarWidth         = 128
	enemyTypeZako      = 0
	enemyTypeBoss      = 1

	BOSSMOVEMODE_WAIT      = -1
	BOSSMOVEMODE_NOMAL     = 0
	BOSSMOVEMODE_EXPLOSION = 1
)

type TekiStruct struct {
	IsUse       bool    //利用中ならtrue
	EnemyWidth  float64 //幅（pngサイズ)
	EnemyHeight float64

	X float64 //座標
	Y float64

	VectorX float64 //ベクトル
	VectorY float64

	enemyType int //0雑魚　1ボス

	movemode     int
	flameCounter int //毎フレーム加算されるカウンタ

	targetDistance float64 //移動ターゲットへの距離

	knockbackCount  int     //ノックバック中プラス
	collisoinRadius float64 //当たり判定用半径

	Hp    int
	HpMax int
}

type EnemyStruct struct {
	EnemyList [200]TekiStruct
}

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージ
	//弾幕クラスのポインタ
	m_This *EnemyStruct
	m_dan  *danmaku.DanmakuStruct
	m_ex   *explosion.ExplosionStruct
	m_s    *sound.SoundStruct

	//プレイヤー座標のキャッシュ
	m_PlayerX float64
	m_PlayerY float64

	//敵出現位置（キャッシュ）
	m_StartDistance float64
	//敵速度
	m_EnemySpeed float64
	//敵接近可能ボーダー
	m_BorderDistance float64

	//全敵制御カウンタ
	m_GlovalCounter int

	//雑魚の弾のスピード
	m_ZakoTamaSpeed float64
)

// ボス
var (
	b_OnHit  bool
	b_x      float64 //汎用バッファ
	b_y      float64
	b_count1 int
	b_count2 int
	b_count3 int
)

// -----------------------------------------------------------------
// 計算
func Update(e *EnemyStruct, p *player.PlayerStruct) {
	m_PlayerX = p.X
	m_PlayerY = p.Y
	m_GlovalCounter++

	//雑魚出現
	if common.GameState == common.GameStateNomal {
		StartZako(e)
	}

	//雑魚処理
	zakoUpdate(e)
	//ボス処理
	bossUpdate(e, p)
}

// 雑魚動作
func zakoUpdate(e *EnemyStruct) {
	//配列は固定長
	for i := 0; i < len(e.EnemyList)-1; i++ {
		t := &e.EnemyList[i]
		//未使用なら飛ばす
		if !t.IsUse {
			continue
		}

		//画面外なら消滅----------
		if t.Y < -120 || t.Y > common.SCREEN_HEIGHT/2 && (t.X < -t.EnemyWidth/2 || common.SCREEN_WIDTH+t.EnemyWidth/2 < t.X) || t.Y > common.SCREEN_HEIGHT+t.EnemyWidth {
			t.IsUse = false
			continue
		}
		//移動ターゲットへの相対距離------------
		t.targetDistance = distance(t.X, t.Y, common.SCREEN_WIDTH/2, common.SCREEN_HEIGHT)

		//移動処理--------------------------------------------------------------------------
		xv := common.SCREEN_WIDTH/2 - t.X
		yv := common.SCREEN_HEIGHT - t.Y
		v := float64(math.Sqrt(xv*xv + yv*yv))
		if t.movemode == 0 {
			//非ノックバック中かつ、ターゲットに接近しすぎていたら

			if t.targetDistance < m_BorderDistance-30 && t.knockbackCount == 0 {
				t.VectorX += ((xv/v)*m_EnemySpeed/10 - t.VectorX) / 10
				t.VectorY += ((0)*m_EnemySpeed/10 - t.VectorY) / 10
				t.movemode = 1
			} else if t.targetDistance < m_BorderDistance && t.knockbackCount == 0 {
				//減速
				t.VectorX += ((xv/v)*m_EnemySpeed/10 - t.VectorX) / 10
				t.VectorY += ((yv/v)*m_EnemySpeed/10 - t.VectorY) / 10

			} else {
				//移動。移動方向は指定方向へ直接ではなく、慣性を持つようになっている
				t.VectorX += ((xv/v)*m_EnemySpeed - t.VectorX) / 10
				t.VectorY += ((yv/v)*m_EnemySpeed - t.VectorY) / 10
			}
		} else {
			t.VectorX += (0*m_EnemySpeed - t.VectorX) / 10
			t.VectorY += 0.1

		}

		//ベクトル加算
		t.X += t.VectorX
		t.Y += t.VectorY

		//ノックバック中以外なら、敵同士が重ならないようにする
		if t.knockbackCount == 0 {
			//敵同士で当たり判定をとる
			i := collisoin(e, i)
			//ぶつかった対象が見つかったら
			if i != -1 {
				//距離が遠い方を外側ベクトルにする
				if t.targetDistance > e.EnemyList[i].targetDistance {
					r := rand.Float64() * float64(60-30)
					xv, yv = rotateVector(-xv, -yv, r)
					t.VectorX = ((xv / v) * m_EnemySpeed)
					t.VectorY = ((yv / v) * m_EnemySpeed)
				}
			}
		} else {
			t.knockbackCount--
		}

		//攻撃処理-----------------------------------------------------------------
		if m_ZakoTamaSpeed != 0 && (t.flameCounter%10) == 0 && t.targetDistance < m_BorderDistance+30 && (m_GlovalCounter%240) > 90 {
			danmaku.ShotZ(m_dan, t.X+rand.ExpFloat64()*t.EnemyWidth-t.EnemyWidth/2, t.Y, m_ZakoTamaSpeed, 0)
		}

		//毎フレーム加算
		t.flameCounter++
	}
}

// ボス動作----------------------------------------------------------------------------------------
func bossUpdate(e *EnemyStruct, p *player.PlayerStruct) {
	t := &e.EnemyList[len(e.EnemyList)-1]
	if !t.IsUse {
		return
	}

	t.flameCounter++
	b_OnHit = false

	if t.movemode == BOSSMOVEMODE_NOMAL {
		//移動処理-------------------
		xv := common.SCREEN_WIDTH/2 - (common.SCREEN_WIDTH/2 - m_PlayerX) - t.X
		yv := float64(common.SCREEN_HEIGHT/4) - t.Y
		v := float64(math.Sqrt(xv*xv + yv*yv))

		t.VectorX += ((xv/v)*m_EnemySpeed/10 - t.VectorX) / 10
		t.VectorY += ((yv/v)*m_EnemySpeed - t.VectorY) / 10

		//ベクトル加算
		t.X += t.VectorX
		t.Y += t.VectorY

		//攻撃処理-----------------------------------------------------------------
		if t.Y > common.SCREEN_HEIGHT/4-60 {
			bossAttack(e, false)
		}

	} else if t.movemode == BOSSMOVEMODE_EXPLOSION { //ボス爆発
		if t.flameCounter%5 == 0 {
			explosion.Explosion(m_ex, t.X+rand.Float64()*64-32, t.Y+rand.Float64()*64-32, 1, 4)
			sound.Play(m_s, sound.EXPLOSION)
		}
		if t.flameCounter >= 60*2 {
			explosion.Explosion(m_ex, t.X, t.Y, 3, 7)
			sound.Stop(m_s, sound.EXPLOSION, 0)
			sound.Play(m_s, sound.EXPLOSION_BOSS)
			t.IsUse = false
			common.StartStage()
		}
	} else {
		if t.flameCounter > 60*5 {
			t.movemode = BOSSMOVEMODE_NOMAL
		}
	}
}

// -----------------------------------------------------------------
// 描画
func Draw(screen *ebiten.Image, d *EnemyStruct) {
	//配列は固定長
	for i := 0; i < len(d.EnemyList); i++ {
		//未使用なら飛ばす
		if !d.EnemyList[i].IsUse {
			continue
		}
		//ebitenライブラリのDrawImageOptions構造体をインスタンス化
		drawImageOption := &ebiten.DrawImageOptions{}
		//中心位置に対する描画のずれを考慮し、配置する
		drawImageOption.GeoM.Translate(d.EnemyList[i].X-d.EnemyList[i].EnemyWidth/2, d.EnemyList[i].Y-d.EnemyList[i].EnemyHeight/2)

		//PNG内位置
		OffsetX := 0
		if d.EnemyList[i].enemyType == enemyTypeZako { //雑魚
			OffsetX = pngOffsetX + ((d.EnemyList[i].flameCounter/5)%4)*int(d.EnemyList[i].EnemyWidth)
		} else {
			OffsetX = pngOffsetBossX + ((d.EnemyList[i].flameCounter/20)%2)*int(d.EnemyList[i].EnemyWidth)
		}
		imageRect := image.Rect(OffsetX, pngOffsetY, OffsetX+int(d.EnemyList[i].EnemyWidth), pngOffsetY+int(d.EnemyList[i].EnemyHeight))
		//イメージ取得
		ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)
	}
}

// 描画
func Draw2(screen *ebiten.Image, e *EnemyStruct) {
	t := e.EnemyList[len(e.EnemyList)-1]
	if !t.IsUse || t.movemode != BOSSMOVEMODE_NOMAL {
		return
	}
	//PNG内バー位置
	imageRect := image.Rect(HpBarPngImgOffsetX, HpBarPngImgOffsetY, HpBarPngImgOffsetX+HpBarWidth, HpBarPngImgOffsetY+8)
	//イメージ取得
	ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
	drawImageOption := &ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(float64(t.X-HpBarWidth/2), float64(t.Y-32))
	//描画
	screen.DrawImage(ebitenImage, drawImageOption)

	//PNG内バー位置
	imageRect = image.Rect(HpBarPngImgOffsetX, HpBarPngImgOffsetY+8, HpBarPngImgOffsetX+int(float64(HpBarWidth)*float64(t.Hp)/float64(t.HpMax)), HpBarPngImgOffsetY+8+8)
	//イメージ取得
	ebitenImage = m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
	drawImageOption = &ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(float64(t.X-HpBarWidth/2), float64(t.Y-32))
	//描画
	screen.DrawImage(ebitenImage, drawImageOption)

}

// 初期化------------------------------------------------------------------
func Init(e *ebiten.Image, en *EnemyStruct, d *danmaku.DanmakuStruct, ex *explosion.ExplosionStruct, s *sound.SoundStruct) {
	//イメージ所持
	m_MainPngImage = e
	//構造体所持
	m_This = en
	m_dan = d
	m_ex = ex
	m_s = s

	//敵速度
	m_EnemySpeed = 1
	//敵接近可能ボーダー
	m_BorderDistance = common.SCREEN_HEIGHT * 2 / 3

	//敵出現位置（画面中央最下段からの円形）
	m_StartDistance = distance(common.SCREEN_WIDTH/2, common.SCREEN_HEIGHT, 0-64, 0-64)
}

func Start(e *EnemyStruct) {
	for i := 0; i < len(e.EnemyList); i++ {
		e.EnemyList[i].IsUse = false
	}

	StartBoss(e)
}

// 出現----------------------------------------------------------------
// 雑魚
func StartZako(e *EnemyStruct) {
	//バッファの空き取得
	t := getTamaBuf(e)
	//もう無いなら無理
	if t == nil {
		return
	}

	//初期配置は、画面最下部中央を中心に、垂直から左右30度、距離m_StartDistanceの位置に配置される
	haba := float64(30)
	r := float64(math.Pi*(-90-haba)/180) + (rand.Float64() * math.Pi * haba * 2 / 180)
	t.IsUse = true
	t.enemyType = 0
	t.X, t.Y = calculateNewCoordinates(common.SCREEN_WIDTH/2, common.SCREEN_HEIGHT, r, m_StartDistance)
	t.EnemyWidth = 32
	t.EnemyHeight = 32
	t.collisoinRadius = 8
	t.movemode = 0
	t.Hp = int(3 + rand.ExpFloat64()*3)

}

// ボス
func StartBoss(e *EnemyStruct) {
	//バッファの空き取得
	t := &e.EnemyList[len(e.EnemyList)-1]
	t.IsUse = true
	t.enemyType = enemyTypeBoss
	t.X = common.SCREEN_WIDTH / 2
	t.Y = -(32 + 32)
	t.EnemyWidth = 32
	t.EnemyHeight = 32 + 16
	t.collisoinRadius = 8
	t.movemode = BOSSMOVEMODE_WAIT
	t.flameCounter = 0
	t.HpMax = 150

	//初期化
	bossAttack(e, true)

	t.Hp = t.HpMax
}

// 空きバッファの返却---------------------------
func getTamaBuf(d *EnemyStruct) *TekiStruct {
	i := 0
	//開いているバッファを探す
	for ; i < len(d.EnemyList)-1; i++ {
		if !d.EnemyList[i].IsUse {
			return &d.EnemyList[i]
		}
	}
	return nil
}

// 指定座標からr角d距離離れた場所の座標を得る--------------------------
func calculateNewCoordinates(x, y, r, d float64) (nx, ny float64) {
	nx = x + d*math.Cos(r)
	ny = y + d*math.Sin(r)
	return nx, ny
}

// 指定座標間の直選距離を求める------------------------------------
func distance(x, y, nx, ny float64) float64 {
	dx := nx - x
	dy := ny - y
	return math.Sqrt(dx*dx + dy*dy)
}

// 敵同士の接触--------------------------------------------------------------
func collisoin(d *EnemyStruct, my int) int {
	x := d.EnemyList[my].X
	y := d.EnemyList[my].Y
	r := d.EnemyList[my].collisoinRadius
	//自分以外ループ
	for i := 0; i < len(d.EnemyList); i++ {
		en2 := d.EnemyList[i]
		if en2.IsUse && i != my {
			dis := distance(en2.X, en2.Y, x, y)
			if dis <= en2.collisoinRadius+r {
				return i
			}
		}
	}
	return -1
}

// ----------------------------------------------------------------------
// 敵同士接触時のベクトル変化。p自分、c対象。cへの接触接線に対し並行方向のベクトルが維持されるように変化する
func collisoinVectorCalc(p *TekiStruct, c *TekiStruct) {
	// 円の中心と点座標の差ベクトル
	dx := p.X - c.X
	dy := p.Y - c.Y

	// 接線方向ベクトル
	tangentX := dy / math.Sqrt(dx*dx+dy*dy)
	tangentY := -dx / math.Sqrt(dx*dx+dy*dy)

	// 点座標の運動量と接線方向ベクトルの内積
	dot := p.VectorX*tangentX + p.VectorY*tangentY

	// 接線方向に並行な運動量
	parallelX := dot * tangentX
	parallelY := dot * tangentY

	// 垂直方向の運動量
	perpendicularX := p.VectorX - parallelX
	perpendicularY := p.VectorY - parallelY

	// 反射後の運動量
	p.VectorX = parallelX - perpendicularX
	p.VectorY = parallelY - perpendicularY
}

// ショットへのヒット------------------------------------------------------------
func Hit(t *TekiStruct) {
	//ボーダーを超えない限り当たらない
	if t.targetDistance > m_BorderDistance {
		return
	}
	if t.enemyType == enemyTypeBoss && b_OnHit {
		return
	}
	//体力減少
	t.Hp--
	b_OnHit = true

	//0なら消滅
	if t.Hp == 0 {
		if t.enemyType == enemyTypeBoss { //ボス撃破
			//爆発モードへ
			t.movemode = BOSSMOVEMODE_EXPLOSION
			t.flameCounter = 0
			//敵弾消去
			danmaku.Start(m_dan)
			//雑魚全滅
			for i := 0; i < len(m_This.EnemyList)-1; i++ {
				if !m_This.EnemyList[i].IsUse {
					continue
				}
				t := &m_This.EnemyList[i]
				t.IsUse = false
			}
			common.StartGameStateBossDefeat()

		} else {
			t.IsUse = false
		}
		explosion.Explosion(m_ex, t.X, t.Y, 1, 1)
		sound.Play(m_s, sound.EXPLOSION)
		return
	}

	//雑魚で体力余りなら
	if t.enemyType == enemyTypeZako {
		sound.Play(m_s, sound.KNOKBACK)
		//プレイヤーと逆方向へノックバック
		xv := m_PlayerX - t.X
		yv := m_PlayerY - t.Y
		v := float64(math.Sqrt(xv*xv + yv*yv))
		r := rand.Float64() * float64(60-30)
		xv, yv = rotateVector(xv, yv, r)
		knockbackSpeed := float64(10)
		t.VectorX = (xv / v) * -knockbackSpeed
		t.VectorY = (yv / v) * -knockbackSpeed
		t.knockbackCount = 60
	} else {
		sound.Play(m_s, sound.BOSSHIT)
	}
}

// ベクトルを回転させる----------------------------
func rotateVector(x, y, angle float64) (float64, float64) {
	rad := angle * (math.Pi / 180.0)
	newX := x*math.Cos(rad) - y*math.Sin(rad)
	newY := x*math.Sin(rad) + y*math.Cos(rad)
	return newX, newY
}

// ==============================================================================================
// ボス攻撃
// ==============================================================================================
func bossAttack(e *EnemyStruct, s bool) {
	t := &e.EnemyList[len(e.EnemyList)-1]

	switch common.Stage {
	case 1:
		stage1(t, s)
	case 2:
		stage2(t, s)
	case 3:
		stage3(t, s)
	case 4:
		stage4(t, s)
	case 5:
		stage5(t, s)
	case 6:
		stage6(t, s)
	case 7:
		stage7(t, s)
	case 8:
		stage8(t, s)
	case 9:
		stage9(t, s)
	case 10:
		stage10(t, s)
	default:
		stage1(t, s)
	}
}
