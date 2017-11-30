package scrapy

import (
	"net/url"
	"testing"
)

func MockLinkExtractor() LinkExtractor {
	return LinkExtractor{
		Allow: []string{
			`search\.php\?.+`,
			`detail\.php\?.+`,
		},
		Deny: []string{
			`catalog\.html\?.+`,
			`map\.html\?.+`,
		},
		AllowDomains: []string{
			`^test\.ru`,
			`^.*\.test\.ru`,
		},
		DenyDomains: []string{
			`^example\.com`,
			`^[a-z]{2,3}\.com`,
		},
	}
}

func ParseUrl(link string) *url.URL {
	u, _ := url.Parse(link)
	return u
}

func TestLinkExtractor_Compile(t *testing.T) {
	lex := MockLinkExtractor()
	lex.Compile()

	for k, v := range lex.Compiled {
		if len(v) != 2 {
			t.Error(
				"Wrong len of compiled regular expression for", k,
				"expected", 2,
				"got", len(v),
			)
		}
	}
}

func TestLinkExtractor_Match(t *testing.T) {
	lex := MockLinkExtractor()
	lex.Compile()

	u := "https://test.ru/search.php?q=search"
	if !lex.Match(ParseUrl(u)) {
		t.Error(
			"Link not matched", u,
			"expected", true,
			"got", false,
		)
	}

	u = "https://test.ru/notpresent.html"
	if lex.Match(ParseUrl(u)) {
		t.Error(
			"Link matched", u,
			"expected", false,
			"got", true,
		)
	}

	u = "https://notpresent.com/search.php?q=search"
	if lex.Match(ParseUrl(u)) {
		t.Error(
			"Link matched", u,
			"expected", false,
			"got", true,
		)
	}

	u = "https://example.com/search.php?q=search"
	if lex.Match(ParseUrl(u)) {
		t.Error(
			"Link matched", u,
			"expected", false,
			"got", true,
		)
	}

	u = "https://example.com/search.php?q=search"
	if lex.Match(ParseUrl(u)) {
		t.Error(
			"Link matched", u,
			"expected", false,
			"got", true,
		)
	}

	u = "http://blabla.com/some.php?bla=bla"
	if lex.Match(ParseUrl(u)) {
		t.Error(
			"Link matched", u,
			"expected", false,
			"got", true,
		)
	}
}
