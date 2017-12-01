package scrapy

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

type RequestErrorType int

// Errors on http requests
const (
	RequestTimeoutError RequestErrorType = iota
	RequestUnknownError
)

type RequestError struct {
	err RequestErrorType
	msg string
}

// Implementation Error interface for custom Request error type
func (e *RequestError) Error() string {
	if e.err == RequestTimeoutError {
		return fmt.Sprintf("http request timeout exceeded")
	}
	return e.msg
}

// Method convert error object to RequestErrorType
func (e *RequestError) Load(err error) {
	if er, ok := err.(net.Error); ok && er.Timeout() {
		e.err = RequestTimeoutError
		return
	}

	e.err = RequestUnknownError
}

// Spider request type
type Request struct {
	Url        string
	Attempt    int
	Depth      int
	Headers    map[string]string
	UserAgent  string
	Config     *SpiderConfig
	HttpClient *http.Client
	ParsedURL  *url.URL
}

type RequestChannel chan Request

// Create new request object
func NewRequest(link string, config *SpiderConfig) *Request {
	parsedUrl, _ := url.Parse(link)

	client := &http.Client{
		Timeout: config.DownloadTimeout,
	}

	return &Request{
		Url:        link,
		Config:     config,
		Headers:    config.RequestHeaders,
		UserAgent:  config.UserAgent,
		ParsedURL:  parsedUrl,
		HttpClient: client,
	}
}

// Method make http request, check http status code and return Response object
func (r *Request) Process() (*Response, error) {
	response := NewResponse(r)
	err := &RequestError{}

	req, e := http.NewRequest("GET", r.Url, nil)
	if e != nil {
		err.Load(e)
		return response, err
	}

	// Set http request headers
	if r.Headers != nil {
		for h, v := range r.Headers {
			req.Header.Set(h, v)
		}
	}

	// Set user agent for http request
	if r.UserAgent != "" {
		req.Header.Set("User-Agent", r.UserAgent)
	}

	resp, e := r.HttpClient.Do(req)
	if e != nil {
		err.Load(e)
		logger.Error(e)
		return response, err
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		err.Load(e)
		logger.Error(e)
		return response, err
	}

	response.StatusCode = resp.StatusCode
	response.Body = body

	return response, nil
}

// Check retry request
func (r *Request) CanRetry() bool {
	if r.Attempt <= r.Config.RetryTimes {
		return true
	}
	return false
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
