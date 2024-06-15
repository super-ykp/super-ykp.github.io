package loader

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	m_URL string
)

func SetUrl(url string) {
	m_URL = url
}

func Open(fname string) (io.Reader, error) {
	if m_URL == "" {
		fileData, err := os.Open(fname)
		return fileData, err
	}

	response, err := http.Get(m_URL + "/" + fname)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	reader := bytes.NewReader(body)
	return reader, nil
}
