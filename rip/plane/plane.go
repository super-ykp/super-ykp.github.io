package plane

import (
	"Rip/common"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlaneStruct struct {
	sp_plane [2]*ebiten.Image
	x        float64
	y        float64
	addSpeed float64
	counter  int
	visible  bool
}

var (
	this *PlaneStruct
)

const (
	PLANE_NOMAL_SPEED = 1
)

// ==========================================================================================
// ひこーき
// ==========================================================================================
// 初期化
func Init(t *PlaneStruct, i *ebiten.Image) {
	this = t
	//インスタンスimgを必要なイメージにスライスする
	//PNG内プレイヤー位置
	imageRect := image.Rect(16, 16, 16+24, 16+16)
	t.sp_plane[0] = i.SubImage(imageRect).(*ebiten.Image)
	imageRect = image.Rect(16+(24), 16, 16+(24)+24, 16+16)
	t.sp_plane[1] = i.SubImage(imageRect).(*ebiten.Image)

	t.x = -64
	t.y = 8 + 16
	t.counter = 0
	t.visible = false
}

// 計算
func Update() {
	//ひこーき移動
	this.x += PLANE_NOMAL_SPEED + this.addSpeed
	this.addSpeed = 0

	if common.SCREEN_WIDTH+(24)/2 < this.x {
		this.visible = false
	}
	//カウンタ加算
	this.counter += 1
}

// 描画
func Draw(screen *ebiten.Image) {
	if !this.visible {
		return
	}
	n := (this.counter / 5) % 2
	s := this.sp_plane[n]
	x := this.x - (24)/2
	y := this.y - 8
	drawImageOption := ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(x, y)
	screen.DrawImage(s, &drawImageOption)
}

// 出現させる
func Start() {
	this.x = -24
	this.visible = true
}

// 加減速
func SetSpped(addSpeed float64) {
	this.addSpeed = addSpeed
}

// 座標取得
func GetXY() (float64, float64) {
	return this.x, this.y
}

// 表示状態取得
func GetVisible() bool {
	return this.visible
}
