package bg

import (
	"Rip/common"
	"Rip/loader"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type BgStruct struct {
	bg_img      *ebiten.Image
	BgPingImage [10]*ebiten.Image
	nowBgIndex  int
}

var (
	this *BgStruct
)

// --===================================================================================
// 　背景
// --===================================================================================
// 初期化
func Init(t *BgStruct) {
	this = t
	this.nowBgIndex = 0
	//最初のBGは、同期で読み込む
	readBg(0)
}

// 描画
func Draw(screen *ebiten.Image) {
	//対象の背景が読込中なら
	if this.BgPingImage[this.nowBgIndex] == nil {
		return
	}
	imageRect := image.Rect(0, 0, common.SCREEN_WIDTH, common.SCREEN_HEIGHT)
	this.bg_img = this.BgPingImage[this.nowBgIndex].SubImage(imageRect).(*ebiten.Image)

	drawImageOption := ebiten.DrawImageOptions{}
	screen.DrawImage(this.bg_img, &drawImageOption)

	//vector.DrawFilledRect(screen, 0, 0, common.SCREEN_WIDTH, common.SCREEN_HEIGHT, color.RGBA{128, 128, 128, 64}, false)

}

// BG読み込みを実施する
func GoAsyncBgLoad(index int) {
	//すでに読み込み済みなら処理しない
	if this.BgPingImage[index] != nil {
		return
	}
	//非同期で読み込み
	ch := make(chan string)
	go AsyncBgLoading(index, ch)
	close(ch)
}

// BG読み込みを非同期実行するラップ関数
func AsyncBgLoading(index int, ret chan string) {
	readBg(index)
}

// BG読み込み。非同期、同期、どちらからでも使われる
func readBg(index int) {
	this.BgPingImage[index] = loader.GetPngImg("img/bg" + fmt.Sprintf("%d", index) + ".png")
}

// 背景を指定する
func SetBg(index int) {
	this.nowBgIndex = index
}
