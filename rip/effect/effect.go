package effect

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type EcStruct struct {
	x       float64
	y       float64
	v_x     float64
	v_y     float64
	counter int
	types   int
	visible bool
}

const (
	E_P10       = 0
	E_P100      = 1
	E_ONE_UP    = 2
	E_PUSHONEUP = 3
)

type EffectStruct struct {
	ec         [50]EcStruct
	effect_img [5]*ebiten.Image
}

var (
	this *EffectStruct
)

// ==========================================================================================
// ==========================================================================================
func Init(t *EffectStruct, i *ebiten.Image) {
	this = t

	//イメージ
	x := 112
	y := 16
	imageRect := image.Rect(x, y, x+16, y+16)
	t.effect_img[E_P10] = i.SubImage(imageRect).(*ebiten.Image)
	x = 112 + 16
	y = 16
	imageRect = image.Rect(x, y, x+32, y+16)
	t.effect_img[E_P100] = i.SubImage(imageRect).(*ebiten.Image)
	x = 112 + 16 + 32
	y = 16
	imageRect = image.Rect(x, y, x+32, y+16)
	t.effect_img[E_ONE_UP] = i.SubImage(imageRect).(*ebiten.Image)

	x = 192
	y = 16
	imageRect = image.Rect(x, y, x+16, y+16)
	t.effect_img[E_PUSHONEUP] = i.SubImage(imageRect).(*ebiten.Image)

}

func Update() {
	for i := 0; i < len(this.ec); i++ {
		f := &this.ec[i]
		if !f.visible {
			continue
		}
		f.counter--
		f.x += f.v_x
		f.y += f.v_y
		if f.counter < 0 {
			f.visible = false
		}
	}
}

func Draw(screen *ebiten.Image) {
	for i := 0; i < len(this.ec); i++ {
		f := &this.ec[i]
		if !f.visible {
			continue
		}
		p := f.types
		x := float64(0)
		if p == E_ONE_UP || p == E_P100 {
			x = f.x - float64(32/2)
		} else {
			x = f.x - float64(16/2)
		}
		y := f.y - 16/2
		drawImageOption := ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(x, y)
		screen.DrawImage(this.effect_img[p], &drawImageOption)
	}
}

// 空きバッファの返却
func GetBuff(x, y float64, types int) *EcStruct {
	for i := 0; i < len(this.ec); i++ {
		f := &this.ec[i]
		if f.visible {
			continue
		}
		f.visible = true
		f.x = x
		f.y = y
		f.types = types
		return f
	}
	return nil
}

// エフェクトスタート
func Start(x, y float64, types int) {
	f := GetBuff(x, y, types)
	if f == nil {
		return
	}
	f.counter = 60
	f.v_x = 0.5 + rand.Float64()*1
	f.v_y = -1.5
	if int(rand.Float64()*2) == 0 {
		f.v_x *= -1
	}
}

// 1Up出現エフェクト
func GetOneUpEffect(x, y float64) {
	f := GetBuff(x, y, E_ONE_UP)
	if f == nil {
		return
	}
	f.counter = 60 * 2
	f.v_y = -0.5
	f.v_x = 0
}

// 1Up取得エフェクト
func PushOneUpEffect(x, y float64) {
	for r := 0; r < 360; r += 20 {
		f := GetBuff(x, y, E_PUSHONEUP)
		if f == nil {
			return
		}
		f.counter = 60
		f.v_x = 1
		f.v_y = 0
		f.v_x, f.v_y = rotateVector(f.v_x, f.v_y, float64(r))
		r += 20
		f.types = E_PUSHONEUP

	}
}

// ベクトル計算
func rotateVector(x, y, angle float64) (float64, float64) {
	rad := angle * (math.Pi / 180.0)
	newX := x*math.Cos(rad) - y*math.Sin(rad)
	newY := x*math.Sin(rad) + y*math.Cos(rad)
	return newX, newY
}

// 全部クリア
func AllClear() {
	for i := 0; i < len(this.ec); i++ {
		this.ec[i].visible = false
	}
}
