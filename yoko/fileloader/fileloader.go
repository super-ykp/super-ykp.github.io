package fileloader

import (
	"Yoko/common"
	"bytes"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	m_URL string
)

// Webブラウザでの起動の場合カレントURLが設定する。
func SetURL(url string) {
	m_URL = url
	common.IsSmartPhone = false
}

// ファイルの読み込み
func Open(fname string) (io.Reader, error) {

	if m_URL == "" { //exe起動の場合
		//普通にファイルを開いて帰る
		fileData, err := os.Open(fname)
		return fileData, err
	} else { //ここからはブラウザ起動の場合
		//URLからファイルを取得
		response, err := http.Get(m_URL + "/" + fname)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		//ボディ部だけ取り出す
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		//開く
		reader := bytes.NewReader(body)
		return reader, nil
	}
}

// pngイメージの展開-----------------------------------------------------
func GetPngImg(fname string) *ebiten.Image {
	//pngをオープン
	fileData, err := Open(fname)
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
