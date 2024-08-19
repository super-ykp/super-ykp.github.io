package magic

import (
	"Yoko/camera"
	"Yoko/common"
	"Yoko/skill"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ==========================================================================================
//即席スキルエフェクト。エフェクト中は一部ゲーム部が止まる
// ==========================================================================================

const (
	pNG_OFFSET_X  = 240
	pNG_OFFSET_Y  = 64
	pNG2_OFFSET_X = 288
	pNG2_OFFSET_Y = 64

	E_NONE            = 0
	E_ELECTRICTRIGGER = 1
	E_THUNDERBOLT     = 2
	E_LIGHTNINGSWORD  = 3
)

type MagicStruct struct {
	magicImage  []*ebiten.Image
	magicImage2 *ebiten.Image
	effectS     int //エフェクトタイプ
	Counter     int
	X           int
	Y           int
	otherVal    float64 //汎用フラグ
}

var (
	this *MagicStruct
)

// --===================================================================================
// 初期化
func Init(t *MagicStruct, p *ebiten.Image) {
	this = t
	//魔法に使うイメージが偶然同じ大きさで横並びだったからできる荒業
	for i := 0; i < 3; i++ {
		offX := pNG_OFFSET_X + i*16
		offW := offX + 16
		offy := pNG_OFFSET_Y
		offH := pNG_OFFSET_Y + 32
		imageRect := image.Rect(offX, offy, offW, offH)
		rect := p.SubImage(imageRect).(*ebiten.Image)
		t.magicImage = append(t.magicImage, rect)
	}
	imageRect := image.Rect(pNG2_OFFSET_X, pNG2_OFFSET_Y, pNG2_OFFSET_X+32, pNG2_OFFSET_Y+32)
	this.magicImage2 = p.SubImage(imageRect).(*ebiten.Image)
}

// 計算
func Update() bool {
	//各魔法は発動から一定フレームで終了となる。エフェクトの終わりをみているのではない
	this.Counter++
	switch this.effectS {
	case E_ELECTRICTRIGGER:
		if this.Counter > 10 {
			this.effectS = E_NONE
		}
		return false
	case E_THUNDERBOLT:
		if this.Counter > 20 {
			this.effectS = E_NONE
		}
		return false
	case E_LIGHTNINGSWORD:
		if this.Counter > 60 {
			this.effectS = E_NONE
		}
		return false
	}

	return true
}

// 描画
func Draw(screen *ebiten.Image) {
	switch this.effectS {
	case E_ELECTRICTRIGGER: //エレクトリッガー
		//チカっと光る
		vector.DrawFilledRect(screen, 0, 0, common.SCREEN_WIDTH, -common.CAM_Y_OFFSET, color.RGBA{255, 255, 255, 255}, false)
	case E_THUNDERBOLT: //サンダーボルト
		//5フレーム光る
		if this.Counter < 5 {
			vector.DrawFilledRect(screen, 0, 0, common.SCREEN_WIDTH, -common.CAM_Y_OFFSET, color.RGBA{255, 255, 255, 255}, false)
		}
		//縦3キャラで構成された稲妻が左右から交差する
		im := (this.Counter / 3) % 3
		for i := 0; i < 3; i++ {
			//稲妻1本目
			drawImageOption := ebiten.DrawImageOptions{}
			ofX := int(this.X) - camera.CamOffsetX + this.Counter*8
			ofY := int(this.Y) - camera.CamOffsetY + 32*i - 32*3
			drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))
			screen.DrawImage(this.magicImage[im], &drawImageOption)

			//稲妻2本目
			drawImageOption = ebiten.DrawImageOptions{}
			ofX = common.SCREEN_WIDTH/2 + int(this.X) - camera.CamOffsetX - this.Counter*8
			drawImageOption.GeoM.Translate(float64(ofX), float64(ofY)+8)
			screen.DrawImage(this.magicImage[im], &drawImageOption)
		}
	case E_LIGHTNINGSWORD: //マジカル艦砲射撃
		sy := float64(1)
		//ビームの太さ変化はリニアだと気持ちよくないため、ゆっくり現れて一気に太くなり、最後また加速するようにしぼむ
		if this.Counter < 60 {
			if 45 < this.Counter {
				this.otherVal += this.otherVal
				sy = 1 - this.otherVal
				if sy < 0 {
					sy = 0
				}
			} else if 40 < this.Counter {
				sy = 1.3
			} else {
				this.otherVal = 0.2
			}
			//横方向に敷き詰める
			for i := 0; i <= common.SCREEN_WIDTH/32+2; i++ {
				drawImageOption := ebiten.DrawImageOptions{}
				drawImageOption.GeoM.Scale(1, sy)
				ofX := float64(int(this.Counter*11+32-camera.CamOffsetX)%32) - 32
				ofy := +common.CAM_Y_OFFSET - float64(camera.CamOffsetY) - 16*sy + 16
				drawImageOption.GeoM.Translate(ofX+float64(i*32), ofy)
				screen.DrawImage(this.magicImage2, &drawImageOption)
			}
		}
		if this.Counter < 60 {
			a := (255 - this.Counter*9)
			if a <= 0 {
				a = 0
			}
			vector.DrawFilledRect(screen, 0, 0, common.SCREEN_WIDTH, -common.CAM_Y_OFFSET, color.RGBA{uint8(a), uint8(a), uint8(a), uint8(a)}, false)
		}
	}
}

// 魔法エフェクトの発動
func MagicGo(mtype int) {
	switch mtype {
	case skill.S_ELECTRICTRIGGER:
		this.effectS = E_ELECTRICTRIGGER
		this.Counter = 0
	case skill.S_THUNDERBOLT:
		this.effectS = E_THUNDERBOLT
		this.X = camera.X
		this.Counter = 0
	case skill.S_LIGHTNINGSWORD:
		this.effectS = E_LIGHTNINGSWORD
		this.X = camera.X
		this.Counter = 0

	}
}
func GetMagicActive() bool {
	return this.effectS != E_NONE
}

func Reset() {
	this.effectS = E_NONE
}
