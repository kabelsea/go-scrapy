SRC_FILES := $(wildcard scrapy/*.go)

deps:
	glide install

test:
	go test $(SRC_FILES)

.PHONY: deps test
