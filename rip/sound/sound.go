package sound

import (
	"Rip/loader"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	FALL      = 0
	MISS      = 1
	MISS2     = 2
	PARACLOSE = 3
	PARAOPEN  = 4
	SUCCESS   = 5
	CATCH     = 6
	GAMEOVER  = 7
	START     = 8
	ONEUP     = 9
	PUSHONEUP = 10
)

type sTrcut struct {
	audioContext *audio.Context //追加
	seBytes      []byte         //追加
	sePlayer     *audio.Player
	muteTime     int
}

type SoundStruct struct {
	sTrcut [11]sTrcut
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

	readMp3(&s.sTrcut[FALL], audioContext, "mp3/fall.mp3")
	readMp3(&s.sTrcut[MISS], audioContext, "mp3/miss.mp3")
	readMp3(&s.sTrcut[MISS2], audioContext, "mp3/miss2.mp3")
	readMp3(&s.sTrcut[PARACLOSE], audioContext, "mp3/paraclose.mp3")
	readMp3(&s.sTrcut[PARAOPEN], audioContext, "mp3/paraopen.mp3")
	readMp3(&s.sTrcut[SUCCESS], audioContext, "mp3/success.mp3")
	readMp3(&s.sTrcut[CATCH], audioContext, "mp3/catch.mp3")
	readMp3(&s.sTrcut[GAMEOVER], audioContext, "mp3/gameover.mp3")
	readMp3(&s.sTrcut[START], audioContext, "mp3/start.mp3")
	readMp3(&s.sTrcut[ONEUP], audioContext, "mp3/oneup.mp3")
	readMp3(&s.sTrcut[PUSHONEUP], audioContext, "mp3/pushoneup.mp3")

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
	file, err := loader.Open(fname)
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
