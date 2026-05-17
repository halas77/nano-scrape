package engine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Request struct {
	Url         string
	Method      string
	Header      http.Header
	ContentType string
}

func InitRequest(url string, method string, header http.Header) Request {
	req := Request{Url: url, Method: method, Header: header}
	return req
}

func (r Request) Execute(body ...map[string][]string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status code %d\n", resp.StatusCode)
		return nil, errors.New("Server responded with status code: " + string(rune(resp.StatusCode)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return bodyBytes, nil
}
