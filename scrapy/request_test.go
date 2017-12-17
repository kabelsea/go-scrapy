package scrapy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	config = SpiderConfig{
		UserAgent: "go-scrapy/test",
		RequestHeaders: http.Header{
			"Test": {"Test"},
		},
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

func TestNewRequest(t *testing.T) {
	u := "https://example.com/path/?search=name#fragment"
	up, _ := url.Parse(u)
	req, _ := NewRequest(u, &config)

	if req.URL == nil || !reflect.DeepEqual(req.URL, up) {
		t.Error(
			"Request object has wrong parsed url",
			"expected", up,
			"got", req.URL,
		)
	}

	if val := req.Headers.Get("User-Agent"); val != config.UserAgent {
		t.Error(
			"Not passed HTTP Header `User-Agent` or wrong value",
			"expected", config.UserAgent,
			"got", val,
		)
	}
}

func TestRequest_CanFollow(t *testing.T) {
	config.Default()
	config.Rules[0].LinkExtractor.Compile()

	req, _ := NewRequest("https://test.com/search.php?q=name", &config)
	if !req.CanFollow() {
		t.Error(
			"Can follow url", req.URL,
			"expected", true,
			"got", false,
		)
	}

	req, _ = NewRequest("https://test.com/search.php?q=name", &config)
	req.Depth = config.MaxDepth + 1
	if req.CanFollow() {
		t.Error(
			"Can not follow url, because max depth more then max value in configuration", req.URL,
			"expected", false,
			"got", true,
		)
	}
	req.Depth = 0

	config.Rules[0].Follow = false
	req, _ = NewRequest("https://test.com/search.php?q=name", &config)
	if req.CanFollow() {
		t.Error(
			"Can not follow url, because follow property has false value", req.URL,
			"expected", false,
			"got", true,
		)
	}
}

func TestRequest_CanParse(t *testing.T) {
	config.Default()
	config.Rules[0].LinkExtractor.Compile()

	req, _ := NewRequest("https://test.com/search.php?q=name", &config)
	if !req.CanParse() {
		t.Error(
			"Can parse url", req.URL,
			"expected", true,
			"got", false,
		)
	}

	config.Rules[0].Handler = nil
	req, _ = NewRequest("https://test.com/search.php?q=name", &config)
	if req.CanParse() {
		t.Error(
			"Can not parse url, because rule handler does not exist", req.URL,
			"expected", false,
			"got", true,
		)
	}
}

func TestRequest_Process(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]map[string]string{}
		data["headers"] = map[string]string{}
		data["cookies"] = map[string]string{}

		for k, v := range r.Header {
			data["headers"][k] = v[0]
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}))
	defer ts.Close()

	req, _ := NewRequest(ts.URL, &config)
	resp, _ := req.Process()

	body := map[string]map[string]string{}
	json.Unmarshal(resp.Body, &body)

	if val, ok := body["headers"]["Test"]; !ok || val != "Test" {
		t.Error(
			"Not passed HTTP Header `Test` or wrong value",
			"expected", "Test",
			"got", val,
		)
	}

	if val, ok := body["headers"]["User-Agent"]; !ok || val != config.UserAgent {
		t.Error(
			"Not passed HTTP Header `User-Agent` or wrong value",
			"expected", config.UserAgent,
			"got", val,
		)
	}

	if resp.StatusCode != 200 {
		t.Error(
			"Wrong HTTP status code in Response",
			"expected", 200,
			"got", resp.StatusCode,
		)
	}

	if resp.Body == nil {
		t.Error("Response body can not be empty")
	}
}
