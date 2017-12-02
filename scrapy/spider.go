package scrapy

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

type Spider struct {
	Config        *SpiderConfig
	Stats         *Stats
	ProcessedUrls map[string]int
}

// Create new spider instance
func NewSpider(config *SpiderConfig) (*Spider, error) {
	config.Default()

	// Setup logger
	logger.SetOutput(os.Stdout)
	if config.Debug {
		logger.SetLevel(logger.DebugLevel)
	} else {
		logger.SetLevel(logger.InfoLevel)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Compile and validate all LinkExtractor regexp
	for _, r := range config.Rules {
		err := r.LinkExtractor.Compile()
		if err != nil {
			return nil, err
		}
	}

	spider := &Spider{
		Config:        config,
		ProcessedUrls: map[string]int{},
	}
	return spider, nil
}

// Run spider
func (s *Spider) Wait() {
	logger.Info("Crawler started")

	semaphore := make(chan bool, s.Config.ConcurrentRequests)
	done := make(chan bool)

	requests := make(RequestChannel)
	responses := make(ResponseChannel)

	// Initialize requests channel from start urls array
	go func() {
		for _, url := range s.Config.StartUrls {
			s.ProcessedUrls[url] = 1
			requests <- *NewRequest(url, s.Config)
		}
	}()

	for {
		select {
		case req := <-requests:
			go func(request Request) {
				semaphore <- true
				defer func() { <-semaphore }()

				logger.Infof("%s started", req.Url)

				resp, err := request.Process()
				if err != nil {
					if req.CanRetry() {
						req.Attempt++
						requests <- req
						logger.Infof("%s retried", req.Url)
					}
				} else {
					responses <- resp
					logger.Infof("%s proceed", req.Url)
				}
			}(req)
		case resp := <-responses:
			handlers := resp.Handlers()
			if len(handlers) != 0 {
				for _, h := range handlers {
					h(resp)
				}
			}

			// Extracts all links from http response and put it into
			//  requests channel if does not processed
			for _, link := range resp.ExtractLinks() {
				req := NewRequest(link, s.Config)
				req.Depth++

				if (req.CanFollow() || req.CanParse()) && !s.CheckProcessUrl(link) {
					go func(req *Request) {
						requests <- *req
					}(req)
				}

				s.ProcessedUrls[link] = 1
			}
		case <-done:
			return
		}
	}
}

// Method check url in processed urls list, if exists return true
func (s *Spider) CheckProcessUrl(url string) bool {
	if _, ok := s.ProcessedUrls[url]; !ok {
		return false
	}
	return true
}
