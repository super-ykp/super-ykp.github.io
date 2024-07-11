package sound

import (
	"Yoko/common"
	"Yoko/fileloader"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	DUMMY      = 0
	EXPLOSION1 = 1
	EXPLOSION2 = 2
	HIT1       = 3
	HIT2       = 4
	FURI       = 5
	PHIT       = 6
	PDOWN      = 7
	PGURU      = 8
	CATCH      = 9
	GAMEOVER   = 10
	HAPPY      = 11
)

type sTrcut struct {
	audioContext *audio.Context //追加
	seBytes      []byte         //追加
	sePlayer     *audio.Player
	muteTime     int
}

type SoundStruct struct {
	sTrcut [20]sTrcut
}

var (
	this *SoundStruct
)

// --===================================================================================
// 　効果音再生
// --===================================================================================
// 初期化
func Init(s *SoundStruct) {
	this = s
	audioContext := audio.NewContext(44100)

	readMp3(&s.sTrcut[DUMMY], audioContext, "mp3/dumy.mp3")
	readMp3(&s.sTrcut[EXPLOSION1], audioContext, "mp3/explosion1.mp3")
	readMp3(&s.sTrcut[EXPLOSION2], audioContext, "mp3/explosion2.mp3")
	readMp3(&s.sTrcut[HIT1], audioContext, "mp3/hit1.mp3")
	readMp3(&s.sTrcut[HIT2], audioContext, "mp3/hit2.mp3")
	readMp3(&s.sTrcut[FURI], audioContext, "mp3/furi.mp3")
	readMp3(&s.sTrcut[PHIT], audioContext, "mp3/phit.mp3")
	readMp3(&s.sTrcut[PDOWN], audioContext, "mp3/pdown.mp3")
	readMp3(&s.sTrcut[PGURU], audioContext, "mp3/pguru.mp3")
	readMp3(&s.sTrcut[CATCH], audioContext, "mp3/catch.mp3")
	readMp3(&s.sTrcut[GAMEOVER], audioContext, "mp3/gameover.mp3")
	readMp3(&s.sTrcut[HAPPY], audioContext, "mp3/happy.mp3")
}

// 計算
func Update() {
	//チャンネル毎にミュート時間が指定できる
	for i := 0; i < len(this.sTrcut); i++ {
		if 0 < this.sTrcut[i].muteTime {
			this.sTrcut[i].muteTime--
		}
	}
}

// 再生開始
func Play(i int) {
	//タッチが確認されるまで鳴らさない
	if !common.IsFastTouch {
		return
	}
	//ミュート中は鳴らさない
	if this.sTrcut[i].muteTime > 0 {
		return
	}
	//プレイ中なら一旦クローズ
	if this.sTrcut[i].sePlayer != nil {
		this.sTrcut[i].sePlayer.Close()
	}
	//プレイ開始
	this.sTrcut[i].sePlayer = this.sTrcut[i].audioContext.NewPlayerFromBytes(this.sTrcut[i].seBytes)
	this.sTrcut[i].sePlayer.Play()
}

// 強制停止
func Stop(i int, time int) {
	//鳴ってない
	if this.sTrcut[i].sePlayer == nil {
		return
	}
	//停止
	this.sTrcut[i].sePlayer.Close()
	this.sTrcut[i].muteTime = time
}

// mp3の読み込み
func readMp3(s *sTrcut, a *audio.Context, fname string) {
	s.audioContext = a
	file, err := fileloader.Open(fname)
	if err != nil {
		panic(err)
	}

	// MP3をゲーム内で使えるようにデコード
	src, err2 := mp3.DecodeWithoutResampling(file)
	if err2 != nil {
		panic(err2)
	}
	// デコードした内容を用意していたbyte変数へ入れておく
	var err3 error
	s.seBytes, err3 = io.ReadAll(src)
	if err3 != nil {
		panic(err3)
	}
}
