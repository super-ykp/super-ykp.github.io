package collision_se

import (
	"Neco/enemy"
	"Neco/shot"
	"math"
)

const (
	pngOffsetX = 16
	pngOffsetY = 8
	pngWidth   = 32
	pngHighte  = 8
)

// ---------------------------------------
// ショットと敵の当たり判定
func Calc(s *shot.ShotStruct, e *enemy.EnemyStruct) {
	//ショット数ループ
	for si := 0; si < len(s.ShotList); si++ {
		//ショット未使用
		if !s.ShotList[si].IsUse {
			continue
		}
		//ショット座標取り出し
		x := s.ShotList[si].X
		y := s.ShotList[si].Y

		//すべての敵ループ
		for ei := 0; ei < len(e.EnemyList); ei++ {
			//敵未使用
			if !e.EnemyList[ei].IsUse {
				continue
			}
			//接触した
			if math.Abs(e.EnemyList[ei].X-x) < e.EnemyList[ei].EnemyWidth/2 && math.Abs(e.EnemyList[ei].Y-y) < e.EnemyList[ei].EnemyHeight/2 {
				enemy.Hit(&e.EnemyList[ei])
				s.ShotList[si].IsUse = false
			}
		}
	}
}
