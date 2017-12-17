package scrapy

import (
	"os"
	"sync/atomic"

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
	var (
		requestCounter  int32 = 0
		responseCounter int32 = 0
	)

	logger.Info("Crawler started")

	semaphore := make(chan bool, s.Config.ConcurrentRequests)
	done := make(chan bool)

	requests := make(RequestChannel, 10)
	responses := make(ResponseChannel, 10)

	// Initialize requests channel from start urls array
	go func() {
		for _, url := range s.Config.StartUrls {
			s.ProcessedUrls[url] = 1
			req, err := NewRequest(url, s.Config)
			if err != nil {
				// TODO: Process request error
				continue
			}
			requests <- req

			atomic.AddInt32(&requestCounter, 1)
		}
	}()

	for {
		select {
		case req := <-requests:
			go func(r *Request) {
				semaphore <- true
				defer func() { <-semaphore }()

				logger.Infof("%s started", r.URL)

				resp, err := r.Process()
				if err != nil {
					logger.Infof("%s error", r.URL)
				} else {
					responses <- resp
					logger.Infof("%s proceed", r.URL)
				}
			}(req)
		case resp := <-responses:
			atomic.AddInt32(&responseCounter, 1)

			handlers := resp.Handlers()
			if len(handlers) != 0 {
				for _, h := range handlers {
					h(resp)
				}
			}

			// Extracts all links from http response and put it into
			//  requests channel if does not processed
			for _, link := range resp.ExtractLinks() {
				req, err := NewRequest(link, s.Config)
				if err != nil {
					// TODO: Process request error
					continue
				}

				req.Depth++

				if (req.CanFollow() || req.CanParse()) && !s.CheckProcessUrl(link) {
					go func(r *Request) {
						requests <- r
					}(req)

					atomic.AddInt32(&requestCounter, 1)
				}
				s.ProcessedUrls[link] = 1
			}

			if requestCounter == responseCounter {
				go func() { done <- true }()
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
