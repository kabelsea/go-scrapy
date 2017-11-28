package scrapy

import (
	"net/http"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

// Spider request type
type Request struct {
	Url       string
	Attempt   int
	Config    *SpiderConfig
	ParsedURL *url.URL
}

type RequestChannel chan Request

// Make request method
func MakeRequest(link string, config *SpiderConfig) *Request {
	parsedUrl, _ := url.Parse(link)

	return &Request{
		Url:       link,
		Attempt:   0,
		Config:    config,
		ParsedURL: parsedUrl,
	}
}

// Return request headers
func (r *Request) Headers() map[string]string {
	return r.Config.RequestHeaders
}

// Method make http request, check http status code and return Response object
func (r *Request) Download() *Response {
	response := MakeResponse(r)

	resp, err := http.Get(r.Url)
	if err != nil {
		r.Attempt++
		logger.Error(err)
		return response
	}
	defer resp.Body.Close()

	response.HttpResponse = resp
	return response
}
