SRCS := $(shell find . -name '*.go')

build: haze

haze: $(SRCS)
	go build .

.PHONY: format
format:
	gofmt -w */*go

.PHONY: test
test:
	go test ./... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
