package danmaku

import (
	"Neco/common"
	"Neco/player"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	pngOffsetX = 16
	pngOffsetY = 0
)

type DanStruct struct {
	IsUse   bool
	X       float64
	Y       float64
	VecterX float64
	VecterY float64
	Axel    float64
	Brake   float64

	Width       float64
	Height      float64
	animCounter int
}

type DanmakuStruct struct {
	DanList [1024 * 12]DanStruct
	TatgetX float64
	TatgetY float64
}

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage      *ebiten.Image //イメージ
	m_RironCounter      int           //理論カウンタ
	m_LastRecodeCounter int
	m_PlayerX           float64
	m_PlayerY           float64
	m_PlayerWidth       float64
	m_PlayerHighte      float64
)

// -----------------------------------------------------------------
// 計算
func Update(d *DanmakuStruct, p *player.PlayerStruct) {
	m_PlayerX = p.X
	m_PlayerY = p.Y

	//弾幕配列は固定長
	for i := 0; i < len(d.DanList); i++ {
		//未使用なら飛ばす
		if !d.DanList[i].IsUse {
			continue
		}
		t := &d.DanList[i]
		//利用中ならベクトル加算
		t.X += t.VecterX * t.Axel
		t.Y += t.VecterY * t.Axel
		if t.Axel > 1 {
			t.Axel -= t.Brake
			if t.Axel < 1 {
				t.Axel = 1
			}
		}
		t.animCounter++

		//画面外に行ったら消す
		if t.X < 0-10 || common.SCREEN_WIDTH+10 < t.X || t.Y < -10 || common.SCREEN_HEIGHT+10 < t.Y {
			t.IsUse = false
		}
	}
	//理論カウンタ加算
	m_RironCounter++
	if m_RironCounter >= math.MaxInt {
		m_RironCounter = 0
		m_LastRecodeCounter = -100
	}
}

// -----------------------------------------------------------------
// 描画
func Draw(screen *ebiten.Image, d *DanmakuStruct) {
	//弾幕配列は固定長
	for i := 0; i < len(d.DanList); i++ {
		//未使用なら飛ばす
		if !d.DanList[i].IsUse {
			continue
		}
		//ebitenライブラリのDrawImageOptions構造体をインスタンス化
		drawImageOption := &ebiten.DrawImageOptions{}
		//中心位置に対する描画のずれを考慮し、配置する
		drawImageOption.GeoM.Translate(d.DanList[i].X-d.DanList[i].Width/2, d.DanList[i].Y-d.DanList[i].Height/2)

		//PNG内位置
		OffsetX := pngOffsetX + ((d.DanList[i].animCounter/3)%4)*int(d.DanList[i].Width)
		imageRect := image.Rect(OffsetX, 0, OffsetX+int(d.DanList[i].Width), pngOffsetY+int(d.DanList[i].Height))
		//イメージ取得
		ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)
	}
}

// 初期化------------------------------------------------------------------
func Init(e *ebiten.Image, d *DanmakuStruct, p *player.PlayerStruct) {
	m_MainPngImage = e
	m_LastRecodeCounter = -100

	m_PlayerWidth = p.Width
	m_PlayerHighte = p.Height
}
func Start(d *DanmakuStruct) {
	for i := 0; i < len(d.DanList); i++ {
		d.DanList[i].IsUse = false
	}
}

// 発射----------------------------------------------------------------
func ShotZ(d *DanmakuStruct, x, y, speed, rotationAngle float64) {
	if (m_PlayerX < m_PlayerWidth && x < m_PlayerWidth*2) || (common.SCREEN_WIDTH-m_PlayerWidth < m_PlayerX && common.SCREEN_WIDTH-m_PlayerWidth*2 < x) {
		return
	}

	//弾幕バッファの空き取得
	t := getDanmakuBuf(d)
	//もう無いなら無理
	if t == nil {
		return
	}

	//初期配置
	t.IsUse = true
	t.Width = 5
	t.Height = 5
	t.X = x
	t.Y = y
	t.Axel = 1
	t.Brake = 0

	//自機方向ベクトルを計算
	interval := 30
	if (m_PlayerX < m_PlayerWidth) || (common.SCREEN_WIDTH-m_PlayerWidth < m_PlayerX) {
		interval = 300
	}
	if m_LastRecodeCounter+interval < m_RironCounter {
		m_LastRecodeCounter = m_RironCounter
		d.TatgetX = m_PlayerX
		d.TatgetY = m_PlayerY
	}

	xv := d.TatgetX - x
	yv := d.TatgetY - y
	v := float64(math.Sqrt(xv*xv + yv*yv))

	//ズレの指定があった場合回転
	if rotationAngle != 0 {
		rx, ry := rotateVector(xv, yv, rotationAngle)
		xv = rx
		yv = ry
	}
	//ベクトル設定
	t.VecterX = (xv / v) * speed
	t.VecterY = (yv / v) * speed
}

// 発射。ボス用----------------------------------------------------------------
func ShotB(d *DanmakuStruct, x, y, speed, rotationAngle float64, isNoAim bool, axcel, brake float64) {
	//弾幕バッファの空き取得
	t := getDanmakuBuf(d)
	//もう無いなら無理
	if t == nil {
		return
	}

	//初期配置
	t.IsUse = true
	t.Width = 5
	t.Height = 5
	t.X = x
	t.Y = y
	t.Axel = axcel
	t.Brake = brake
	//自機方向ベクトルを計算
	xv := float64(0)
	yv := float64(1)
	v := float64(0)

	//自機を狙うか
	if isNoAim { //狙わない
		xv = float64(0)
		yv = float64(1)
		v = float64(math.Sqrt(xv*xv + yv*yv))

	} else {
		xv = m_PlayerX - x
		yv = m_PlayerY - y
		v = float64(math.Sqrt(xv*xv + yv*yv))

	}

	//ズレの指定があった場合回転
	if rotationAngle != 0 {
		rx, ry := rotateVector(xv, yv, rotationAngle)
		xv = rx
		yv = ry
	}
	//ベクトル設定
	t.VecterX = (xv / v) * speed
	t.VecterY = (yv / v) * speed
}

// 弾幕空きバッファの返却---------------------------
func getDanmakuBuf(d *DanmakuStruct) *DanStruct {
	i := 0
	//開いているバッファを探す
	for ; i < len(d.DanList); i++ {
		if !d.DanList[i].IsUse {
			break
		}
	}
	//もう無いなら無理
	if i >= len(d.DanList) {
		return nil
	}
	return &d.DanList[i]
}

// ベクトルを回転させる----------------------------
func rotateVector(x, y, angle float64) (float64, float64) {
	rad := angle * (math.Pi / 180.0)
	newX := x*math.Cos(rad) - y*math.Sin(rad)
	newY := x*math.Sin(rad) + y*math.Cos(rad)
	return newX, newY
}
