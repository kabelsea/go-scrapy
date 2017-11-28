package scrapy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var hrefTokens = []string{"a", "area"}

// Spider response object
type Response struct {
	Request      *Request
	HttpResponse *http.Response
}

type ResponseChannel chan *Response

// Make Response object, initialize with default values
func MakeResponse(req *Request) *Response {
	return &Response{Request: req}
}

func (r *Response) Body() []byte {
	body, err := ioutil.ReadAll(r.HttpResponse.Body)
	if err != nil {
		logger.Error(err)
		return []byte{}
	}
	return body
}

// Method return http response status code, int
func (r *Response) StatusCode() int {
	return r.HttpResponse.StatusCode
}

// Return Response status code
func (r *Response) Success() bool {
	for _, code := range r.Request.Config.RetryHttpCodes {
		if code == r.HttpResponse.StatusCode {
			return false
		}
	}
	return true
}

// Method parse and return all rules from http response
func (r *Response) ExtractLinks() []string {
	var links []string

	body := string(r.Body())
	req := r.Request

	if body != "" {
		z := html.NewTokenizer(strings.NewReader(body))

		for {
			tt := z.Next()

			switch {
			case tt == html.ErrorToken:
				return links
			case tt == html.StartTagToken:
				t := z.Token()

				for _, name := range hrefTokens {
					if t.Data == name {
						for _, a := range t.Attr {
							if a.Key == "href" && a.Val != "" {
								link := a.Val
								u, _ := url.Parse(link)
								if u.Host == "" {
									link = fmt.Sprintf("%s://%s%s", req.ParsedURL.Scheme, req.ParsedURL.Host, a.Val)
								}
								links = append(links, link)
							}
						}
					}
				}
			}
		}
	}
	return links
}
