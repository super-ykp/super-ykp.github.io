package main

import (
	"image/color"
	"log"

	"Yoko/bg"
	"Yoko/camera"
	"Yoko/common"
	"Yoko/effect"
	"Yoko/enemy"
	"Yoko/fileloader"
	"Yoko/input"
	"Yoko/item"
	"Yoko/magic"
	"Yoko/manager"
	"Yoko/mytext"
	"Yoko/myui"
	"Yoko/player"
	"Yoko/powerup"
	"Yoko/saveload"
	"Yoko/skill"
	"Yoko/sound"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 構造体
type GameStruct struct {
	bg      bg.BgStruct
	player  player.PlayerStruct
	enemy   enemy.EnemyStruct
	effect  effect.EffectStruct
	myui    myui.UIStruct
	item    item.ItemStruct
	input   input.InputStruct
	sound   sound.SoundStruct
	manager manager.ManagerStruct
	powerup powerup.PowerupStruct
	skill   skill.SkillStruct
	magic   magic.MagicStruct
}

const (
	GameModeInit      = 0
	GameModeLoading   = 1
	GameModePushStart = 2
	GameModeGame      = 3
)

// グローバル変数---------------------------------------
var (
	m_TitlePngImage *ebiten.Image //イメージ
	m_MainPngImage  *ebiten.Image //イメージ
	GameMode        int
)

// ======================================================================================
// エントリポイント
// ======================================================================================
func main() {
	fileloader.Init()

	//ウインドウのサイズ
	ebiten.SetWindowSize(common.SCREEN_WIDTH*2, common.SCREEN_HEIGHT*2)
	//タイトルバーの名前
	ebiten.SetWindowTitle("Yoko")

	//メインインスタンス
	gameStruct := &GameStruct{}

	GameMode = GameModeInit

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
	case GameModeInit:
		GameMode = GameModeLoading
	case GameModePushStart:
		UpdatePushStart(g)
	case GameModeLoading:
		UpdateLoading(g)
	case GameModeGame:
		UpdateGame(g)
	}

	return nil
}

// ローディング
func UpdateLoading(g *GameStruct) error {
	if fileloader.IsUrlStandBy() {
		Init(g)
		GameMode = GameModePushStart
	}
	return nil
}

// プッシュスタート
func UpdatePushStart(g *GameStruct) error {

	if input.GetPress() || input.GetKeyPress() {
		common.IsFastTouch = true
		sound.Play(sound.DUMMY)
		GameMode = GameModeGame
	}
	return nil
}

// ゲームモード
func UpdateGame(g *GameStruct) error {
	for i := 0; i < common.DEBUGSPEED; i++ {
		camera.Update()
		bg.Update()
		if manager.Update() {
			if magic.Update() {
				item.Update()
				player.Update()
				enemy.Update()
				effect.Update()
			} else {
				enemy.Update2()
			}
		}
		myui.Update()
		powerup.Update()
	}
	return nil
}

// ebitenが要求する仮想画面サイズ。----------------------------------------------
// ウインドウ内で、ここまでが描画領域になる
func (g *GameStruct) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.SCREEN_WIDTH, common.SCREEN_HEIGHT
}

// ======================================================================================
// 描画処理
// ======================================================================================
func (g *GameStruct) Draw(screen *ebiten.Image) {
	switch GameMode {
	case GameModeInit, GameModeLoading:
		DrawLoading(screen, g)
	case GameModePushStart:
		DrawPushStart(screen, g)
	case GameModeGame:
		DrawGame(screen, g)
	}
	//DEBUG
	ebitenutil.DebugPrint(screen, common.DebugStr)
}

// ローディング
func DrawLoading(screen *ebiten.Image, g *GameStruct) error {
	ebitenutil.DebugPrint(screen, "Now Loading ...")
	return nil
}

// プッシュスタート
func DrawPushStart(screen *ebiten.Image, g *GameStruct) error {

	drawImageOption := ebiten.DrawImageOptions{}
	ofX := +common.SCREEN_WIDTH/2 - 240/2
	drawImageOption.GeoM.Translate(float64(ofX), float64(-15))
	screen.DrawImage(m_TitlePngImage, &drawImageOption)
	c := color.RGBA{0, 0, 255, 255}
	mytext.DrawG(screen, 30, "NECONECO RAID", 0, c)
	mytext.DrawG(screen, 200, "TapScreen or", 0, c)
	mytext.DrawG(screen, 220, "Press Space Key", 0, c)
	return nil
}

// ゲームモード
func DrawGame(screen *ebiten.Image, g *GameStruct) error {
	bg.Draw(screen)
	myui.Draw1(screen)
	item.Draw(screen)
	enemy.Draw(screen)
	effect.Draw(screen)
	manager.Draw(screen)
	myui.Draw2(screen)
	player.Draw2(screen)
	magic.Draw(screen)
	powerup.Draw(screen)
	myui.Draw3(screen)

	return nil
}

// ======================================================================================
// 他
// ======================================================================================
// 初期化--------------------------------------------------------------------
func Init(g *GameStruct) {
	mytext.Init()
	//ロード
	m_MainPngImage = fileloader.GetPngImg("img/main.png")
	m_TitlePngImage = fileloader.GetPngImg("img/title.png")

	input.Init(&g.input)
	sound.Init(&g.sound)
	bg.Init(&g.bg, m_MainPngImage)
	effect.Init(&g.effect, m_MainPngImage)
	item.Init(&g.item, m_MainPngImage)
	skill.Init(&g.skill)
	player.Init(&g.player, m_MainPngImage)
	enemy.Init(&g.enemy, m_MainPngImage, &g.player)
	magic.Init(&g.magic, m_MainPngImage)
	manager.Init(&g.manager)
	powerup.Init(&g.powerup, m_MainPngImage, &g.player, &g.manager.Gparam)
	myui.Init(&g.myui, m_MainPngImage, &g.player, &g.manager.Gparam)

	manager.SetClass(&g.input, &g.sound, &g.bg, &g.player, &g.enemy, &g.effect, &g.item, &g.powerup, &g.skill, &g.magic)

	//データロード-----------------------------------------------------
	g.manager.Gparam.FirstPowerUp = true

	gm, pl, en := saveload.LoadAllData()
	if gm != nil {
		g.manager.Gparam = *gm
		g.player.PState = *pl
		g.enemy.EList = en.EList

		for i := 0; i < len(g.enemy.EList); i++ {
			ec := &g.enemy.EList[i]
			if ec.HP <= 0 {
				ec.IsUse = false
			}
			ec.Y = 0
			ec.EStete = enemy.S_STAND

			if int(ec.X)-ec.Radius < g.player.PState.X+31 {
				ec.X = float64(g.player.PState.X + 31 + ec.Radius)
			}
		}

		if g.player.PState.HP <= 0 {
			manager.ForceGameOver()

		}
		myui.SetUIMode(myui.UI_NOMAL)
	}
}
