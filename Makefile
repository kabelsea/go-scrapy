SRC_FILES := $(wildcard scrapy/*.go)

deps:
	glide install

test:
	go test $(SRC_FILES) -cover -race -coverprofile=coverage.txt -covermode=atomic -v

.PHONY: deps test
