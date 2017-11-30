package scrapy

import (
	"io/ioutil"
	"net/http"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

// Spider request type
type Request struct {
	Url       string
	Attempt   int
	Depth     int
	Config    *SpiderConfig
	ParsedURL *url.URL
}

type RequestChannel chan Request

// Create new request object
func NewRequest(link string, config *SpiderConfig) *Request {
	parsedUrl, _ := url.Parse(link)

	return &Request{
		Url:       link,
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
	response := NewResponse(r)

	resp, err := http.Get(r.Url)
	if err != nil {
		r.Attempt++
		logger.Error(err)
		return response
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Attempt++
		logger.Error(err)
		return response
	}

	response.StatusCode = resp.StatusCode
	response.Body = body

	return response
}

// Return true if target request can follow
func (r *Request) CanFollow() bool {
	if r.Depth >= r.Config.MaxDepth {
		return false
	}

	if len(r.Config.Rules) != 0 {
		for _, lex := range r.Config.Rules {
			if lex.Follow {
				return lex.LinkExtractor.Match(r.ParsedURL)
			}
		}
	}
	return false
}

// Returns true if target request can parse
func (r *Request) CanParse() bool {
	if len(r.Config.Rules) != 0 {
		for _, lex := range r.Config.Rules {
			if lex.Handler != nil {
				return lex.LinkExtractor.Match(r.ParsedURL)
			}
		}
	}
	return false
}
