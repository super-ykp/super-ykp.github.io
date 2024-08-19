package bg

import (
	"Yoko/camera"
	"Yoko/common"
	"Yoko/fileloader"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PNG_OFFSET_X = 0
	PNG_OFFSET_Y = 0
	PNG_WIDTH    = 64
	PNG_HIGTHE   = 32
)

type BgStruct struct {
	bg_img      *ebiten.Image
	BgPingImage []*ebiten.Image
	nowBgIndex  int

	groundObject *ebiten.Image
}

var (
	this *BgStruct
)

// --===================================================================================
// 　背景
// --===================================================================================
// 初期化
func Init(t *BgStruct, p *ebiten.Image) {
	this = t
	this.nowBgIndex = 0
	//最初のBGは、同期で読み込む
	readBg(0)

	imageRect := image.Rect(PNG_OFFSET_X, PNG_OFFSET_Y, PNG_OFFSET_X+PNG_WIDTH, PNG_OFFSET_Y+PNG_HIGTHE)
	this.groundObject = p.SubImage(imageRect).(*ebiten.Image)

}

// 計算
func Update() {
}

// 描画
func Draw(screen *ebiten.Image) {
	//対象の背景が読込中なら
	if len(this.BgPingImage) <= this.nowBgIndex || this.BgPingImage[this.nowBgIndex] == nil {
		//何もしない
		return
	} else {//読み込み終わっていたら、表示
		imageRect := image.Rect(0, 0, common.SCREEN_WIDTH, common.SCREEN_HEIGHT)
		this.bg_img = this.BgPingImage[this.nowBgIndex].SubImage(imageRect).(*ebiten.Image)

		drawImageOption := ebiten.DrawImageOptions{}
		screen.DrawImage(this.bg_img, &drawImageOption)
	}
	//地面、カメラ位置に応じて生成される
	for i := 0; i <= common.SCREEN_WIDTH/PNG_WIDTH+2; i++ {
		drawImageOption := ebiten.DrawImageOptions{}
		ofX := float64(int(PNG_WIDTH-camera.CamOffsetX)%PNG_WIDTH) - PNG_WIDTH
		drawImageOption.GeoM.Translate(ofX+float64(i*PNG_WIDTH), (-16/2)-float64(camera.CamOffsetY))
		screen.DrawImage(this.groundObject, &drawImageOption)
	}
}

// BG読み込みを実施する
func GoAsyncBgLoad(index int) {
	//すでに読み込み済みなら処理しない
	if len(this.BgPingImage) > index && this.BgPingImage[index] != nil {
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
	bg := fileloader.GetPngImg("img/bg" + fmt.Sprintf("%d", index) + ".png")
	this.BgPingImage = append(this.BgPingImage, bg)
}

// 背景を指定する
func SetBg(index int) {
	this.nowBgIndex = index
}
