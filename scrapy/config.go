package scrapy

import (
	"errors"
	"time"
)

var (
	ConcurrentRequests       = 5
	MaxDepth                 = 2
	DownloadTimeout          = time.Duration(30) * time.Second
	DownloadMaxSize    int32 = 1024 * 1024 * 10
	UserAgent                = "go-scrapy/1.0"
	RequestHeaders           = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "en",
	}
	RetryHttpCodes = []int{500, 502, 503, 504, 408, 404}
	RetryTimes     = 2
)

// Scrapy spider configuration
type SpiderConfig struct {
	Debug          bool
	Name           string
	AllowedDomains []string
	StartUrls      []string
	Rules          []Rule

	MaxDepth int

	DownloadMiddlewares map[DownloadMiddleware]int

	// Concurrent settings
	ConcurrentRequests int

	// Download settings
	DownloadTimeout time.Duration
	DownloadMaxSize int32

	// Requests params
	UserAgent      string
	RequestHeaders map[string]string

	// Attempt settings
	RetryEnabled   bool
	RetryHttpCodes []int
	RetryTimes     int
}

// Load default value into spider configuration
func (c *SpiderConfig) Default() {
	if c.MaxDepth == 0 {
		c.MaxDepth = MaxDepth
	}

	if c.ConcurrentRequests == 0 {
		c.ConcurrentRequests = ConcurrentRequests
	}

	if c.UserAgent == "" {
		c.UserAgent = UserAgent
	}

	if c.DownloadTimeout == 0 {
		c.DownloadTimeout = DownloadTimeout
	}

	if c.DownloadMaxSize == 0 {
		c.DownloadMaxSize = DownloadMaxSize
	}

	if len(c.RequestHeaders) == 0 {
		c.RequestHeaders = RequestHeaders
	}

	if c.RetryEnabled {
		if len(c.RetryHttpCodes) == 0 {
			c.RetryHttpCodes = RetryHttpCodes
		}

		if c.RetryTimes == 0 {
			c.RetryTimes = RetryTimes
		}
	}
}

// Validate spider configuration
func (c *SpiderConfig) Validate() error {
	if len(c.Rules) == 0 {
		return errors.New("not found rules in configuration")
	}
	return nil
}
