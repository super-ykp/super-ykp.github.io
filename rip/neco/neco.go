package neco

import (
	"Rip/common"
	"Rip/sound"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	NECO_RAID        = 0
	NECO_FALL        = 1
	NECO_PARA        = 2
	NECO_MISS        = 4
	NECO_SUCCESS     = 5
	NECO_RESTARTOK   = 6
	NECO_NEXTSTAGEOK = 7

	NECO_MISS_STATE1 = 0
	NECO_MISS_STATE2 = 1
	NECO_MISS_STATE3 = 4
)
const (
	LRSPPED    = 0.05
	LRSPPEDMAX = 1.5
)

type NecoStruct struct {
	sound     *sound.SoundStruct
	neco      [3]*ebiten.Image
	neco_para *ebiten.Image
	neco_miss *ebiten.Image

	x       float64
	y       float64
	v_x     float64
	v_y     float64
	counter int

	neco_mode    int
	NMS          int
	neco_imgType int
}

var (
	this *NecoStruct
)

// ==========================================================================================
//  ねこ
//==========================================================================================

// 初期化/////////////////////////////////////////////////////////////////
func Init(p *NecoStruct, i *ebiten.Image, so *sound.SoundStruct) {
	this = p
	this.sound = so
	//ねこイメージ
	imageRect := image.Rect(16, 0, 16+16, 16)
	p.neco[0] = i.SubImage(imageRect).(*ebiten.Image)
	imageRect = image.Rect(16*2, 0, 16*2+16, 16)
	p.neco[1] = i.SubImage(imageRect).(*ebiten.Image)
	imageRect = image.Rect(16*3, 0, 16*3+16, 16)
	p.neco[2] = i.SubImage(imageRect).(*ebiten.Image)

	//パラシュート
	imageRect = image.Rect(0, 16, 16, 16+16)
	p.neco_para = i.SubImage(imageRect).(*ebiten.Image)

	//MISS吹き出し
	imageRect = image.Rect(64, 0, 64+32, 16)
	p.neco_miss = i.SubImage(imageRect).(*ebiten.Image)

	//ステータス--------
	p.neco_mode = NECO_RAID
	p.neco_imgType = 0

	p.x = 0
	p.y = 0
	p.v_x = 0
	p.v_y = 0
	p.counter = 0
}

// 計算
func Update() {
	//カウンター加算
	this.counter += 1

	//モードによる制御--------------------------------------------------------
	if this.neco_mode == NECO_FALL { //落下................
		//重力
		this.v_y += 0.1

	} else if this.neco_mode == NECO_PARA { //パラシュート............
		this.v_y = 0.5

	} else if this.neco_mode == NECO_SUCCESS { //着地成功............
		this.neco_imgType = (this.counter / 10) % 2
		//わーい
		if this.counter > 60*2 {
			this.neco_mode = NECO_NEXTSTAGEOK
		}

	} else if this.neco_mode == NECO_MISS { //ミス....................
		if this.NMS == NECO_MISS_STATE1 {
			this.v_y += 0.1
			//画面外へ消えたら
			if this.y > common.SCREEN_HEIGHT+30 {
				//ウエイト開始
				this.counter = 0
				this.NMS = NECO_MISS_STATE2
				sound.Play(sound.MISS2)
			}
		} else if this.NMS == NECO_MISS_STATE2 {
			//ミス表示ウエイト到達
			if this.counter > 60*1 {
				this.neco_mode = NECO_RESTARTOK
			}
		}
	}

	//画面端を考慮した移動--------------------------------------------------------
	if (this.v_x < 0 && this.x < 8) || (this.v_x > 0 && common.SCREEN_WIDTH-8 < this.x) || this.neco_mode == NECO_SUCCESS || this.neco_mode == NECO_MISS {
		this.v_x = 0
	}
	this.x += this.v_x
	this.y += this.v_y
}

// 描画
func Draw(screen *ebiten.Image) {
	//パラシュート
	if this.neco_mode == NECO_PARA {
		x := this.x - 8
		y := this.y - 8 - 15

		drawImageOption := ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(x, y)
		screen.DrawImage(this.neco_para, &drawImageOption)
	}

	//ねこ
	if this.neco_mode != NECO_RAID {
		drawImageOption := ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(this.x-8, this.y-8)
		screen.DrawImage(this.neco[this.neco_imgType], &drawImageOption)
	}

	//ミス吹き出し
	if this.neco_mode == NECO_MISS && this.NMS == NECO_MISS_STATE2 {
		x := this.x - 32
		if x < 0 {
			x = 0
		}
		if common.SCREEN_WIDTH-32 < x {
			x = common.SCREEN_WIDTH - 32
		}

		y := common.SCREEN_WIDTH - this.counter*4
		if y < common.SCREEN_HEIGHT-32+8 {
			y = common.SCREEN_HEIGHT - 32 + 8
		}

		drawImageOption := ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(x, float64(y))
		screen.DrawImage(this.neco_miss, &drawImageOption)
	}
}

// 外部呼び出し////////////////////////////////////////////////////////////
// 状態取得
func GetState() int {
	return this.neco_mode
}

// 座標取得
func GetXY() (float64, float64) {
	return this.x, this.y
}

// ひこーきに乗せる
func Start() {
	this.neco_mode = NECO_RAID
	this.neco_imgType = 0
}

// 落下開始
func Fall(x, y float64) {
	this.neco_mode = NECO_FALL
	this.x = x
	this.y = y

	if this.x < 8 {
		this.x = 8
	}
	if common.SCREEN_WIDTH-8 < this.x {
		this.x = common.SCREEN_WIDTH - 8
	}

	this.v_x = 0 //ベクトル
	this.v_y = 0
}

// パラシュートオン・オフ
func Para(on bool) {
	if on {
		this.neco_mode = NECO_PARA
	} else {
		this.neco_mode = NECO_FALL
	}
}

// 左右移動
func Move(ac float64) {
	if ac < 0 && -LRSPPEDMAX < this.v_x {
		this.v_x -= LRSPPED
	} else if 0 < ac && this.v_x < LRSPPEDMAX {
		this.v_x += LRSPPED
	}

}

// ミス
func Miss(v_y float64) {
	this.neco_mode = NECO_MISS
	this.neco_imgType = 2
	this.v_y = v_y
	this.NMS = NECO_MISS_STATE1
}

// 着地成功
func Success() {
	this.neco_mode = NECO_SUCCESS
	this.y = common.SCREEN_WIDTH - 16 - 8
	this.v_x = 0
	this.v_y = 0
	this.counter = 0
}
