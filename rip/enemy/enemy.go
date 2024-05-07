package enemy

import (
	"Rip/common"
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnStruct struct {
	X       float64
	Y       float64
	v_x     float64
	counter int
	Visible bool
}
type EnemyStruct struct {
	ec        [100]EnStruct
	enemy_img [2]*ebiten.Image
}

var (
	this *EnemyStruct
)

// ==========================================================================================
// ==========================================================================================
func Init(t *EnemyStruct, i *ebiten.Image) {
	this = t
	//イメージ
	imageRect := image.Rect(64, 16, 64+24, 16+16)
	t.enemy_img[0] = i.SubImage(imageRect).(*ebiten.Image)
	imageRect = image.Rect(88, 16, 88+24, 16+16)
	t.enemy_img[1] = i.SubImage(imageRect).(*ebiten.Image)
}

func Update() {
	for i := 0; i < len(this.ec); i++ {
		e := &this.ec[i]
		if !e.Visible {
			continue
		}
		e.counter += 1
		e.X += e.v_x
		if e.X < -64 {
			e.X = common.SCREEN_WIDTH + 64
		} else if common.SCREEN_WIDTH+64 < e.X {
			e.X = -64
		}
	}
}

func Draw(screen *ebiten.Image) {
	for i := 0; i < len(this.ec); i++ {
		e := &this.ec[i]
		if !e.Visible {
			continue
		}

		p := (e.counter / 5) % 2
		x := e.X - 24/2
		y := e.Y - 16/2

		drawImageOption := ebiten.DrawImageOptions{}
		if e.v_x > 0 {
			drawImageOption.GeoM.Scale(-1.0, 1.0)
			x += 24
		} else {
			drawImageOption.GeoM.Scale(1.0, 1.0)
		}
		drawImageOption.GeoM.Translate(x, y)

		screen.DrawImage(this.enemy_img[p], &drawImageOption)
	}
}

func Start(max int) {
	for i := 0; i < 100; i++ {
		e := &this.ec[i]
		if i <= max {
			e.v_x = rand.Float64()*0.5 + 0.25
			if int(rand.Float64()*2) == 0 {
				e.v_x *= -1
			}
			e.X = rand.Float64() * common.SCREEN_WIDTH
			e.Y = rand.Float64()*(common.SCREEN_HEIGHT-100) + 40
			e.Visible = true
		} else {
			e.Visible = false
		}
	}
}

func GetEnemyArray() *[100]EnStruct {
	return &this.ec
}
