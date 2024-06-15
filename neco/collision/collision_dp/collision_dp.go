package collision_dp

import (
	"Neco/danmaku"
	"Neco/player"
	"math"
)

const (
	pngOffsetX = 16
	pngOffsetY = 8
	pngWidth   = 32
	pngHighte  = 8
)

// ---------------------------------------
// プレイヤーと弾の当たり判定
func Calc(d *danmaku.DanmakuStruct, p *player.PlayerStruct) {
	hit := bool(false)
	x := p.X
	y := p.Y
	//すべての弾ループ
	for _, dan := range d.DanList {
		if !dan.IsUse {
			continue
		}
		//接触した
		if math.Abs(dan.X-x) < player.HitDot && math.Abs(dan.Y-y) < player.HitDot {
			hit = true
			break
		}
	}
	//接触していた場合、接触を通知
	if hit {
		player.HitDan(p)
	}
}
