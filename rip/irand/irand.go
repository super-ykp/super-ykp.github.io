package irand

import (
	"Rip/common"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type IrandStruct struct {
	irand   *ebiten.Image
	x       float64
	visible bool
}

var (
	this *IrandStruct
)

// ==========================================================================================
// 島
// ==========================================================================================
// 初期化
func Init(t *IrandStruct, i *ebiten.Image) {
	this = t
	imageRect := image.Rect(96, 0, 96+32, 16)
	this.irand = i.SubImage(imageRect).(*ebiten.Image)
	this.visible = false
}

// 描画
func Draw(screen *ebiten.Image) {
	if !this.visible {
		return
	}
	y := float64(common.SCREEN_HEIGHT - 16)
	x := this.x - 16
	drawImageOption := ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(x, y)
	screen.DrawImage(this.irand, &drawImageOption)
}

// 島座標取得
func GetX() float64 {
	return this.x
}

// 島配置
func Start() {
	this.x = rand.Float64()*float64(common.SCREEN_WIDTH/3) + common.SCREEN_WIDTH/3
	this.visible = true
}

// 消去
func Delete() {
	this.visible = false
}
