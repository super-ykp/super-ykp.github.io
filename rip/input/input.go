package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputStruct struct {
	IsTouch     bool
	firstTouchX int
	firstTouchY int
}

var (
	this *InputStruct
)

// ==========================================================================================
// 操作入力
// ==========================================================================================
// 初期化
func Init(t *InputStruct) {
	this = t
}

// 左右移動。neco =false ひこーき移動　=true ねこ移動
func GetPressLR(neco bool) float64 {
	//カーソルキーで移動.................
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		return -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		return 1
	}

	//タッチ移動................
	IDs := TouchIDs()
	if len(IDs) > 0 { //タッチを検出...............
		//今までタッチしていなかった
		if !this.IsTouch {
			//タッチし直ししたらその箇所を中心点に
			this.firstTouchX, this.firstTouchY = GetTouchXY(IDs[0])
		}
		this.IsTouch = true
		//タッチ座標取得
		newX, _ := GetTouchXY(IDs[0])
		//前回タッチ箇所との差分が移動距離
		sa := float64(newX-this.firstTouchX) * 2

		//ねこ移動時は差分計算は常に前フレームとの差にする（ひこーきは初回位置との相対)
		if neco {
			this.firstTouchX = newX
		}
		return sa
	} else { //タッチしていない
		this.IsTouch = false
	}

	return 0
}

// スペース押し続け
func IsPressSpace() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// タッチ有無
func IsTouch() bool {
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])
	return len(touchIDs) > 0
}

// タッチされている箇所のIDを全部返す
func TouchIDs() []ebiten.TouchID {
	//タッチ
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)
	return touchIDs
}

// 指定IDのタッチ座標を返す
func GetTouchXY(touchID ebiten.TouchID) (x, y int) {
	x, y = ebiten.TouchPosition(touchID)
	return x, y
}
