package text

import (
	"Neco/common"
	"image/color"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	m_mplusNormalFont font.Face

	m_Counter int
)

func Init() {
	//フォントの初期化
	const dpi = 72
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	m_mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func Update() {
	m_Counter++
}

func DrawText(screen *ebiten.Image, s string, y int, fonttype int) {
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	c := black
	if fonttype == 0 || (m_Counter%20) < 10 {
		c = black
	} else {
		c = white
	}

	text.Draw(screen, s, m_mplusNormalFont, common.SCREEN_WIDTH/2-(len(s)*14)/2, y, c)

}
