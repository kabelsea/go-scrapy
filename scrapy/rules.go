package scrapy

type Rule struct {
	LinkExtractor *LinkExtractor
	Follow        bool
	Handler       func(response *Response)
}
