package shot

import (
	"Neco/common"
	"Neco/sound"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pngOffsetX = 16 * 3
	pngOffsetY = 0
	ShotWidth  = 16
	ShotHeight = 16
	ShotSpeed  = 6
)

type TamaStruct struct {
	IsUse   bool
	X       float64
	Y       float64
	VecterX float64
	VecterY float64
	Pattern int
}

type ShotStruct struct {
	ShotList [1024]TamaStruct
}

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージ

	m_Sound *sound.SoundStruct

	m_Counter int
)

// -----------------------------------------------------------------
// 計算
func Update(d *ShotStruct) {
	//配列は固定長
	for i := 0; i < len(d.ShotList); i++ {
		//未使用なら飛ばす
		if !d.ShotList[i].IsUse {
			continue
		}
		//利用中ならベクトル加算
		d.ShotList[i].X += d.ShotList[i].VecterX
		d.ShotList[i].Y += d.ShotList[i].VecterY

		//画面外に行ったら消す
		if d.ShotList[i].X < 0-ShotWidth || common.SCREEN_WIDTH+ShotWidth < d.ShotList[i].X || d.ShotList[i].Y < -ShotHeight || common.SCREEN_HEIGHT+ShotHeight < d.ShotList[i].Y {
			d.ShotList[i].IsUse = false
		}
	}
	m_Counter++
}

// -----------------------------------------------------------------
// 描画
func Draw(screen *ebiten.Image, d *ShotStruct) {
	//配列は固定長
	for i := 0; i < len(d.ShotList); i++ {
		//未使用なら飛ばす
		if !d.ShotList[i].IsUse {
			continue
		}
		//ebitenライブラリのDrawImageOptions構造体をインスタンス化
		drawImageOption := &ebiten.DrawImageOptions{}
		//中心位置に対する描画のずれを考慮し、配置する
		drawImageOption.GeoM.Translate(d.ShotList[i].X-ShotWidth/2, d.ShotList[i].Y-ShotHeight/2)

		//PNG内位置
		OffsetX := pngOffsetX + d.ShotList[i].Pattern*16
		imageRect := image.Rect(OffsetX, 0, OffsetX+int(ShotWidth), pngOffsetY+int(ShotHeight))
		//イメージ取得
		ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)
	}
}

// 初期化------------------------------------------------------------------
func Init(e *ebiten.Image, d *ShotStruct, s *sound.SoundStruct) {
	m_MainPngImage = e
	m_Sound = s
}

func Start(d *ShotStruct) {
	for i := 0; i < len(d.ShotList); i++ {
		d.ShotList[i].IsUse = false
	}
}

// 発射----------------------------------------------------------------
func Shot(d *ShotStruct, x, y float64) {
	if (m_Counter % 4) != 0 {
		return
	}
	sound.Play(m_Sound, sound.SHOT)

	var rotats []float64 = []float64{0, 10, 20, 30, 40, 50}
	for i := -1; i < 2; i += 2 {
		for r := 0; r < 6; r++ {
			//バッファの空き取得
			t := getTamaBuf(d)
			//もう無いなら無理
			if t == nil {
				return
			}

			//初期配置
			t.IsUse = true
			t.X = x + float64(6*i+(m_Counter%3)) - 1
			t.Y = y
			t.Pattern = 5 + r*i
			rotationAngle := float64(rotats[r]) * float64(i)
			xv, yv := rotateVector(0, -1, rotationAngle)

			//ベクトル設定
			t.VecterX = (xv / 1) * ShotSpeed
			t.VecterY = (yv / 1) * ShotSpeed
		}
	}
}

// 空きバッファの返却---------------------------
func getTamaBuf(d *ShotStruct) *TamaStruct {
	i := 0
	//開いているバッファを探す
	for ; i < len(d.ShotList); i++ {
		if !d.ShotList[i].IsUse {
			return &d.ShotList[i]
		}
	}
	return nil
}

// ベクトルを回転させる----------------------------
func rotateVector(x, y, angle float64) (float64, float64) {
	rad := angle * (math.Pi / 180.0)
	newX := x*math.Cos(rad) - y*math.Sin(rad)
	newY := x*math.Sin(rad) + y*math.Cos(rad)
	return newX, newY
}
