package scrapy

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
)

type RequestMethod string

// Available HTTP methods
const (
	GetMethod     RequestMethod = "GET"
	PostMethod    RequestMethod = "POST"
	PutMethod     RequestMethod = "PUT"
	DeleteMethod  RequestMethod = "DELETE"
	OptionsMethod RequestMethod = "OPTIONS"
	HeadMethod    RequestMethod = "HEAD"
)

// Request is the representation of a HTTP request made by a Spider
type Request struct {
	Config *SpiderConfig

	// Unique identifier of the request
	ID uuid.UUID

	// URL is the parsed URL of the HTTP request
	URL *url.URL

	// Headers contains the Request's HTTP headers
	Headers *http.Header

	// Depth is the number of the parents of the request
	Depth int

	// Method is the HTTP method of the request
	Method RequestMethod

	// Body is the request body which is used on POST/PUT requests
	Body io.Reader

	// Ctx is a context between a Request and a Response
	Ctx *Context
}

type RequestChannel chan *Request

// NewRequest creates a Request instance
func NewRequest(link string, config *SpiderConfig) (*Request, error) {
	ctx := NewContext()

	// Copy default headers from config
	//
	headers := http.Header{}
	for k, v := range config.RequestHeaders {
		headers[k] = v
	}

	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if config.UserAgent != "" {
		headers.Set("User-Agent", config.UserAgent)
	}

	req := &Request{
		ID:      uuid.NewV4(),
		URL:     u,
		Method:  GetMethod,
		Config:  config,
		Headers: &headers,
		Ctx:     ctx,
	}
	return req, nil
}

// Process makes http request, check http status code and return Response object
func (r *Request) Process() (*Response, error) {
	response := NewResponse(r, r.Ctx)

	client := &http.Client{
		Timeout: r.Config.DownloadTimeout,
	}

	req, err := http.NewRequest("GET", r.URL.String(), nil)
	if err != nil {
		return response, err
	}

	// Set http request headers
	if r.Headers != nil {
		req.Header = *r.Headers
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return response, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return response, err
	}

	response.StatusCode = resp.StatusCode
	response.Body = body

	return response, nil
}

// CanFollow returns true if target Request can follow
func (r *Request) CanFollow() bool {
	if r.Depth >= r.Config.MaxDepth {
		return false
	}

	if len(r.Config.Rules) != 0 {
		for _, lex := range r.Config.Rules {
			if lex.Follow {
				return lex.LinkExtractor.Match(r.URL)
			}
		}
	}
	return false
}

// CanParse returns true if target Request can parse
func (r *Request) CanParse() bool {
	if len(r.Config.Rules) != 0 {
		for _, lex := range r.Config.Rules {
			if lex.Handler != nil {
				return lex.LinkExtractor.Match(r.URL)
			}
		}
	}
	return false
}
