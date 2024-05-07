package food

import (
	"Rip/common"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ONEUP = 20
)

type FdStruct struct {
	X       float64
	Y       float64
	Type    int
	Visible bool
}

type FoodStruct struct {
	fc       [200]FdStruct
	food_img []*ebiten.Image
	counter  uint
}

var (
	this *FoodStruct
)

// ==========================================================================================
// フード
// ==========================================================================================
// 初期化
func Init(t *FoodStruct, im *ebiten.Image) {
	this = t
	//イメージ
	for i := 0; i < 22; i++ {
		imageRect := image.Rect(0+(i%10)*16, 32+(i/10)*16, 0+(i%10)*16+16, 32+(i/10)*16+16)
		t.food_img = append(t.food_img, im.SubImage(imageRect).(*ebiten.Image))
	}
}

// 計算
func Update() {
	this.counter++
}

// 描画
func Draw(screen *ebiten.Image) {
	for _, f := range this.fc {
		if !f.Visible {
			continue
		}
		p := f.Type

		x := f.X - 16/2
		y := f.Y - 16/2
		drawImageOption := ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(x, y)

		if p == ONEUP && (this.counter/2)%2 == 0 {
			p = ONEUP + 1
		}

		screen.DrawImage(this.food_img[p], &drawImageOption)
	}
}

// フードを配置する。xは中心位置
func Start(x float64, max int) {
	AllDelete()
	//配列最上位は1UP用に開けておく
	for i := 0; i < len(this.fc)-1; i++ {
		if i >= max {
			break
		}
		f := &this.fc[i]
		f.X = rand.Float64()*(common.SCREEN_WIDTH/2) + float64(x) - common.SCREEN_WIDTH/4
		f.Y = rand.Float64()*(common.SCREEN_WIDTH-80) + 40
		f.Type = i % 20
		f.Visible = true
	}
}

// 1Up配置
func Push1Up(x float64) (float64, float64) {
	for i := len(this.fc) - 1; 0 <= i; i-- {
		f := &this.fc[i]
		if f.Visible {
			continue
		}
		f.X = rand.Float64()*(common.SCREEN_WIDTH/4) + float64(x) - common.SCREEN_WIDTH/8
		f.Y = common.SCREEN_WIDTH * 2 / 3
		f.Type = 20
		f.Visible = true
		return f.X, f.Y
	}
	return 0, 0
}

// 全部クリア
func AllDelete() {
	for i := 0; i < len(this.fc); i++ {
		f := &this.fc[i]
		f.Visible = false
	}
}

// 配列取得
func GetFoorArray() *[200]FdStruct {
	return &this.fc
}
