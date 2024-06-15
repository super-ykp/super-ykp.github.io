package explosion

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pngOffsetX = 0
	pngOffsetY = 64
	pngWidth   = 64
	pngHighte  = 64
)

type bakuStruct struct {
	IsUse bool //利用中ならtrue

	X     float64 //座標
	Y     float64
	Width int
	Hight int
	Scale float64

	flameCounter    int
	maxflameCounter int
	animCounter     int //アニメーションカウンタ
}

type ExplosionStruct struct {
	BakuList [150]bakuStruct
}

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージ
)

// -----------------------------------------------------------------
// 計算
func Update(e *ExplosionStruct) {
	//配列は固定長
	for i := 0; i < len(e.BakuList); i++ {
		b := &e.BakuList[i]
		//未使用なら飛ばす
		if !b.IsUse {
			continue
		}

		b.flameCounter++
		if b.flameCounter >= b.maxflameCounter {
			b.flameCounter = 0
			b.animCounter++
			if b.animCounter == 16 {
				b.IsUse = false
				continue
			}
		}
	}
}

// 描画---------------------------------------------------------------
func Draw(screen *ebiten.Image, e *ExplosionStruct) {
	//配列は固定長
	for i := 0; i < len(e.BakuList); i++ {
		b := &e.BakuList[i]
		//未使用なら飛ばす
		if !b.IsUse {
			continue
		}

		//ebitenライブラリのDrawImageOptions構造体をインスタンス化
		drawImageOption := &ebiten.DrawImageOptions{}
		//中心位置に対する描画のずれを考慮し、配置する
		drawImageOption.GeoM.Scale(b.Scale, b.Scale)
		drawImageOption.GeoM.Translate(b.X-float64(b.Width)/2, b.Y-float64(b.Hight)/2)
		//PNG内位置
		animOffset := (b.animCounter)

		offsetX := pngOffsetX + (animOffset%4)*pngWidth
		offsetY := pngOffsetY + (animOffset/4)*pngHighte
		imageRect := image.Rect(offsetX, offsetY, offsetX+pngWidth, offsetY+pngHighte)
		//イメージ取得
		ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)

	}
}

// ---------------------------------------------------------
// 初期化
func Init(e *ebiten.Image, ex *ExplosionStruct) {
	//イメージ所持
	m_MainPngImage = e
}

func Start(e *ExplosionStruct) {
	for i := 0; i < len(e.BakuList); i++ {
		e.BakuList[i].IsUse = false
	}
}

// ------------------------------------------------------------
// 爆発
func Explosion(e *ExplosionStruct, x, y float64, scale float64, maxflameCount int) {
	for i := 0; i < len(e.BakuList); i++ {
		b := &e.BakuList[i]
		//使用中なら飛ばす
		if b.IsUse {
			continue
		}
		b.IsUse = true
		b.maxflameCounter = maxflameCount
		b.flameCounter = -1
		b.animCounter = 0
		b.X = x
		b.Y = y
		b.Scale = scale
		b.Width = int(pngWidth * b.Scale)
		b.Hight = int(pngHighte * b.Scale)
		return
	}
}
