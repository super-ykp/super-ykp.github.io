package sound

import (
	"Neco/loader"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	EXPLOSION      = 0
	EXPLOSION_BOSS = 1
	SHOT           = 2
	KNOKBACK       = 3
	DAMAGE         = 4
	BOSSHIT        = 5
	PLAYERDEATH    = 6
	START          = 7
)

type sTrcut struct {
	audioContext *audio.Context //追加
	seBytes      []byte         //追加
	sePlayer     *audio.Player
	muteTime     int
}

type SoundStruct struct {
	sTrcut [10]sTrcut
}

func Init(s *SoundStruct) {
	audioContext := audio.NewContext(44100)

	readMp3(&s.sTrcut[EXPLOSION], audioContext, "mp3/explosion.mp3")
	readMp3(&s.sTrcut[EXPLOSION_BOSS], audioContext, "mp3/explosionboss.mp3")
	readMp3(&s.sTrcut[SHOT], audioContext, "mp3/shot.mp3")
	readMp3(&s.sTrcut[KNOKBACK], audioContext, "mp3/knokback.mp3")
	readMp3(&s.sTrcut[DAMAGE], audioContext, "mp3/damage.mp3")
	readMp3(&s.sTrcut[BOSSHIT], audioContext, "mp3/bosshit.mp3")
	readMp3(&s.sTrcut[PLAYERDEATH], audioContext, "mp3/playerdeath.mp3")
	readMp3(&s.sTrcut[START], audioContext, "mp3/start.mp3")
}
func Update(s *SoundStruct) {
	for i := 0; i < len(s.sTrcut); i++ {
		if 0 < s.sTrcut[i].muteTime {
			s.sTrcut[i].muteTime--
		}
	}
}

func Play(s *SoundStruct, i int) {
	if s.sTrcut[i].muteTime > 0 {
		return
	}
	if s.sTrcut[i].sePlayer != nil {
		s.sTrcut[i].sePlayer.Close()
	}
	s.sTrcut[i].sePlayer = s.sTrcut[i].audioContext.NewPlayerFromBytes(s.sTrcut[i].seBytes)
	s.sTrcut[i].sePlayer.Play()

}
func Stop(s *SoundStruct, i int, time int) {
	if s.sTrcut[i].sePlayer == nil {
		return
	}
	s.sTrcut[i].sePlayer.Close()
	s.sTrcut[i].muteTime = time
}

// -------------------
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
