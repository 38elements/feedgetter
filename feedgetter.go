package feedgetter

import (
	"bytes"
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

func encode(h *http.Header, body []byte) (encode string) {
	contentType := ""
	for k, v := range *h {
		if strings.ToLower(k) == "content-type" {
			contentType = v[0]
			break
		}
	}
	_, encode, _ = charset.DetermineEncoding(body, contentType)
	return
}

func webpage(url string) (body string, err error) {
	res, err := response(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	enc := encode(&res.Header, byteBody)
	if err != nil {
		return
	}
	s := bytes.NewReader(byteBody)
	r, err := charset.NewReaderLabel(enc, s)
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	body = string(bytes)
	return
}

// Get feeds
func Get(targetURL string) (feeds []string, err error) {
	body, err := webpage(targetURL)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}
	selectors := []string{
		`link[type="application/atom+xml"]`,
		`link[type="text/xml"]`,
		`link[type="application/rss+xml"]`,
		`link[type="application/x.atom+xml"]`,
		`link[type="application/x-atom+xml"]`,
	}
	u, _ := url.Parse(targetURL)
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
