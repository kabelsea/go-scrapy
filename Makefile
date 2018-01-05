SRC_FILES := $(wildcard scrapy/*.go)

deps:
	glide install

test:
	go test ./scrapy -coverprofile=coverage.txt -v && go tool cover -html=coverage.txt -o coverage.html

.PHONY: deps test
