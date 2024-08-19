package item

import (
	"Yoko/camera"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	//main.png内でのイメージ位置
	pNG_OFFSET_X = 288
	pNG_OFFSET_Y = 16

	//アイテム種類
	I_EXP   = 1 //アイテム種類、EXP
	I_COIN  = 2 //アイテム種類、COIN
	I_FORCE = 3 //アイテム種類、FORCE
	I_MAX   = 4 //アイテム最大数

	//アイテム状態
	M_FREE        = 1 //ベクトルに応じて自由移動
	M_COLLENCT    = 2 //プレイヤーに集まる
	M_KAISYUMACHI = 3 //上位存在による回収待ち
)

type ItStruct struct {
	IsUse    bool //利用中フラグ
	Visible  bool //表示フラグ
	ItemType int  //アイテムタイプ

	MoveMode    int //動作モード
	MoveCounter int //その動作の制御カウンタ

	X  float64 //座標
	Y  float64
	Z  int
	vx float64 //運動ベクトル
	vy float64

	animIndex    uint //アニメ種類
	animCounter  uint //アニメカウンタ
	dispPngIndex uint //表示するpngのインデックス
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
	//main.pngからイメージを切り出す
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

	//それを構造体に突っ込む
	for i := 0; i < len(offs); i++ {
		offX := pNG_OFFSET_X + offs[i][0]
		offy := pNG_OFFSET_Y + offs[i][1]
		offW := offX + offs[i][2]
		offH := offy + offs[i][3]

		imageRect := image.Rect(offX, offy, offW, offH)
		rect := e.SubImage(imageRect).(*ebiten.Image)
		this.ItemImage = append(this.ItemImage, rect)
	}
	//アニメーション事にタイムシートを作る
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

		switch ec.MoveMode {
		case M_FREE: //自由ベクトル移動モード
			//ベクトルに応じて移動
			ec.X += ec.vx
			ec.Y += ec.vy
			//ベクトルは空気抵抗を受け、急速に遅くなる
			ec.vx *= 0.90
			ec.vy *= 0.90
			//地面にあたった場合跳ね返る
			if 0 < ec.Y {
				ec.Y = 0
				ec.vy = -ec.vy / 2
			}
			//カウンタ加算
			ec.MoveCounter--
			//一定時間が立っていたら回収モードへ
			if ec.MoveCounter < 0 {
				ec.MoveMode = M_COLLENCT
			}
		case M_COLLENCT: //回収モード
			//プレイヤーのほうへ向かっていく
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

			//プレイヤーと一定距離内になったら非表示に。あとはマネージャが回収してくれる
			if distance(ec.X, ec.Y, float64(camera.X), -16) < 4 {
				ec.MoveMode = M_KAISYUMACHI
				ec.Visible = false
			}

		}
		//下二桁がアニメ時間
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
