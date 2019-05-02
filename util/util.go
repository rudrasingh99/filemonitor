package util

import (
	"bufio"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 60 * time.Second}

func CreateRequest(url string) (*bufio.Reader, error) {

	resp, err := client.Get(url)
	resp.Header.Set("User-Agent", "Mozilla/5.0 (compatible; filemonitor/0.1; +github.com/kapytein/filemonitor)")

	if err != nil {
		return nil, err
	}

	return bufio.NewReader(resp.Body), nil

}

func ReadBuffer(b *bufio.Reader) []byte {

	scanner := bufio.NewScanner(b)

	var str string
	for scanner.Scan() {
		str += scanner.Text()
	}

	return []byte(str)

}
