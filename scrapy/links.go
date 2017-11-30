package scrapy

import (
	"net/url"
	"reflect"
	"regexp"
)

// Link extractor struct
//
// Compiled map contains all compiled regular expressions
//      and available by key, equals struct field name.
type LinkExtractor struct {
	Allow        []string
	Deny         []string
	AllowDomains []string
	DenyDomains  []string
	Compiled     map[string][]*regexp.Regexp
}

var (
	compileFields = []string{
		"Allow", "Deny", "AllowDomains", "DenyDomains",
	}
)

// Method compile all regular expressions
func (l *LinkExtractor) Compile() error {
	if len(l.Compiled) == 0 {
		l.Compiled = make(map[string][]*regexp.Regexp)
	}

	for _, act := range compileFields {
		val := reflect.Indirect(reflect.ValueOf(l)).FieldByName(act).Interface().([]string)
		if len(val) != 0 {
			for _, v := range val {
				reg, err := regexp.Compile(v)
				if err != nil {
					return err
				}
				l.Compiled[act] = append(l.Compiled[act], reg)
			}
		}
	}
	return nil
}

// Method checks whether the http url to the link extractor rule is appropriate
func (l *LinkExtractor) Match(u *url.URL) bool {
	matched := false

	// check if domain allow
	regexps, ok := l.Compiled["AllowDomains"]
	if ok {
		for _, r := range regexps {
			if r.MatchString(u.Host) {
				matched = true
			}
		}
	}

	// check if domain deny
	regexps, ok = l.Compiled["DenyDomains"]
	if ok && matched {
		for _, r := range regexps {
			if r.MatchString(u.Host) {
				matched = false
			}
		}
	}

	// check if uri allow
	regexps, ok = l.Compiled["Allow"]
	if ok && matched {
		for _, r := range regexps {
			if r.MatchString(u.RequestURI()) {
				matched = true
			}
		}
	}

	// check if uri deny
	regexps, ok = l.Compiled["Deny"]
	if ok && matched {
		for _, r := range regexps {
			if r.MatchString(u.RequestURI()) {
				matched = false
			}
		}
	}
	return matched
}
