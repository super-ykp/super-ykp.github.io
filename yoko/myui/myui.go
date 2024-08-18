package myui

import (
	"Yoko/common"
	"Yoko/mytext"
	"Yoko/player"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	pNG_OFFSET_X = 336
	pNG_OFFSET_Y = 0

	hPNG_OFFSET_X = 176
	hPNG_OFFSET_Y = 96

	UI_TITLE = 0
	UI_NOMAL = 1
	UI_MENU  = 2
	UI_HAPPY = 3
)

type UIStruct struct {
	UiMode int

	bar   *ebiten.Image
	sys   *ebiten.Image
	happy []*ebiten.Image

	player      *player.PlayerStruct
	gameparam   *common.GameParam
	isTouchLast bool
	counter     uint
}

var (
	this *UIStruct
)

// --===================================================================================
// 初期化
func Init(t *UIStruct, e *ebiten.Image, pl *player.PlayerStruct, ga *common.GameParam) {
	this = t
	this.player = pl
	this.gameparam = ga
	imageRect := image.Rect(288, 0, 288+32, 1)
	this.bar = e.SubImage(imageRect).(*ebiten.Image)

	imageRect = image.Rect(pNG_OFFSET_X, pNG_OFFSET_Y, pNG_OFFSET_X+16, pNG_OFFSET_Y+16)
	this.sys = e.SubImage(imageRect).(*ebiten.Image)

	offs := [][]int{
		{0, 0, 32, 32},       //
		{32, 0, 32, 32},      //
		{64, 0, 16, 16},      //
		{64 + 16, 0, 16, 16}, //
	}
	for i := 0; i < len(offs); i++ {
		offX := hPNG_OFFSET_X + offs[i][0]
		offy := hPNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]

		imageRect := image.Rect(offX, offy, offW, offH)
		rect := e.SubImage(imageRect).(*ebiten.Image)
		this.happy = append(this.happy, rect)
	}
}

func Update() {
	this.counter++
}

// 描画
func Draw1(screen *ebiten.Image) {
	p := message.NewPrinter(language.Japanese)
	//背景
	vector.DrawFilledRect(screen, 4, 4, 102, 19, color.RGBA{200, 200, 200, 255}, false)
	//HP
	drawBar(screen, 5, 5, 1, color.RGBA{255, 0, 0, 255})
	drawBar(screen, 5, 5, float64(this.player.PState.HP)/float64(this.player.PState.MaxHP), color.RGBA{0, 128, 255, 255})
	mytext.DrawShadowed(screen, 5, 5, p.Sprintf("%d/%d", this.player.PState.HP, this.player.PState.MaxHP))
	//SP
	drawBar(screen, 5, 14, 1, color.RGBA{255, 0, 0, 255})
	drawBar(screen, 5, 14, float64(this.player.PState.SP)/float64(this.player.PState.MaxSP), color.RGBA{56, 255, 56, 255})
	mytext.DrawShadowed(screen, 5, 14, p.Sprintf("%d/%d", this.player.PState.SP, this.player.PState.MaxSP))

	//Sys
	drawImageOption := ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(common.SCREEN_WIDTH-17, 1)
	screen.DrawImage(this.sys, &drawImageOption)
}

func Draw2(screen *ebiten.Image) {
	startY := -float32(common.CAM_Y_OFFSET) + 16

	p := message.NewPrinter(language.Japanese)
	//マスク
	vector.DrawFilledRect(screen, 0, startY, common.SCREEN_WIDTH, common.SCREEN_HEIGHT-startY, color.RGBA{0, 0, 0, 0xff}, false)

	//ゲージ
	yOffsetg := float32(startY) + 5
	vector.DrawFilledRect(screen, 0, yOffsetg-5, common.SCREEN_WIDTH, 17, color.RGBA{255, 255, 255, 0xff}, false)
	vector.DrawFilledRect(screen, 1, yOffsetg-4, common.SCREEN_WIDTH-2, 15, color.RGBA{128, 128, 128, 0xff}, false)
	ps := float64(this.player.PState.EXP%common.G_MAX) / common.G_MAX
	drawGauge(screen, yOffsetg+0, 1, color.RGBA{0, 0, 0, 255})
	drawGauge(screen, yOffsetg+0, ps, color.RGBA{168, 255, 255, 255})
	ps = float64(this.player.PState.GOLD%common.G_MAX) / common.G_MAX
	drawGauge(screen, yOffsetg+3, 1, color.RGBA{0, 0, 0, 255})
	drawGauge(screen, yOffsetg+3, ps, color.RGBA{255, 168, 0, 255})
	ps = float64(this.player.PState.FORCE%common.G_MAX) / common.G_MAX
	drawGauge(screen, yOffsetg+6, 1, color.RGBA{0, 0, 0, 255})
	drawGauge(screen, yOffsetg+6, ps, color.RGBA{56, 255, 56, 255})

	switch this.UiMode {
	case UI_TITLE:
		yOffset := int(startY) + 64 + 15
		yindex := 0
		mytext.DrawN(screen, 5, yOffset+8*yindex, p.Sprintf("操作はタッチかマウス(前進のみスペースキーでも可)"))
		yindex++
		mytext.DrawN(screen, 5, yOffset+8*yindex, p.Sprintf("タッチで前進。後退不能。進むのみ"))
		yindex++
		mytext.DrawN(screen, 5, yOffset+8*yindex, p.Sprintf("ランク100を目指しましょう"))
	case UI_NOMAL, UI_HAPPY:
		//ステータス
		yOffset := int(startY) + 64 + 15
		yindex := 0
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("RANK         %8.3f", this.gameparam.Rank))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("ATTACK   %8d", this.player.PState.AT))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("RES HP   %8d", this.player.PState.RES_HP))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("RES SP   %8d", this.player.PState.RES_SP))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("EXP   RATE %5d%%", this.player.PState.RATE_EXP))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("COIN  RATE %5d%%", this.player.PState.RATE_COIN))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("FORCE RATE %5d%%", this.player.PState.RATE_FORCE))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("SPIKE    %8d", this.player.PState.SPIKE))

		yindex += 2
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("TOTAL EX %8d", this.player.PState.T_EXP/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("TOTAL GD %8d", this.player.PState.T_GOLD/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("TOTAL FC %8d", this.player.PState.T_FORCE/100))

		//DEBUG
		if common.DEBUGSPEED > 1 {
			yOffset = int(startY) + 64 + 50 - 16
			yindex = 0
			mytext.Draw(screen, common.SCREEN_WIDTH/2+30, yOffset+8*yindex, "ENEMY")
			yindex++
			mytext.Draw(screen, common.SCREEN_WIDTH/2+30, yOffset+8*yindex, p.Sprintf(" HP%4d", this.gameparam.EN_BASE_HP))
			yindex++
			mytext.Draw(screen, common.SCREEN_WIDTH/2+30, yOffset+8*yindex, p.Sprintf(" DF%4d", this.gameparam.EN_BASE_DF))
			yindex++
			mytext.Draw(screen, common.SCREEN_WIDTH/2+30, yOffset+8*yindex, p.Sprintf(" AT%4d", this.gameparam.EN_BASE_AT))
		}

	case UI_MENU:
		//ステータス
		yOffset := int(startY) + 64 + 15
		yindex := 0
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("HI RANK         %8.3f", this.gameparam.H_Rank))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("   RANK         %8.3f", this.gameparam.Rank))
		yindex++
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("HI TOTAL EX %8d", this.player.PState.H_T_EXP/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("   TOTAL EX %8d", this.player.PState.T_EXP/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("HI TOTAL GD %8d", this.player.PState.H_T_GOLD/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("   TOTAL GD %8d", this.player.PState.T_GOLD/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("HI TOTAL FC %8d", this.player.PState.H_T_FORCE/100))
		yindex++
		mytext.Draw(screen, 5, yOffset+8*yindex, p.Sprintf("   TOTAL FC %8d", this.player.PState.T_FORCE/100))

		offx := common.SCREEN_WIDTH / 2
		offy := common.SCREEN_HEIGHT * 7 / 8
		vector.DrawFilledRect(screen, float32(offx)-9, float32(offy)-1, 16+2, 10, color.RGBA{255, 0, 0, 0xff}, false)
		mytext.DrawN(screen, offx-8, offy, "自爆")

	}

}
func Draw3(screen *ebiten.Image) {
	if this.UiMode == UI_HAPPY {
		vector.DrawFilledRect(screen, 4, -float32(common.CAM_Y_OFFSET)+20, common.SCREEN_WIDTH-8, 68, color.RGBA{255, 255, 255, 0xff}, false)
		vector.DrawFilledRect(screen, 6, -float32(common.CAM_Y_OFFSET)+22, common.SCREEN_WIDTH-12, 64, color.RGBA{128, 128, 128, 0xff}, false)

		p := message.NewPrinter(language.Japanese)
		yOffset := int(-float32(common.CAM_Y_OFFSET) + 24)
		yindex := 0
		mytext.DrawN(screen, 8, yOffset+8*yindex, p.Sprintf("ランク100到達おめでとう！"))
		yindex++
		mytext.DrawN(screen, 8, yOffset+8*yindex, p.Sprintf("あとは行けるところまで行ってみよう"))
		yindex++
		mytext.DrawN(screen, 8, yOffset+8*yindex, p.Sprintf("「続く」を押してね"))
		//---
		drawImageOption := ebiten.DrawImageOptions{}
		ofY := float64(-float32(common.CAM_Y_OFFSET) + 24 + 24)
		drawImageOption.GeoM.Translate(common.SCREEN_WIDTH-32-8, ofY)
		screen.DrawImage(this.happy[(this.counter/30)%2], &drawImageOption)

		drawImageOption = ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(8, ofY+16)
		screen.DrawImage(this.happy[((this.counter/30)%2)+2], &drawImageOption)
		//----
		offx := common.SCREEN_WIDTH / 2
		offy := common.SCREEN_HEIGHT * 3 / 4
		vector.DrawFilledRect(screen, float32(offx)-10, float32(offy)-2, 16+4, 12, color.RGBA{255, 255, 0, 0xff}, false)
		vector.DrawFilledRect(screen, float32(offx)-9, float32(offy)-1, 16+2, 10, color.RGBA{0, 0, 255, 0xff}, false)
		mytext.DrawN(screen, offx-8, offy, "続く")
	}
}

// 体力バー
func drawBar(screen *ebiten.Image, x, y float32, ps float64, c color.Color) {
	vector.DrawFilledRect(screen, x, y, float32(float64(100)*ps), 8, c, false)
}

// 経験値バー
func drawGauge(screen *ebiten.Image, y float32, ps float64, c color.Color) {
	vector.DrawFilledRect(screen, 5, y, float32(float64(common.SCREEN_WIDTH-10)*ps), 1, c, false)
}

// UIシステム切り替え
func SetUIMode(uiMode int) {
	if this != nil {
		this.UiMode = uiMode
	}
}

func GetUIselect(touchX, touchY int, isTouch bool) int {
	ret := 0
	switch this.UiMode {
	case UI_NOMAL:
		if this.isTouchLast && !isTouch {
			if common.SCREEN_WIDTH-17 < touchX && touchY < 17 {
				ret = -1
			}
		}
	case UI_MENU: //メニュー
		if this.isTouchLast && !isTouch {
			if common.SCREEN_WIDTH/2-8 < touchX && touchX < common.SCREEN_WIDTH/2+8 &&
				common.SCREEN_HEIGHT*7/8 < touchY && touchY < common.SCREEN_HEIGHT*7/8+10 {
				ret = -1
			} else {
				ret = 1
			}
		}
	case UI_HAPPY:
		if this.isTouchLast && !isTouch {
			if common.SCREEN_WIDTH/2-8 < touchX && touchX < common.SCREEN_WIDTH/2+8 &&
				common.SCREEN_HEIGHT*3/4 < touchY && touchY < common.SCREEN_HEIGHT*3/4+10 {
				ret = 1
			}
		}
	}
	this.isTouchLast = isTouch
	return ret
}
