package main

import (
	"Neco/background"
	"Neco/collision/collision_dp"
	"Neco/collision/collision_se"
	"Neco/common"
	"Neco/danmaku"
	"Neco/enemy"
	"Neco/explosion"
	"Neco/loader"
	"Neco/player"
	"Neco/shot"
	"Neco/sound"
	"Neco/text"
	"fmt"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// グローバル変数みたいなものらしい。変数名が先、定義が後---------------------------------------
var (
	m_MainPngImage *ebiten.Image //イメージ
	m_BgPingImage  *ebiten.Image
	m_TitleImage   *ebiten.Image

	m_fpsCounter time.Time
	m_fps        int
)

// 構造体-------------------------------------------------------------------
type GameStruct struct {
	background background.BackgroundStruct
	player     player.PlayerStruct
	danmaku    danmaku.DanmakuStruct
	shot       shot.ShotStruct
	enemy      enemy.EnemyStruct
	explosion  explosion.ExplosionStruct

	sound sound.SoundStruct
}

// ======================================================================================
// エントリポイント
// ======================================================================================
func main() {
	loader.Init()

	//ウインドウのサイズ。描画領域はこの範囲すべてでは無い
	ebiten.SetWindowSize(common.SCREEN_WIDTH*2, common.SCREEN_HEIGHT*2)
	//タイトルバーの名前
	ebiten.SetWindowTitle("Neco")

	//メインインスタンス
	gameStruct := &GameStruct{}

	//メインスレッドスタート。エラーなら返ってくるのでログに
	if err := ebiten.RunGame(gameStruct); err != nil {
		log.Fatal(err)
	}
}

// ebitenが要求する仮想画面サイズ。----------------------------------------------
// ウインドウ内で、ここまでが描画領域になる
func (g *GameStruct) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.SCREEN_WIDTH, common.SCREEN_HEIGHT
}

// ======================================================================================
// 計算処理
// ======================================================================================
func (g *GameStruct) Update() error { //関数の引数も、前にあるのが変数名で、後ろが型

	//前回からの時間差分
	m_fps = int(1e9 / float64(time.Since(m_fpsCounter).Nanoseconds()))
	//今回の時間記録
	m_fpsCounter = time.Now()

	//各種計算-------------------------------------------
	switch common.GameMode {
	case common.GameMOdeUrlWaiting:
		UpdateUrlWaiting(g)
	case common.GameModeLoading:
		UpdateLoading(g)
	case common.GameModeLogo:
		UpdateLogo(g)
	case common.GameModeGame:
		UpdateGame(g)
	}

	return nil
}

// URL登録待ち
func UpdateUrlWaiting(g *GameStruct) error {
	if loader.IsUrlStandBy() {
		common.GameMode = common.GameModeLoading
	}
	return nil
}

// ローディング
func UpdateLoading(g *GameStruct) error {
	Init(g)
	common.GameMode = common.GameModeLogo
	return nil
}

// ロゴモード
func UpdateLogo(g *GameStruct) error {

	//BGループ
	background.Update(&g.background)

	//キー入力があったらゲーム開始
	if IsEnter() && common.KeyRelese {
		sound.Play(&g.sound, sound.START)

		common.Stage = 0
		common.StartStage()

		//モード設定
		common.GameMode = common.GameModeGame
		common.GameState = common.GameStateNomal
		common.GameState = common.GameStateStartStage
		//各処理初期化
		player.Start(&g.player)
		enemy.Start(&g.enemy)
		shot.Start(&g.shot)
		danmaku.Start(&g.danmaku)

	} else { //キー押しっぱなしで遷移するのを避けるため、一旦離すまで待つ
		if !IsEnter() {
			common.KeyRelese = true
		}
	}
	return nil
}

func UpdateGame(g *GameStruct) error {
	//計算
	background.Update(&g.background)
	shot.Update(&g.shot)
	player.Update(&g.player)
	danmaku.Update(&g.danmaku, &g.player)
	enemy.Update(&g.enemy, &g.player)
	explosion.Update(&g.explosion)
	sound.Update(&g.sound)
	text.Update()

	if common.GameState == common.GameStateNomal { //ゲーム中
		collision_dp.Calc(&g.danmaku, &g.player) //当たり判定
		collision_se.Calc(&g.shot, &g.enemy)

	} else if common.GameState == common.GameStateStartStage { //ゲーム開始演出
		if common.Counter < 1000*60 {
			common.Counter++
		}

		if common.Counter > 60*2 { //2秒経ったら敵出力
			common.GameState = common.GameStateNomal
			enemy.Start(&g.enemy)
			player.RrecoveryHp(&g.player)
		}
	} else if common.GameState == common.GameStateGameOver { //ゲームオーバー
		if common.Counter < 1000*60 {
			common.Counter++
		}
		if common.Counter > 60 { //1秒後からキー受付
			if IsEnter() {
				common.GameMode = common.GameModeLogo
				common.KeyRelese = false
			}
		}
	} else if common.GameState == common.GameStateAllClear { //クリア
		if common.Counter < 60*1000 {
			common.Counter++
		}
		if common.Counter > 60*5 { //5秒後からキー受付
			if IsEnter() {
				common.GameMode = common.GameModeLogo
				common.KeyRelese = false
			}
		}
	}

	return nil
}

// ======================================================================================
// 描画処理
// ======================================================================================
func (g *GameStruct) Draw(screen *ebiten.Image) {
	switch common.GameMode {
	case common.GameMOdeUrlWaiting:
	case common.GameModeLoading:
		DrawLoading(screen, g)
	case common.GameModeLogo:
		DrawLogo(screen, g)
	case common.GameModeGame:
		DrawGame(screen, g)
	}

}

// ローディング
func DrawLoading(screen *ebiten.Image, g *GameStruct) error {
	ebitenutil.DebugPrint(screen, "Now Loading ...")
	return nil
}

// ロゴモード
func DrawLogo(screen *ebiten.Image, g *GameStruct) error {
	background.Draw(screen, &g.background)
	//ロゴ
	drawImageOption := &ebiten.DrawImageOptions{}
	drawImageOption.GeoM.Translate(common.SCREEN_WIDTH/2-248/2, 50)
	imageRect := image.Rect(0, 0, 248, 248)
	ebitenImage := m_TitleImage.SubImage(imageRect).(*ebiten.Image)
	screen.DrawImage(ebitenImage, drawImageOption)

	return nil
}

// ゲームモード
func DrawGame(screen *ebiten.Image, g *GameStruct) error {
	background.Draw(screen, &g.background)
	shot.Draw(screen, &g.shot)
	player.Draw(screen, &g.player)
	danmaku.Draw(screen, &g.danmaku)
	enemy.Draw(screen, &g.enemy)
	explosion.Draw(screen, &g.explosion)
	enemy.Draw2(screen, &g.enemy)
	player.Draw2(screen, &g.player)

	if common.GameState == common.GameStateGameOver { //ゲームオーバー
		y := common.SCREEN_WIDTH/3 + (common.SCREEN_WIDTH - common.Counter*3)
		if y < common.SCREEN_WIDTH/3 {
			y = common.SCREEN_WIDTH / 3
		}
		text.DrawText(screen, "GAME OVER", y, 0)

	} else if common.GameState == common.GameStateAllClear { //オールクリア
		y := common.SCREEN_WIDTH/3 + (-common.SCREEN_WIDTH/2 + common.Counter)
		if y > common.SCREEN_WIDTH/3 {
			y = common.SCREEN_WIDTH / 3
		}
		text.DrawText(screen, "ALL CLEAR !", y, 1)

	} else if common.GameState == common.GameStateStartStage { //ゲーム開始

		y := common.SCREEN_HEIGHT/3 + (-300 + common.Counter*10)
		if y >= common.SCREEN_HEIGHT/3 {
			y = common.SCREEN_HEIGHT / 3
		}

		str := fmt.Sprintf("STAGE %d/10 START", common.Stage)
		text.DrawText(screen, str, y, 0)
	}

	//DEBUG--------------------
	//str := fmt.Sprintf("%d ", m_fps)
	//ebitenutil.DebugPrint(screen, str)

	return nil
}

// 初期化--------------------------------------------------------------------
func Init(g *GameStruct) {
	//ロード
	m_MainPngImage = getPngImg("img/main.png")
	m_BgPingImage = getPngImg("img/bg.png")
	m_TitleImage = getPngImg("img/title.png")

	sound.Init(&g.sound)
	background.Init(m_BgPingImage, &g.background)
	player.Init(m_MainPngImage, &g.player, &g.shot, &g.explosion, &g.sound)
	shot.Init(m_MainPngImage, &g.shot, &g.sound)
	danmaku.Init(m_MainPngImage, &g.danmaku, &g.player)
	enemy.Init(m_MainPngImage, &g.enemy, &g.danmaku, &g.explosion, &g.sound)
	explosion.Init(m_MainPngImage, &g.explosion)
	text.Init()

	common.GameMode = common.GameModeLogo
	common.GameState = common.GameStateNomal
}

// pngイメージのロード-----------------------------------------------------
func getPngImg(fname string) *ebiten.Image {
	//go言語の特徴として、戻り値は複数持てるらしい。pngをオープン
	fileData, err := loader.Open(fname)
	//エラーがあったらログに書く
	if err != nil {
		log.Fatal(err)
	}

	//ファイルデータをpngライブラリでイメージに変換
	img, err := png.Decode(fileData)
	//またエラー判定
	if err != nil {
		log.Fatal(err)
	}
	//ebiten用に更にイメージ化
	return ebiten.NewImageFromImage(img)
}

// キー押下---------------------------------------------
func IsEnter() bool {
	//Enterキー
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		return true
	}
	//マウス
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	//タッチ
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])
	if len(touchIDs) > 0 {
		return true
	}
	return false
}
