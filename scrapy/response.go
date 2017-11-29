package scrapy

import (
	"fmt"
	"net/url"
	"regexp"
)

var regexpURL, _ = regexp.Compile("<a[^<]*href=['\"](.*?)['\"][^<]*>")

// Spider response object
type Response struct {
	Url        string
	StatusCode int
	Body       []byte
	Request    *Request
}

type ResponseChannel chan *Response

// Create new response object, initialize with default values
func NewResponse(req *Request) *Response {
	return &Response{Request: req}
}

// Return Response status code
func (r *Response) Success() bool {
	for _, code := range r.Request.Config.RetryHttpCodes {
		if code == r.StatusCode {
			return false
		}
	}
	return true
}

// Method parse and return all rules from http response
func (r *Response) ExtractLinks() []string {
	var links []string

	req := r.Request

	if len(r.Body) != 0 {
		matches := regexpURL.FindAllStringSubmatch(string(r.Body), -1)
		if len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 0 {
					link := match[1]

					u, _ := url.Parse(link)
					if u.Host == "" {
						link = fmt.Sprintf("%s://%s%s", req.ParsedURL.Scheme, req.ParsedURL.Host, link)
					}
					links = append(links, link)
				}
			}
		}
	}
	return links
}
