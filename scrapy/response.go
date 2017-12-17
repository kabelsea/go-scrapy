package scrapy

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var regexpURL, _ = regexp.Compile("<a[^<]*href=['\"](.*?)['\"][^<]*>")

// Response is the representation of a HTTP response made by a Spider
type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int

	// Body is the content of the Response
	Body []byte

	// Ctx is a context between a Request and a Response
	Ctx *Context

	// Request is the Request object of the response
	Request *Request

	// Headers contains the Response's HTTP headers
	Headers *http.Header
}

type ResponseChannel chan *Response

// NewResponse creates Response instance and initialized it default values
func NewResponse(req *Request, ctx *Context) *Response {
	return &Response{
		Request: req,
		Ctx:     ctx,
	}
}

// Method parse and return all links from http response
func (r *Response) ExtractLinks() []string {
	var links []string

	if len(r.Body) != 0 {
		matches := regexpURL.FindAllStringSubmatch(string(r.Body), -1)
		if len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 0 {
					link := match[1]

					u, _ := url.Parse(link)
					if u.Host == "" {
						req := r.Request
						link = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, link)
					}
					links = append(links, link)
				}
			}
		}
	}
	return links
}

// Handlers returns list callback functions for Response instance
func (r *Response) Handlers() []Handler {
	var handlers []Handler

	for _, lex := range r.Request.Config.Rules {
		if lex.Handler != nil {
			if lex.LinkExtractor.Match(r.Request.URL) {
				handlers = append(handlers, lex.Handler)
			}
		}
	}
	return handlers
}
