package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
)

func main() {
	resp, _ := http.Get("http://example.com")
	body, _ := ioutil.ReadAll(resp.Body)
	s := bytes.NewReader(body)
	r, _ := charset.NewReaderLabel("euc-jp", s)
	bytes, _ := ioutil.ReadAll(r)
	fmt.Println(string(bytes))
}
