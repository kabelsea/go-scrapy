package scrapy

type Handler func(response *Response)

type Rule struct {
	LinkExtractor *LinkExtractor
	Follow        bool
	Handler       Handler
}

type HandlerChannel chan Handler

// Method call handler for processing data
//func (r *Rule) ProcessItem(resp *Response) {
//	if r.Handler != nil {
//		r.Handler(resp)
//	}
//}
