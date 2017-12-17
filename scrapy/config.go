package scrapy

import (
	"errors"
	"net/http"
	"time"
)

// Varibales with default spider configuration
var (
	ConcurrentRequests       = 5
	MaxDepth                 = 2
	DownloadTimeout          = time.Duration(30) * time.Second
	DownloadMaxSize    int32 = 1024 * 1024 * 10
	RequestHeaders           = &http.Header{}
	UserAgent                = "go-scrapy/1.0"
)

// Scrapy spider configuration struct
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
	RequestHeaders http.Header

	// Spider stats collector
	Stats SpiderStats
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

	if c.RequestHeaders == nil {
		c.RequestHeaders = *RequestHeaders
	}

	if c.Stats == nil {
		c.Stats = NewStats()
	}
}

// Validate spider configuration
func (c *SpiderConfig) Validate() error {
	if len(c.Rules) == 0 {
		return errors.New("not found rules in configuration")
	}
	return nil
}
