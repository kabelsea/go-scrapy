package scrapy

import (
	"io/ioutil"
	"net/http"
)

type Request struct {
	Url     string
	Body    []byte
	Meta    map[string]string
	Headers map[string]string
	Retry   int
}

// Method send get http request
func (r *Request) Download() bool {
	resp, err := http.Get(r.Url)
	if err != nil {
		r.Retry++
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Retry++
		return false
	}
	r.Body = body

	return true
}
