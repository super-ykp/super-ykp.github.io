package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputStruct struct {
	IsTouch bool
	X       int
	Y       int
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

// キー押下
func GetKeyPress() bool {
	//キーで移動.................
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		return true
	}
	return false
}

// マウスまたはタッチ
func GetPress() bool {
	//マウス押下
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	//タッチで移動................
	IDs := touchIDs()
	return len(IDs) > 0
}

func GetTouchPos() (int, int, bool) {
	//カーソル
	this.IsTouch = false
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		this.X, this.Y = ebiten.CursorPosition()
		this.IsTouch = true
	} else { //タッチ
		IDs := touchIDs()
		if len(IDs) > 0 {
			this.X, this.Y = touchXY(IDs[0])
			this.IsTouch = true
		} else {
			this.IsTouch = false
		}
	}

	//----------------------------------
	return this.X, this.Y, this.IsTouch
}

// マルチタッチ用
func GetTouchPos2() (int, int, bool) {
	//タッチ
	IDs := touchIDs()
	if len(IDs) >= 2 {
		X, Y := touchXY(IDs[1])
		return X, Y, true
	} else {
		return 0, 0, false
	}
	//----------------------------------
}

// タッチされている箇所のIDを全部返す
func touchIDs() []ebiten.TouchID {
	//タッチ
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)
	return touchIDs
}

// 指定IDのタッチ座標を返す
func touchXY(touchID ebiten.TouchID) (x, y int) {
	x, y = ebiten.TouchPosition(touchID)
	return x, y
}
