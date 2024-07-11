package item

import (
	"Yoko/camera"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pNG_OFFSET_X = 288
	pNG_OFFSET_Y = 16
	GRAVITY      = 0.1

	M_FREE       = 1
	M_COLLENCT   = 2
	M_KETTEIMACH = 3

	I_EXP   = 1
	I_COIN  = 2
	I_FORCE = 3

	I_MAX = 4
)

type ItStruct struct {
	IsUse    bool
	Visible  bool
	ItemType int //アイテムタイプ

	MoveMode    int
	MoveCounter int

	X  float64
	Y  float64
	Z  int
	vx float64
	vy float64

	animIndex    uint
	animCounter  uint
	dispPngIndex uint
}

type ItemStruct struct {
	ItemImage []*ebiten.Image
	ItemList  []ItStruct
}

var (
	this  *ItemStruct
	animC [][][]uint
)

func Reset() {
	this.ItemList = []ItStruct{}
}

// --===================================================================================
// 初期化
func Init(t *ItemStruct, e *ebiten.Image) {
	this = t

	offs := [][]int{
		{0, 0, 16, 16},   //0 EXP
		{0, 16, 16, 16},  //1
		{0, 32, 16, 16},  //2
		{16, 0, 16, 16},  //3 COIN
		{16, 16, 16, 16}, //4
		{16, 32, 16, 16}, //5
		{32, 0, 16, 16},  //6 FC
		{32, 16, 16, 16}, //7
	}

	for i := 0; i < len(offs); i++ {
		offX := pNG_OFFSET_X + offs[i][0]
		offy := pNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]

		imageRect := image.Rect(offX, offy, offW, offH)
		rect := e.SubImage(imageRect).(*ebiten.Image)
		this.ItemImage = append(this.ItemImage, rect)
	}
	animC = make([][][]uint, I_MAX)
	animC[I_EXP] = [][]uint{{0, 3}, {1, 3}, {2, 3}}
	animC[I_COIN] = [][]uint{{3, 3}, {4, 3}, {5, 3}, {4, 3}}
	animC[I_FORCE] = [][]uint{{6, 1}, {7, 1}}
}

func Update() {
	//アニメーション制御---------------------------------------------------------------------
	for i := 0; i < len(this.ItemList); i++ {
		//バッファが利用中でないなら次
		ec := &this.ItemList[i]
		if !ec.IsUse {
			continue
		}
		//自由ベクトル移動モード
		switch ec.MoveMode {
		case M_FREE:
			//ec.vy += GRAVITY
			ec.X += ec.vx
			ec.Y += ec.vy
			ec.vx *= 0.90
			ec.vy *= 0.90
			if 0 < ec.Y {
				ec.Y = 0
				ec.vy = -ec.vy / 2
			}
			ec.MoveCounter--
			if ec.MoveCounter < 0 {
				ec.MoveMode = M_COLLENCT
			}
		case M_COLLENCT: //回収モード

			xv := float64(camera.X) - (ec.X)
			yv := -16 - (ec.Y)
			v := float64(math.Sqrt(xv*xv + yv*yv))

			//ベクトル設定
			sp := float64(15)
			if distance(ec.X, ec.Y, float64(camera.X), -16) > 8 {
				sp = 15
			} else {
				sp = 5
			}
			ec.vx = (xv / v) * sp
			ec.vy = (yv / v) * sp
			ec.X += ec.vx
			ec.Y += ec.vy

			if distance(ec.X, ec.Y, float64(camera.X), -16) < 4 {
				ec.MoveMode = M_KETTEIMACH
				ec.Visible = false
			}

		}

		countMax := animC[ec.ItemType][ec.animIndex][1] % 100
		ec.animCounter++
		//アニメーション表示フレームを超えたら
		if ec.animCounter > countMax && countMax != 0 {
			//次のアニメパターンへ
			ec.animCounter = 0
			ec.animIndex++

			//アニメーションが最初まで行ったら0に
			if ec.animIndex == uint(len(animC[ec.ItemType])) {
				ec.animIndex = 0
			}
		}
		//元アニメのPng番号を指定
		ec.dispPngIndex = animC[ec.ItemType][ec.animIndex][0]
	}
}

// 描画
func Draw(screen *ebiten.Image) {
	for i := 0; i < len(this.ItemList); i++ {
		//バッファが利用中でないなら次
		ec := this.ItemList[i]
		if !ec.IsUse || !ec.Visible {
			continue
		}

		//アイテム表示
		drawImageOption := ebiten.DrawImageOptions{}
		ofX := int(ec.X-(32)/2) - camera.CamOffsetX
		ofY := ec.Z + int(ec.Y-32/2) - camera.CamOffsetY
		drawImageOption.GeoM.Translate(float64(ofX), float64(ofY))
		screen.DrawImage(this.ItemImage[ec.dispPngIndex], &drawImageOption)
	}
}

// 外部からアイテムをセットする
func SetItem(itemType int, x, y, vx, vy float64) {
	ec := getBuffer()
	ec.IsUse = true
	ec.Visible = true
	ec.ItemType = itemType
	ec.animIndex = 0
	ec.X = x
	ec.Y = y
	ec.Z = 0
	ec.MoveMode = M_FREE
	ec.MoveCounter = 30 + rand.Intn(10)
	ec.vx = vx
	ec.vy = vy
}

// バッファの空きを探し返却、なければ追加する
func getBuffer() *ItStruct {
	for i := 0; i < len(this.ItemList); i++ {
		ec := &this.ItemList[i]
		if ec.IsUse {
			continue
		}
		return ec
	}
	ep := ItStruct{}
	this.ItemList = append(this.ItemList, ep)
	return &this.ItemList[len(this.ItemList)-1]
}

// 指定座標間の直選距離を求める------------------------------------
func distance(x, y, nx, ny float64) float64 {
	dx := nx - x
	dy := ny - y
	return math.Sqrt(dx*dx + dy*dy)
}

func GetThis() *ItemStruct {
	return this
}
