package feedgetter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	defer res.Body.Close()
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

// Get feeds
func Get(targetURL string) (feeds []string, err error) {
	u, _ := url.Parse(targetURL)
	selectors := []string{
		`link[type="application/atom+xml"]`,
		`link[type="text/xml"]`,
		`link[type="application/rss+xml"]`,
		`link[type="application/x.atom+xml"]`,
		`link[type="application/x-atom+xml"]`,
	}
	body, err := Body(targetURL)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}
	for _, sel := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			feed, exist := s.Attr("href")
			if exist {
				if strings.Contains(feed, "://") {
					feeds = append(feeds, feed)
				} else {
					_f, _ := url.Parse(feed)
					feeds = append(feeds, u.ResolveReference(_f).String())
				}
			}
		})
	}
	return
}
