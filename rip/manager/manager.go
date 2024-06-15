package manager

import (
	"Rip/bg"
	"Rip/common"
	"Rip/effect"
	"Rip/enemy"
	"Rip/food"
	"Rip/input"
	"Rip/irand"
	"Rip/mytext"
	"Rip/neco"
	"Rip/plane"
	"Rip/sendjs"
	"Rip/sound"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCORE_VER = 1 //スコアの記録バージョン

	MANAGER_WAIT       = 0
	MANAGER_FIRST_INIT = 1
	MANAGER_TITLE      = 2
	MANAGER_GAME_INIT  = 3
	MANAGER_NOMAL_GAME = 4
	MANAGER_GAME_OVER  = 5
)

type ManagerStruct struct {
	manager_status int

	input  *input.InputStruct
	bg     *bg.BgStruct
	plane  *plane.PlaneStruct
	neco   *neco.NecoStruct
	irand  *irand.IrandStruct
	effect *effect.EffectStruct
	enemy  *enemy.EnemyStruct
	food   *food.FoodStruct
	sound  *sound.SoundStruct

	lastFlameisTouch      bool
	lastFlameisPressSpace bool
	inpuTypeTouch         bool

	score     int
	highScore int
	zanki     int
	stage     int

	counter uint

	bar1up     int
	bar1upMax  int
	enemyMax   int
	nowBgIndex int

	left *ebiten.Image
}

var (
	this *ManagerStruct

	m_fpsLastTime time.Time
	m_fpsCount    int
	m_DispFps     int
)

// --===================================================================================
// マネージャのみが全クラスに対する参照を行える
// --===================================================================================
// 初期化
func Init(t *ManagerStruct, i *ebiten.Image) {
	this = t

	imageRect := image.Rect(128, 0, 128+8, 8)
	t.left = i.SubImage(imageRect).(*ebiten.Image)

	this.manager_status = MANAGER_FIRST_INIT
}

// クラスの保持
func SetClass(bg *bg.BgStruct, p *plane.PlaneStruct, n *neco.NecoStruct, i *irand.IrandStruct, e *effect.EffectStruct, en *enemy.EnemyStruct, fd *food.FoodStruct, inp *input.InputStruct, so *sound.SoundStruct) {
	this.bg = bg
	this.plane = p
	this.neco = n
	this.irand = i
	this.effect = e
	this.enemy = en
	this.food = fd
	this.input = inp
	this.sound = so
}

// 計算
func Update() {
	//fps計算（デバッグ用)
	m_fpsCount++
	currentTime := time.Now()
	if m_fpsLastTime.IsZero() || currentTime.Sub(m_fpsLastTime) >= time.Second {
		m_fpsLastTime = time.Now()
		m_DispFps = m_fpsCount
		m_fpsCount = 0
	}

	//-----------------------------------------------
	//初回初期化
	if this.manager_status == MANAGER_FIRST_INIT {
		//webasmなら、javascript経由でクッキーにアクセスし、ハイスコアを取得する
		ver := sendjs.Get("VERSION")
		high := sendjs.Get("HISCORE")
		_, err1 := strconv.Atoi(ver)
		_, err2 := strconv.Atoi(high)
		//取得できない(exe起動)場合ハイスコア0
		if err1 != nil || err2 != nil {
			this.highScore = 0
		} else {
			//スコアバージョンが違っていたらハイスコアリセット
			ver, _ := strconv.Atoi(ver)
			if ver == SCORE_VER {
				this.highScore, _ = strconv.Atoi(high)
			}
		}
		//ゲームの初期化に
		this.manager_status = MANAGER_GAME_INIT
	}

	//デフォルト画面初期化
	if this.manager_status == MANAGER_GAME_INIT {
		//タイトルへ
		this.manager_status = MANAGER_TITLE
		enemy.Start(5)
	}

	//........................................
	//スペースキー押下でtrue
	isPressSpace := input.IsPressSpace()
	//タッチされていたらtrue
	isTouch := input.IsTouch()
	//デフォルトはスペースキーのトグル切り替え。タッチが検出されたらそちらに
	if isTouch {
		this.inpuTypeTouch = true
	}

	if this.manager_status == MANAGER_TITLE { //タイトル.................................
		//スペースキーかタッチで開始
		if (!isTouch && this.lastFlameisTouch) || (isPressSpace && !this.lastFlameisPressSpace) {
			//ブラウザはアクティブ化しないと音が鳴らせない。このタイミングでダミー音を鳴らす
			sound.Play(sound.START)
			//ゲームオーバー音楽がなっていたら止める
			sound.Stop(sound.GAMEOVER, 60)
			//開始
			Start()
			//初期化
			this.inpuTypeTouch = false
			//ゲームへ
			this.manager_status = MANAGER_NOMAL_GAME
		}
	} else if this.manager_status == MANAGER_NOMAL_GAME { //ゲーム本編
		//ねこ状態取得
		neco_mode := neco.GetState()

		if neco_mode == neco.NECO_RESTARTOK { //ミス後リスタート要請中............................
			//残機を減らす
			this.zanki -= 1
			//残0、ゲームオーバー
			if this.zanki < 0 {
				this.manager_status = MANAGER_GAME_OVER
				this.counter = 0
				//ゲームオーバー音楽再生
				sound.Play(sound.GAMEOVER)
				//ハイスコア記録
				sendjs.Send("VERSION", strconv.Itoa(SCORE_VER))
				sendjs.Send("HISCORE", strconv.Itoa(this.highScore))
			} else { //まだリトライ可能
				//敵を減らす
				if this.enemyMax > 9 {
					this.enemyMax -= 6
				}
				//敵とひこーき再起動
				plane.Start()
				neco.Start()
				enemy.Start(this.enemyMax)
			}
		} else if neco_mode == neco.NECO_NEXTSTAGEOK { //クリア後次ステージ要請............................
			//次ステージへ
			this.stage += 1
			//敵増える
			this.enemyMax += 3
			//ステージ再起動
			StageStart()
		} else if neco_mode == neco.NECO_RAID { //ひこーきにねこが乗っている.........................................
			//ひこーきは画面外で消える。ただしねこがのっているので再出現
			if !plane.GetVisible() {
				plane.Start()
			}

			//プレイヤーによる操作
			if plane.GetVisible() {
				//ひこーきはねこが乗っている限り左右で加減速できる
				if input.GetPressLR(false) < 0 {
					plane.SetSpped(-0.75)
				} else if input.GetPressLR(false) > 0 {
					plane.SetSpped(1)
				}
				//タップ➝離す、又は　スペースキー押下で　ねこフォール
				if (!isTouch && this.lastFlameisTouch) || (isPressSpace && !this.lastFlameisPressSpace) {
					x, y := plane.GetXY()
					neco.Fall(x, y)
					sound.Play(sound.FALL)
				}
			}
		} else if neco_mode == neco.NECO_FALL || neco_mode == neco.NECO_PARA {
			if neco_mode == neco.NECO_FALL { //ねこ落下中...........
				//タップ中、又は　スペースキー押下で　パラシュート
				if (this.inpuTypeTouch && isTouch) || (isPressSpace && !this.lastFlameisPressSpace) {
					neco.Para(true)
					sound.Play(sound.PARAOPEN)
				}
			} else if neco_mode == neco.NECO_PARA { //ねこパラ中...........
				//タップ中、又は　スペースキー押下で　パラシュートたたむ
				if this.inpuTypeTouch && !isTouch || (isPressSpace && !this.lastFlameisPressSpace) {
					neco.Para(false)
					sound.Play(sound.PARACLOSE)
				}
			}

			//左右移動.......................................................
			neco.Move(input.GetPressLR(true))

			//接地判定用に取得
			neco_x, neco_y := neco.GetXY()
			irand_x := irand.GetX()

			if neco_y > common.SCREEN_HEIGHT-16-8 { //島の高さを超えたら....
				//島に乗っている
				if irand_x-16-6 < neco_x && neco_x < irand_x+16+6 { //成功
					//得点に応じてBGを変える。その場合非同期で先行読込する
					AddScore(100)
					NextBgChangeCheck()
					//ねこ歓喜の舞
					neco.Success()
					sound.Play(sound.SUCCESS)
					//100点エフェクト
					effect.Start(neco_x, neco_y, effect.E_P100)
				} else { //乗っていない
					//ミス中除く
					if neco_mode != neco.NECO_MISS {
						//ミス開始
						neco.Miss(0)
						sound.Play(sound.MISS)
					}
				}
			} else { //島の高さに来ていない。（落下中）.............
				//敵との接触判定
				enemyArray := enemy.GetEnemyArray()
				hit := false
				for _, e := range enemyArray {
					if e.Visible && math.Abs(e.X-neco_x) < 14/2 && math.Abs(e.Y-neco_y) < 12/2 {
						hit = true
					}
				}
				//敵と接触
				if hit {
					//ミス
					neco.Miss(-1.5)
					sound.Play(sound.MISS)
				}

				//フードとの接触判定
				foodArray := food.GetFoorArray()
				for i := 0; i < len(foodArray); i++ {
					f := &foodArray[i]
					if !f.Visible {
						continue
					}
					//接触
					if math.Abs(f.X-neco_x) < 8+4 && math.Abs(f.Y-neco_y) < 8+4 {
						//表示から消す
						f.Visible = false
						//取ったフードが1Up
						if f.Type == food.ONEUP {
							this.zanki++
							effect.GetOneUpEffect(neco_x, neco_y)
							sound.Play(sound.ONEUP)
						} else { //フード
							effect.Start(neco_x, neco_y, effect.E_P10)
							AddScore(10)
							sound.Play(sound.CATCH)
						}
					}
				}
			}
		}
	} else if this.manager_status == MANAGER_GAME_OVER { //ゲームオーバー
		this.counter++
		//1秒間は誤タッチを考慮して操作無視
		if this.counter > 60*1 {
			//キーが離されて再度押される、またはタッチ状態がオフ状態からオンになったら、タイトルへ
			if (!isTouch && this.lastFlameisTouch) || (isPressSpace && !this.lastFlameisPressSpace) {
				this.manager_status = MANAGER_GAME_INIT
				this.nowBgIndex = 0
				bg.SetBg(this.nowBgIndex)
				irand.Delete()
				food.AllDelete()
				enemy.Start(5)
			}
		}
	}
	//前回フレームのキー情報を次のために保持
	this.lastFlameisTouch = isTouch
	this.lastFlameisPressSpace = isPressSpace
}

// --===================================================================================
// 描画
// --===================================================================================
func Draw(screen *ebiten.Image) {
	c := color.RGBA{255, 255, 255, 255}

	//ハイスコアは常に表示
	st := "HIGH  " + fmt.Sprintf("%06d", this.highScore)
	mytext.ScoreDraw(screen, 0, 9, st, c)

	//ゲーム中とゲームオーバー時のみスコア表示
	if this.manager_status == MANAGER_NOMAL_GAME || this.manager_status == MANAGER_GAME_OVER {
		//スコア
		st = "SCORE " + fmt.Sprintf("%06d", this.score)
		mytext.ScoreDraw(screen, 0, 17, st, c)

		//残機
		for i := 0; i < this.zanki; i++ {
			x := common.SCREEN_WIDTH - 10 - 10*i
			y := 3
			drawImageOption := ebiten.DrawImageOptions{}
			drawImageOption.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(this.left, &drawImageOption)
		}
	}

	if this.manager_status == MANAGER_GAME_OVER { //ゲームオーバーの文字...........
		st = "GAME OVER"
		of := common.SCREEN_HEIGHT/2 + 16 - int(this.counter)
		if of < 0 {
			of = 0
		}
		c := color.RGBA{255, 0, 0, 255}
		mytext.Draw(screen, common.SCREEN_HEIGHT/2+int(of), st, 0, c)
	} else if this.manager_status == MANAGER_TITLE { //タイトル文字..................
		st = "NECO Parachute"
		c := color.RGBA{0, 0, 255, 255}
		mytext.Draw(screen, common.SCREEN_HEIGHT/2-40, st, 0, c)

		st = "Press space key or tap screen."
		c = color.RGBA{0, 255, 255, 255}
		mytext.Draw(screen, common.SCREEN_HEIGHT/2+40, st, 1, c)
		//DEBUG-------------
		//if common.Builddate != "" {
		//	mytext.ScoreDraw(screen, 0, common.SCREEN_WIDTH, "BUILD>"+common.Builddate, color.RGBA{255, 255, 255, 255})
		//}
	}
	//DEBUG--------------------
	//str := fmt.Sprintf("%d ", m_DispFps)
	//text.Draw(screen, str, mplusNormalFont1, 0, common.SCREEN_HEIGHT, color.RGBA{0, 0, 0, 255})
}

// --===================================================================================
// 他
// --===================================================================================
// スコア加算
func AddScore(ad int) {
	this.score += ad
	if this.highScore < this.score {
		this.highScore = this.score
	}
	this.bar1up += ad
}

// ゲーム開始
func Start() {
	this.bar1upMax = 1000
	this.bar1up = 0
	this.score = 0
	this.zanki = 2
	this.stage = 1
	this.enemyMax = 4
	StageStart()
}

// ステージの開始
func StageStart() {
	plane.Start()
	neco.Start()
	irand.Start()
	enemy.Start(this.enemyMax)
	food.Start(irand.GetX(), this.stage*4)
	bg.SetBg(this.nowBgIndex)
	CheckOneUp()
}

// 1Up出現チェック及びアイテムの出現
func CheckOneUp() {
	if this.bar1upMax <= this.bar1up {
		this.bar1up -= this.bar1upMax
		this.bar1upMax += 1000
		x, y := food.Push1Up(irand.GetX())
		if x != 0 {
			sound.Play(sound.PUSHONEUP)
			effect.PushOneUpEffect(x, y)
		}
	}
}

func NextBgChangeCheck() {
	var boder = [3]int{3000, 6000, 9000}
	for i := 0; i < len(boder); i++ {
		if this.nowBgIndex == i && boder[i] <= this.score {
			this.nowBgIndex = i + 1
			bg.GoAsyncBgLoad(this.nowBgIndex)
		}
	}

}
