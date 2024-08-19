package camera

import (
	"Yoko/common"
	"math/rand"
)

var (
	//カメラ座標
	X int
	Y int
	//カメラが実際に映すズレ
	CamOffsetX int
	CamOffsetY int
	//振動
	Vibration float64
)

func Update() {
	//カメラは基本カメラ座標に追従する
	CamOffsetX = X - common.SCREEN_WIDTH/2
	CamOffsetY = Y + common.CAM_Y_OFFSET
	//振動エフェクト発生中は、オフセットにランダム値が加算される
	if Vibration > 0 {
		CamOffsetX += int(rand.Float64()*Vibration*2 - Vibration)
		CamOffsetY += int(rand.Float64()*Vibration*2 - Vibration)
		Vibration -= 0.3
	}

}
//カメラ振動
func SetViblation(v float64) {
	Vibration = v
}
