package scrapy

import (
	"testing"
)

var (
	config = SpiderConfig{
		Rules: []Rule{
			{
				LinkExtractor: &LinkExtractor{
					Allow:        []string{`search\.php.+`},
					AllowDomains: []string{`^test\.com`},
				},
				Follow:  true,
				Handler: func(resp *Response) {},
			},
		},
	}
)

func TestRequest_CanFollow(t *testing.T) {
	config.LoadDefault()
	config.Rules[0].LinkExtractor.Compile()

	req := NewRequest("https://test.com/search.php?q=name", &config)
	if !req.CanFollow() {
		t.Error(
			"Can follow url", req.Url,
			"expected", true,
			"got", false,
		)
	}

	req = NewRequest("https://test.com/search.php?q=name", &config)
	req.Depth = config.MaxDepth + 1
	if req.CanFollow() {
		t.Error(
			"Can not follow url, because max depth more then max value in configuration", req.Url,
			"expected", false,
			"got", true,
		)
	}
	req.Depth = 0

	config.Rules[0].Follow = false
	req = NewRequest("https://test.com/search.php?q=name", &config)
	if req.CanFollow() {
		t.Error(
			"Can not follow url, because follow property has false value", req.Url,
			"expected", false,
			"got", true,
		)
	}
}

func TestRequest_CanParse(t *testing.T) {
	config.LoadDefault()
	config.Rules[0].LinkExtractor.Compile()

	req := NewRequest("https://test.com/search.php?q=name", &config)
	if !req.CanParse() {
		t.Error(
			"Can parse url", req.Url,
			"expected", true,
			"got", false,
		)
	}

	config.Rules[0].Handler = nil
	req = NewRequest("https://test.com/search.php?q=name", &config)
	if req.CanParse() {
		t.Error(
			"Can not parse url, because rule handler does not exist", req.Url,
			"expected", false,
			"got", true,
		)
	}
}
