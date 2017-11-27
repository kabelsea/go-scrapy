package main

import (
	"../../scrapy"
)

func main() {
	// Init spider configuration
	config := &scrapy.SpiderConfig{
		Name:               "HabraBot",
		MaxDepth:           2,
		ConcurrentRequests: 2,
		RetryEnabled:       true,
		StartUrls: []string{
			"https://habrahabr.ru/top/",
			"https://habrahabr.ru/all/",
		},
		Rules: []scrapy.Rule{
			scrapy.Rule{
				LinkExtractor: scrapy.LinkExtractor{
					Allow:        []string{"/post/\\d+/"},
					AllowDomains: []string{"habrahabr.ru"},
					DenyDomains:  []string{"some.name"},
				},
				Follow: true,
			},
			scrapy.Rule{
				LinkExtractor: scrapy.LinkExtractor{
					Allow:        []string{"/post/\\d+/"},
					AllowDomains: []string{"habrahabr.ru"},
					DenyDomains:  []string{"some.name"},
				},
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
