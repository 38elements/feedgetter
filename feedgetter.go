package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func response(url string) (res *http.Response, err error) {
	c := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	return c.Get(url)
}

func encode(h *http.Header) (encode string, err error) {
	err = nil
	for k, v := range *h {
		if strings.ToLower(k) == "content-type" {
			temp := strings.Split(v[0], "=")
			if len(temp) != 2 {
				err = errors.New("There is not charset in header")
			}
			encode = temp[1]
			return
		}
	}
	err = errors.New("There is not content-type in header")
	return
}

// Body gets content of url
func Body(url string) (body string, err error) {
	res, err := response(url)
	if err != nil {
		return
	}
	enc, err := encode(&res.Header)
	if err != nil {
		return
	}
	_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	s := bytes.NewReader(_body)
	r, err := charset.NewReaderLabel(enc, s)
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(r)
	body = string(bytes)
	return
}

func main() {
	url := "http://example.com"
	body, _ := Body(url)
	fmt.Println(body)
}
