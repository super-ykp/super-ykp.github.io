package player

import (
	"Neco/common"
	"Neco/explosion"
	"Neco/shot"
	"Neco/sound"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerStruct struct {
	X      float64
	Y      float64
	Width  float64
	Height float64

	HP    int
	MaxXP int
	HPBar float64

	IsUse bool

	InvincibleCounter int

	animCounter int
}

const (
	HitDot         = 2
	INVICIBLE_TIME = 60 * 3
	speed          = float64(1)

	HpBarPngImgOffsetX = 16
	HpBarPngImgOffsetY = 8
	HpBarWidth         = 32
)

// クラス変数的なもの-------------------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージを扱うのだろう
	m_sound        *sound.SoundStruct
	m_shot         *shot.ShotStruct
	m_explosion    *explosion.ExplosionStruct
	m_Touch        bool
	m_LastTouchX   int
	m_LastTouchY   int
)

// -----------------------------------------------------------------
func Update(p *PlayerStruct) {
	if !p.IsUse {
		return
	}

	//キー移動
	updateKey(p)
	//タッチ移動
	updateTouch(p)

	//ショットの発射
	if common.GameState == common.GameStateNomal {
		shot.Shot(m_shot, p.X, p.Y)
	}

	//無敵時間
	if 0 < p.InvincibleCounter {
		p.InvincibleCounter--
		p.HPBar += (float64(HpBarWidth)*float64(p.HP)/float64(p.MaxXP) - p.HPBar) / 10
	}
	p.animCounter++

}

func updateKey(p *PlayerStruct) {
	//キー入力に応じてプレイヤー位置を移動
	Xvector := float64(0)
	Yvector := float64(0)

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && (p.Width/2) < p.X {
		Xvector -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && (common.SCREEN_WIDTH)-(p.Width/2) > p.X {
		Xvector += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && (p.Height/2) < p.Y {
		Yvector -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && (common.SCREEN_HEIGHT)-(p.Height/2) > p.Y {
		Yvector += speed
	}
	//縦横方向同時押しなら、斜め方向に1になるよう補正
	if Xvector != 0 && Yvector != 0 {
		Xvector *= 0.700
		Yvector *= 0.700
	}

	p.X += Xvector
	p.Y += Yvector
}

// タッチ移動
func updateTouch(p *PlayerStruct) {
	IDs := TouchIDs()
	if len(IDs) > 0 {
		//今までタッチしていなかった
		if !m_Touch {
			m_LastTouchX, m_LastTouchY = GetTouchXY(IDs[0])
		}
		m_Touch = true

		newX, newY := GetTouchXY(IDs[0])
		p.X += float64(newX - m_LastTouchX)
		p.Y += float64(newY - m_LastTouchY)
		m_LastTouchX = newX
		m_LastTouchY = newY

	} else {
		m_Touch = false
	}

	if p.X < (p.Width / 2) {
		p.X = (p.Width / 2)
	}
	if (common.SCREEN_WIDTH)-(p.Width/2) < p.X {
		p.X = (common.SCREEN_WIDTH) - (p.Width / 2)
	}
	if p.Y < (p.Height / 2) {
		p.Y = (p.Height / 2)
	}
	if (common.SCREEN_HEIGHT)-(p.Height/2) < p.Y {
		p.Y = (common.SCREEN_HEIGHT) - (p.Height / 2)
	}
}

func TouchIDs() []ebiten.TouchID {
	//タッチ
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)
	return touchIDs
}
func GetTouchXY(touchID ebiten.TouchID) (x, y int) {
	x, y = ebiten.TouchPosition(touchID)
	return x, y
}

// --======================================================================================
// --======================================================================================
func Draw(screen *ebiten.Image, p *PlayerStruct) {
	if !p.IsUse {
		return
	}

	//ebitenライブラリのDrawImageOptions構造体をインスタンス化
	drawImageOption := &ebiten.DrawImageOptions{}
	//プレイヤーの中心位置に対する描画のずれを考慮し、配置する
	drawImageOption.GeoM.Translate(float64(p.X-p.Width/2), float64(p.Y-p.Height/2))

	//PNG内プレイヤー位置
	imageOffsetX := 0
	imageRect := image.Rect(imageOffsetX, 0, imageOffsetX+int(p.Width), 0+int(p.Height))
	//イメージ取得
	ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
	//描画
	screen.DrawImage(ebitenImage, drawImageOption)

}

func Draw2(screen *ebiten.Image, p *PlayerStruct) {
	if !p.IsUse {
		return
	}
	if 0 < p.InvincibleCounter && INVICIBLE_TIME-60 < p.InvincibleCounter {
		//PNG内バー位置
		imageRect := image.Rect(HpBarPngImgOffsetX, HpBarPngImgOffsetY, HpBarPngImgOffsetX+HpBarWidth, HpBarPngImgOffsetY+4)
		//イメージ取得
		ebitenImage := m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		drawImageOption := &ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(float64(p.X-HpBarWidth/2), float64(p.Y-32))
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)

		//PNG内バー位置
		imageRect = image.Rect(HpBarPngImgOffsetX, HpBarPngImgOffsetY+4, HpBarPngImgOffsetX+int(p.HPBar), HpBarPngImgOffsetY+4+4)
		//イメージ取得
		ebitenImage = m_MainPngImage.SubImage(imageRect).(*ebiten.Image)
		drawImageOption = &ebiten.DrawImageOptions{}
		drawImageOption.GeoM.Translate(float64(p.X-HpBarWidth/2), float64(p.Y-32))
		//描画
		screen.DrawImage(ebitenImage, drawImageOption)

	}
}

// --======================================================================================
// --======================================================================================
func Init(e *ebiten.Image, player *PlayerStruct, shot *shot.ShotStruct, explosion *explosion.ExplosionStruct, sound *sound.SoundStruct) {
	m_sound = sound
	m_shot = shot
	m_explosion = explosion

	m_MainPngImage = e
}
func Start(player *PlayerStruct) {
	player.IsUse = true
	player.Width = 16
	player.Height = 16
	player.X = common.SCREEN_WIDTH / 2
	player.Y = common.SCREEN_HEIGHT - player.Height/2

	player.MaxXP = 2
	player.HP = player.MaxXP
}

// ------------------------------------------------------------------
// 弾に衝突
func HitDan(p *PlayerStruct) {
	if !p.IsUse {
		return
	}
	//接触通知が来ていたら
	//HPがあり無敵じゃなければ
	if 0 < p.HP && p.InvincibleCounter == 0 {
		sound.Stop(m_sound, sound.BOSSHIT, 30)
		sound.Stop(m_sound, sound.KNOKBACK, 30)
		sound.Play(m_sound, sound.DAMAGE)
		//無敵時間設定
		p.InvincibleCounter = INVICIBLE_TIME
		//HPバーの長さを設定
		p.HPBar = float64(HpBarWidth) * float64(p.HP) / float64(p.MaxXP)
		//HP現小
		p.HP--
	} else if p.HP == 0 && p.InvincibleCounter == 0 {
		explosion.Explosion(m_explosion, p.X, p.Y, 1, 6)
		p.IsUse = false
		sound.Stop(m_sound, sound.BOSSHIT, 30)
		sound.Stop(m_sound, sound.KNOKBACK, 30)
		sound.Stop(m_sound, sound.EXPLOSION, 30)
		sound.Play(m_sound, sound.PLAYERDEATH)
		common.StartGameOver()
	}
}

func RrecoveryHp(p *PlayerStruct) {
	p.HP = 2
}
