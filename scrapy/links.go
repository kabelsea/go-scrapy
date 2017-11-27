package scrapy

type LinkExtractor struct {
	Allow          []string
	Deny           []string
	AllowDomains   []string
	DenyDomains    []string
	DenyExtensions []string
}

// Compile all regular expressions in Rule
func (l *LinkExtractor) MakeCompile() error {
	return nil
}
