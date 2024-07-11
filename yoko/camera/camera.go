package camera

import (
	"Yoko/common"
	"math/rand"
)

var (
	X int
	Y int

	CamOffsetX int
	CamOffsetY int

	Vibration float64
)

func Update() {
	CamOffsetX = X - common.SCREEN_WIDTH/2
	CamOffsetY = Y + common.CAM_Y_OFFSET
	if Vibration > 0 {
		CamOffsetX += int(rand.Float64()*Vibration*2 - Vibration)
		CamOffsetY += int(rand.Float64()*Vibration*2 - Vibration)
		Vibration -= 0.3
	}

}

func SetViblation(v float64) {
	Vibration = v
}
