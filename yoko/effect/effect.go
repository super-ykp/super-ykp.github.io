package effect

import (
	"Yoko/camera"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pNG_OFFSET_X = 0
	pNG_OFFSET_Y = 128

	E_EXPLOTION = 0//爆発
	E_HITEFFEC  = 1//ヒット
	E_DAMAGE    = 2//ダメージ
	E_MAX       = 3 //エフェクト最大数
)

type EfStruct struct {
	IsUse         bool
	EffectType    int     //エフェクトタイプ
	Magnification float64 //拡大率
	Rotate        float64 //回転角
	X             float64 //座標X
	Y             float64 //座標Y
	Z             int     //座標Z（仮想奥行き）
	vx            float64 //ベクトルx
	vy            float64 //ベクトルy
	dispPngIndex  int     //表示PNGオフセット
	AnimCounter   int     //アニメーションカウンタ
	animIndex     int     //再生中アニメ種類インデックス
	animSpeed     int     //アニメーション速度
}

type EffectStruct struct {
	effectImage []*ebiten.Image
	EList       []EfStruct
}

var (
	this  *EffectStruct
	animC [][]int
)

// --===================================================================================
// 初期化
func Init(t *EffectStruct, p *ebiten.Image) {
	//インスタンスの取り込み
	this = t
	//main.pngの中のどの場所の絵か
	offs := [][]int{
		{0, 0, 48, 48}, //0EXPROSION
		{48, 0, 48, 48},
		{96, 0, 48, 48},
		{144, 0, 48, 48},
		{192, 0, 48, 48},
		{240, 0, 48, 48},
		{288, 0, 48, 48},
		{336, 0, 48, 48},
		{336, 0, 48, 48},
		{0, 48, 64, 16}, //9Hit
		{64, 48, 6, 6},  //10DMG
	}
	//上記を構造体へ
	for i := 0; i < len(offs); i++ {
		offX := pNG_OFFSET_X + offs[i][0]
		offy := pNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]
		//イメージも切り出して保持
		imageRect := image.Rect(offX, offy, offW, offH)
		rect := p.SubImage(imageRect).(*ebiten.Image)
		this.effectImage = append(this.effectImage, rect)
	}
	//アニメ種類に応じたキャラパターン
	animC = make([][]int, E_MAX)
	animC[E_EXPLOTION] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}//爆発
	animC[E_HITEFFEC] = []int{9}//ヒットエフェクト
	animC[E_DAMAGE] = []int{10}//ダメージ
}

func Update() {
	//バッファ全周
	for i := 0; i < len(this.EList); i++ {
		//未使用ならスキップ
		ec := &this.EList[i]
		if !ec.IsUse {
			continue
		}

		//エフェクトタイプ毎の固有処理
		switch ec.EffectType {
		case E_HITEFFEC://ヒットエフェクト
			//少しずつ小さくなる
			ec.Magnification -= 0.1
		}
		//ベクトル加算
		ec.X += ec.vx
		ec.Y += ec.vy

		//アニメーションフレーム。下位2桁ががフレーム
		ec.AnimCounter++

		//アニメーション表示フレームを超えたら
		if ec.AnimCounter > ec.animSpeed {
			//次のアニメパターンへ
			ec.AnimCounter = 0
			//アニメーションが最初まで行ったら終わり
			ec.animIndex++
			if ec.animIndex == len(animC[ec.EffectType]) {
				ec.IsUse = false
				ec.animIndex = 0
				continue
			}
		}
		//元アニメのPng番号を指定
		ec.dispPngIndex = animC[ec.EffectType][ec.animIndex]
	}
}

// 描画
func Draw(screen *ebiten.Image) {
	for i := 0; i < len(this.EList); i++ {
		//バッファが利用中でないなら次
		ec := &this.EList[i]
		if !ec.IsUse {
			continue
		}
		//爆発フェクトとダメージエフェクトは、これといって特殊なことをしない
		if ec.EffectType == E_EXPLOTION || ec.EffectType == E_DAMAGE {
			//エフェクト表示
			drawImageOption := ebiten.DrawImageOptions{}
			r := this.effectImage[ec.dispPngIndex].Bounds().Size()
			widthX := float64(r.X)
			widthY := float64(r.Y)
			drawImageOption.GeoM.Translate(-widthX/2, -widthY/2)
			drawImageOption.GeoM.Scale(ec.Magnification, ec.Magnification)
			drawImageOption.GeoM.Translate(widthX/2, widthY/2)

			ofX := int(ec.X-(widthX)/2) - camera.CamOffsetX
			ofY := int(ec.Y-(widthY)/2) + ec.Z - camera.CamOffsetY
			drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))

			screen.DrawImage(this.effectImage[ec.dispPngIndex], &drawImageOption)
		} else if ec.EffectType == E_HITEFFEC {//ヒットエフェクトは半透明で縮小していく
			//エフェクト表示
			drawImageOption := ebiten.DrawImageOptions{}
			r := this.effectImage[ec.dispPngIndex].Bounds().Size()
			widthX := float64(r.X)
			widthY := float64(r.Y)
			drawImageOption.GeoM.Translate(-widthX/2, -widthY/2)
			drawImageOption.GeoM.Scale(4, ec.Magnification)

			drawImageOption.GeoM.Rotate(ec.Rotate)
			drawImageOption.ColorScale.Scale(0.5, 0.5, 0.5+rand.Float32()+0.5, 0.5)

			drawImageOption.GeoM.Translate(widthX/2, widthY/2)

			ofX := int(ec.X-(widthX)/2) - camera.CamOffsetX
			ofY := int(ec.Y-(widthY)/2) + ec.Z - camera.CamOffsetY
			drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))

			screen.DrawImage(this.effectImage[ec.dispPngIndex], &drawImageOption)
		}
	}

}

// エフェクトを配置する
//x, y, z, 発生座標
//etype int, エフェクトタイプ
//sub1 float64 エフェクト大きさ
func SetEffect(x, y, z, etype int, sub1 float64) {
	ec := getBuffer()
	ec.EffectType = etype
	ec.IsUse = true
	ec.X = float64(x)
	ec.Y = float64(y)
	ec.Z = z
	ec.vx = 0
	ec.vy = 0
	ec.Magnification = 1
	ec.Rotate = 0

	switch etype {
	case E_EXPLOTION://爆発
		ec.Magnification = sub1//爆発の大きさ
		ec.animSpeed = int(sub1 * 3)//爆発が大きいほど爆発もゆっくり
	case E_HITEFFEC://ヒットエフェクト
		ec.Magnification = 0.5
		ec.animSpeed = 3
		ec.Rotate = rand.Float64() * 360//ランダムで回転
	case E_DAMAGE://ダメージ
		ec.animSpeed = 10 + rand.Intn(30)//#模様にベクトルを設定
		ec.vx = -1 - rand.Float64()*2
		ec.vy = rand.Float64()*2 - 1
	}
	//他基本初期化
	ec.animIndex = 0
	ec.AnimCounter = 0
	ec.dispPngIndex = animC[ec.EffectType][ec.animIndex]
}

// バッファの空きを探し返却、なければ追加する
func getBuffer() *EfStruct {
	for i := 0; i < len(this.EList); i++ {
		ec := &this.EList[i]
		if ec.IsUse {
			continue
		}
		return ec
	}
	ep := EfStruct{}
	this.EList = append(this.EList, ep)
	return &this.EList[len(this.EList)-1]
}
