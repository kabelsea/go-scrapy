package scrapy

import (
	"log"
)

type Spider struct {
	config *SpiderConfig
}

type SpiderRequest struct {
}

type SpiderResponse struct {
}

// Create new spider instance
func NewSpider(config *SpiderConfig) (*Spider, error) {
	config.LoadDefault()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	spider := &Spider{
		config: config,
	}
	return spider, nil
}

// Run spider
func (s *Spider) Wait() {
	semaphore := make(chan bool, s.config.ConcurrentRequests)
	requests := make(chan Request)

	// Initial requests channel from start urls slice
	go func() {
		for _, url := range s.config.StartUrls {
			requests <- Request{
				Url:     url,
				Headers: s.config.RequestHeaders,
				Meta:    make(map[string]string),
			}
		}
	}()

	for req := range requests {
		semaphore <- true

		go func(req Request) {
			defer func() { <-semaphore }()

			if ok := req.Download(); ok {
				log.Println(string(req.Url))
			}
		}(req)
	}

	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}
}
