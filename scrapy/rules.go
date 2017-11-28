package scrapy

type Rules []*Rule

type Rule struct {
	LinkExtractor LinkExtractor
	Follow        bool
	Handler       func(response *Response) error
}

// Return true if checked link is followed
func (r *Rules) IsFollow(link string) {

}

// Return true if checked link is being processed
func (r *Rules) IsProcessed(link string) {

}
