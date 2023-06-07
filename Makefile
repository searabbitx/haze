SRCS := $(shell find . -name '*.go')
GOFLAGS=-trimpath -ldflags "-s -w"

define goos
$(word 2,\
	$(subst _, ,\
	$(subst .exe, ,\
	$(subst macOS,darwin,$(1)))))
endef

define goarch
$(word 3,\
	$(subst _, ,\
	$(subst .exe, ,\
	$(subst macOS,darwin,$(1)))))
endef

build: haze

haze: $(SRCS)
	go build $(GOFLAGS) .

all: build/haze_linux_amd64 build/haze_linux_386 build/haze_linux_arm \
	build/haze_windows_amd64.exe build/haze_windows_386.exe \
	build/haze_macOS_amd64 build/haze_macOS_arm64

build/haze_%: $(SRCS)
	env GOOS=$(call goos,$@) GOARCH=$(call goarch,$@) go build -o $@ $(GOFLAGS) .

.PHONY: format
format:
	gofmt -w */*go

.PHONY: test
test:
	go test ./... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
