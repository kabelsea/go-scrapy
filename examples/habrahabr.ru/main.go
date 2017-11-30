package main

import (
	"log"

	"../../scrapy"
)

func main() {
	// Init spider configuration
	config := &scrapy.SpiderConfig{
		Name:               "HabraBot",
		MaxDepth:           2, // TODO: Needed implementation
		ConcurrentRequests: 2,
		RetryEnabled:       true,
		RetryTimes:         2,
		StartUrls: []string{
			"https://habrahabr.ru/",
		},
		Rules: []scrapy.Rule{
			scrapy.Rule{
				LinkExtractor: scrapy.LinkExtractor{
					Allow:        []string{"/post/\\d+/"},
					AllowDomains: []string{"habrahabr.ru"},
					DenyDomains:  []string{"twitter.com"},
				},
				Follow: true,
			},
			scrapy.Rule{
				LinkExtractor: scrapy.LinkExtractor{
					Allow:        []string{"/post/\\d+/"},
					AllowDomains: []string{"habrahabr.ru"},
				},
				Handler: ProcessItem,
			},
		},
	}

	// Create new spider
	spider, err := scrapy.NewSpider(config)
	if err != nil {
		panic(err)
	}

	// Run spider and wait
	spider.Wait()
}

func ProcessItem(resp *scrapy.Response) {
	log.Println(resp.Url, resp.StatusCode)
}
