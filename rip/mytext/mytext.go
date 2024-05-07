package mytext

import (
	"Rip/common"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont1 font.Face
	mplusNormalFont2 font.Face
)

// --===================================================================================
// テキスト病が
// --===================================================================================
// 初期化
func Init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont1, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	mplusNormalFont2, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// 描画、スコア用
func ScoreDraw(screen *ebiten.Image, x, y int, txt string, c color.Color) {
	text.Draw(screen, txt, mplusNormalFont1, x+1, y+1, color.Black)
	text.Draw(screen, txt, mplusNormalFont1, x, y, c)
}

// 描画
func Draw(screen *ebiten.Image, y int, txt string, ftyp int, c color.Color) {
	cb := color.RGBA{0, 0, 0, 255}
	cw := color.RGBA{255, 255, 255, 255}
	switch ftyp {
	case 0:
		text.Draw(screen, txt, mplusNormalFont2, common.SCREEN_WIDTH/2-len(txt)*15/2-1, y-1, cw)
		text.Draw(screen, txt, mplusNormalFont2, common.SCREEN_WIDTH/2-len(txt)*15/2+1, y+1, cb)
		text.Draw(screen, txt, mplusNormalFont2, common.SCREEN_WIDTH/2-len(txt)*15/2, y, c)
	case 1:
		text.Draw(screen, txt, mplusNormalFont1, common.SCREEN_WIDTH/2-len(txt)*8/2-1, common.SCREEN_HEIGHT/2+40-1, cw)
		text.Draw(screen, txt, mplusNormalFont1, common.SCREEN_WIDTH/2-len(txt)*8/2+1, common.SCREEN_HEIGHT/2+40+1, cb)
		text.Draw(screen, txt, mplusNormalFont1, common.SCREEN_WIDTH/2-len(txt)*8/2, common.SCREEN_HEIGHT/2+40, c)
	}
}
