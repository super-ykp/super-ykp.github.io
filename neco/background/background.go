package background

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BG_WIDTH  = 320
	BG_HEIGHT = 640
)

type BackgroundStruct struct {
	animCounter int
}

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage *ebiten.Image
)

// -----------------------------------------------------------------
func Update(b *BackgroundStruct) {

	b.animCounter++
}

func Draw(screen *ebiten.Image, b *BackgroundStruct) {
	//ebitenライブラリのDrawImageOptions構造体をインスタンス化
	drawImageOption := &ebiten.DrawImageOptions{}

	//1枚目----------
	//PNG内位置
	imageRect := image.Rect(0, 0, BG_WIDTH, BG_HEIGHT)
	//イメージ取得
	ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
	//スクロール位置
	OffsetY := float64(b.animCounter % BG_HEIGHT)
	drawImageOption.GeoM.Translate(0, OffsetY)
	//描画
	screen.DrawImage(ebitenImage, drawImageOption)

	//2枚目---------
	drawImageOption = &ebiten.DrawImageOptions{}
	//PNG内位置
	imageRect = image.Rect(0, 0, BG_HEIGHT, BG_HEIGHT)
	//イメージ取得
	ebitenImage = m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
	//スクロール位置
	drawImageOption.GeoM.Translate(0, OffsetY-BG_HEIGHT)
	//描画
	screen.DrawImage(ebitenImage, drawImageOption)
}

func Init(e *ebiten.Image, b *BackgroundStruct) {
	m_MainPngImage = e
}
