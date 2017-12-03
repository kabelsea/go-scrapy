package scrapy

// Spider downloader interface
type DownloadMiddleware interface {
	ProcessRequest(req Request, spider Spider) interface{}
	ProcessResponse(req Request, resp Response, spider Spider) interface{}
	ProcessException(req Request, err error, spider Spider) interface{}
}

// Retry middleware
type RetryMiddleware struct {
	DownloadMiddleware

	// Which HTTP response codes to retry.
	HttpCodes []int

	// Maximum number of times to retry, in addition to the first download.
	Times int
}

// Load default middleware configuration
func (m *RetryMiddleware) Default() {
	if m.HttpCodes == nil {
		m.HttpCodes = []int{500, 502, 503, 504, 408}
	}

	if m.Times == 0 {
		m.Times = 2
	}
}

func (m *RetryMiddleware) ProcessRequest(req Request, spider Spider) interface{} {
	return nil
}
func (m *RetryMiddleware) ProcessResponse(req Request, resp Response, spider Spider) interface{} {
	return nil
}
func (m *RetryMiddleware) ProcessException(req Request, err error, spider Spider) interface{} {
	return nil
}
