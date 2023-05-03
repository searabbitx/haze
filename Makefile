SRCS := $(shell find . -name '*.go')

build: haze

haze: $(SRCS)
	go build .

.PHONY: format
format:
	gofmt -w */*go

.PHONY: test
test:
	go test ./... -v
