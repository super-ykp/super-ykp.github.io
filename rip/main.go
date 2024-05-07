package main

import (
	"image/color"
	"log"

	"Rip/bg"
	"Rip/common"
	"Rip/effect"
	"Rip/enemy"
	"Rip/food"
	"Rip/input"
	"Rip/irand"
	"Rip/loader"
	"Rip/manager"
	"Rip/mytext"
	"Rip/neco"
	"Rip/plane"
	"Rip/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 外部クラスとの通信を行うためのインターフェース
type GameStruct struct {
	manager manager.ManagerStruct
	plane   plane.PlaneStruct
	neco    neco.NecoStruct
	enemy   enemy.EnemyStruct
	irand   irand.IrandStruct
	effect  effect.EffectStruct
	food    food.FoodStruct
	bg      bg.BgStruct
	sound   sound.SoundStruct
	input   input.InputStruct
}

const (
	GameModeUrlWaiting = 0
	GameModeLoading    = 1
	GameModeGame       = 2
)

// グローバル変数---------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージ
	GameMode       int
)

// ======================================================================================
// エントリポイント
// ======================================================================================
func main() {
	loader.Init()

	//ウインドウのサイズ。+60はタッチ領域
	ebiten.SetWindowSize(common.SCREEN_WIDTH*2, (common.SCREEN_HEIGHT+60)*2)
	//タイトルバーの名前
	ebiten.SetWindowTitle("Neco")

	//メインインスタンス
	gameStruct := &GameStruct{}

	//メインスレッドスタート。エラーなら返ってくるのでログに
	if err := ebiten.RunGame(gameStruct); err != nil {
		log.Fatal(err)
	}
}

// ======================================================================================
// 計算処理
// ======================================================================================
func (g *GameStruct) Update() error {
	//各種計算
	switch GameMode {
	case GameModeUrlWaiting:
		UpdateUrlWaiting(g)
	case GameModeLoading:
		UpdateLoading(g)
	case GameModeGame:
		UpdateGame(g)
	}

	return nil
}

// URL登録待ち
func UpdateUrlWaiting(g *GameStruct) error {
	if loader.IsUrlStandBy() {
		GameMode = GameModeLoading
		//スマホでない場合、タッチ領域はいらない
		if !common.IsSmartPhone {
			ebiten.SetWindowSize(common.SCREEN_WIDTH*2, (common.SCREEN_HEIGHT)*2)
		}
	}
	return nil
}

// ローディング
func UpdateLoading(g *GameStruct) error {
	Init(g)
	GameMode = GameModeGame
	return nil
}

// ゲームモード
func UpdateGame(g *GameStruct) error {
	manager.Update()
	enemy.Update()
	effect.Update()
	plane.Update()
	food.Update()
	neco.Update()
	sound.Update()
	return nil
}

// ebitenが要求する仮想画面サイズ。----------------------------------------------
// ウインドウ内で、ここまでが描画領域になる
func (g *GameStruct) Layout(outsideWidth, outsideHeight int) (int, int) {
	//スマホの場合タッチ領域を作る
	if common.IsSmartPhone {
		return common.SCREEN_WIDTH, common.SCREEN_HEIGHT + 60
	} else {
		return common.SCREEN_WIDTH, common.SCREEN_HEIGHT

	}
}

// ======================================================================================
// 描画処理
// ======================================================================================
func (g *GameStruct) Draw(screen *ebiten.Image) {
	switch GameMode {
	case GameModeUrlWaiting:
		DrawLoading(screen, g)
	case GameModeLoading:
		DrawLoading(screen, g)
	case GameModeGame:
		DrawGame(screen, g)
	}

}

// ローディング
func DrawLoading(screen *ebiten.Image, g *GameStruct) error {
	ebitenutil.DebugPrint(screen, "Now Loading ...")
	return nil
}

// ゲームモード
func DrawGame(screen *ebiten.Image, g *GameStruct) error {
	bg.Draw(screen)
	irand.Draw(screen)
	food.Draw(screen)
	enemy.Draw(screen)
	effect.Draw(screen)
	neco.Draw(screen)
	plane.Draw(screen)
	manager.Draw(screen)

	if common.IsSmartPhone {
		vector.DrawFilledRect(screen, 0, common.SCREEN_HEIGHT, common.SCREEN_WIDTH, 60, color.RGBA{128, 128, 128, 0xff}, false)
	}
	return nil
}

// ======================================================================================
// 他
// ======================================================================================
// 初期化--------------------------------------------------------------------
func Init(g *GameStruct) {
	//ロード
	m_MainPngImage = loader.GetPngImg("img/main.png")

	sound.Init(&g.sound)
	manager.Init(&g.manager, m_MainPngImage)
	neco.Init(&g.neco, m_MainPngImage, &g.sound)
	plane.Init(&g.plane, m_MainPngImage)
	enemy.Init(&g.enemy, m_MainPngImage)
	irand.Init(&g.irand, m_MainPngImage)
	effect.Init(&g.effect, m_MainPngImage)
	food.Init(&g.food, m_MainPngImage)
	bg.Init(&g.bg)
	input.Init(&g.input)
	mytext.Init()

	manager.SetClass(&g.bg, &g.plane, &g.neco, &g.irand, &g.effect, &g.enemy, &g.food, &g.input, &g.sound)
}
