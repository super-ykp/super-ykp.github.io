package mytext

import (
	"Yoko/common"
	"Yoko/fileloader"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusNormalFont1 *text.GoTextFace
	mplusNormalFont2 *text.GoTextFace
	mplusNormalFont3 *text.GoTextFace
)

// --===================================================================================
// テキスト病が
// --===================================================================================
// 初期化
func Init() {
	//-------------------------------------------
	//io.Reader を得る
	f, err := fileloader.Open("font/pressstart2p.ttf")
	if err != nil {
		panic(err)
	}

	// フォントを読み込む
	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		panic(err)
	}
	mplusNormalFont1 = &text.GoTextFace{Source: src, Size: 8}
	mplusNormalFont3 = &text.GoTextFace{Source: src, Size: 14}
	//----------------------------------------------------------------
	f, err = fileloader.Open("font/misaki_gothic.ttf")
	if err != nil {
		panic(err)
	}

	// フォントを読み込む
	src, err = text.NewGoTextFaceSource(f)
	if err != nil {
		panic(err)
	}
	mplusNormalFont2 = &text.GoTextFace{Source: src, Size: 8}

}

// 描画
func Draw(screen *ebiten.Image, x, y int, txt string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, txt, mplusNormalFont1, op)
}

// 描画影付き。HPとか
func DrawShadowed(screen *ebiten.Image, x, y int, txt string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+1, float64(y)+1)
	op.ColorScale.Scale(0, 0, 0, 255)
	text.Draw(screen, txt, mplusNormalFont1, op)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+1, float64(y))
	text.Draw(screen, txt, mplusNormalFont1, op)
}

// 日本語
func DrawN(screen *ebiten.Image, x, y int, txt string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = 10
	text.Draw(screen, txt, mplusNormalFont2, op)
}

// 描画
func DrawG(screen *ebiten.Image, y int, txt string, ftyp int, c color.Color) {

	r, g, b, a := c.RGBA()

	// RGBAの値を0-255の範囲に変換
	r /= 257 // 65535 / 255 = 257
	g /= 257
	b /= 257
	a /= 257

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(common.SCREEN_WIDTH/2-len(txt)*15/2+1), float64(y)+1)
	op.ColorScale.Scale(0, 0, 0, 255)
	text.Draw(screen, txt, mplusNormalFont3, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(common.SCREEN_WIDTH/2-len(txt)*15/2-1), float64(y)-1)
	op.ColorScale.Scale(255, 255, 255, 255)
	text.Draw(screen, txt, mplusNormalFont3, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(common.SCREEN_WIDTH/2-len(txt)*15/2), float64(y))
	op.ColorScale.Scale(float32(r), float32(g), float32(b), float32(a))
	text.Draw(screen, txt, mplusNormalFont3, op)

}
