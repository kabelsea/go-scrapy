package scrapy

import (
	"log"
	"os"
	"sync"

	logger "github.com/sirupsen/logrus"
)

type Spider struct {
	Config        *SpiderConfig
	ProcessedUrls map[string]int
	Mutex         *sync.Mutex
}

// Create new spider instance
func NewSpider(config *SpiderConfig) (*Spider, error) {
	config.LoadDefault()

	// Setup logger
	logger.SetOutput(os.Stdout)
	if config.Debug {
		logger.SetLevel(logger.DebugLevel)
	}

	if err := config.Validate(); err != nil {
		return nil, err
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
			requests <- *MakeRequest(url, s.Config)
		}
	}()

	for {
		select {
		case req := <-requests:
			logger.Infof("%s started", req.Url)

			go func(request Request) {
				semaphore <- true
				defer func() { <-semaphore }()
				responses <- request.Download()
				logger.Infof("%s proceed", req.Url)
			}(req)
		case resp := <-responses:
			log.Println(resp.Request.Url, resp.StatusCode())

			if !resp.Success() {
				req := resp.Request

				// Increase request attempt
				if req.Attempt < s.Config.RetryTimes {
					req.Attempt++

					logger.Infof("%s retried", req.Url)

					go func(request *Request) {
						requests <- *request
					}(req)
				}
			} else {
				for _, url := range resp.ExtractLinks() {
					if !s.CheckProcessUrl(url) {
						s.ProcessedUrls[url] = 1

						go func(link string) {
							requests <- *MakeRequest(link, s.Config)
						}(url)
					}
				}
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
