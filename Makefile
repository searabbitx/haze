SRCS := $(shell find . -name '*.go')
GOFLAGS=-trimpath -ldflags "-s -w"

build: haze

haze: $(SRCS)
	go build $(GOFLAGS) .

all: build/haze_linux_amd64 build/haze_linux_386

build/haze_%: $(SRCS)
	env GOOS=$(word 1,$(subst _, ,$*)) GOARCH=$(word 2,$(subst _, ,$*)) go build -o $@ $(GOFLAGS) .

.PHONY: format
format:
	gofmt -w */*go

.PHONY: test
test:
	go test ./... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
