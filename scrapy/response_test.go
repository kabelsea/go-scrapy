package scrapy

import (
	"fmt"
	"reflect"
	"testing"
)

// Mock Response object
func MockResponse() *Response {
	config := &SpiderConfig{}
	config.Default()

	req, _ := NewRequest("http://test.com", config)

	return &Response{
		StatusCode: 200,
		Body:       []byte{},
		Request:    req,
	}
}

func TestResponse_ExtractLinks(t *testing.T) {
	testLinks := []interface{}{
		"http://test.com/bla",
		"http://test.com/bla-bla",
		"http://test.com/bla-bla?some=param",
		"http://test.com/bla-bla-bla#fragment",
	}

	resp := MockResponse()
	resp.Body = []byte(fmt.Sprintf(`
        <div id="block">
            <a href="%s">Link 1</a>
            <a href="%s">Link 2</a>
            <a href="%s">Link 3</a>
            <a href="%s">Link 4</a>
        </div>
    `, testLinks...))

	// Convert interfaces slice to string slice
	tt := make([]string, 0)
	for _, t := range testLinks {
		tt = append(tt, t.(string))
	}

	links := resp.ExtractLinks()

	if !reflect.DeepEqual(links, tt) {
		t.Error(
			"Wrong extracted links",
			"expected", testLinks,
			"got", links,
		)
	}
}
