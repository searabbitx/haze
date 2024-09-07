SRCS := $(shell find . -name '*.go')
GOFLAGS=-trimpath -ldflags "-s -w" -buildvcs=false

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

all: build/dist/haze_linux_amd64.tar.gz \
	build/dist/haze_linux_386.tar.gz \
	build/dist/haze_linux_arm.tar.gz \
	build/dist/haze_windows_amd64.exe.tar.gz \
	build/dist/haze_windows_386.exe.tar.gz \
	build/dist/haze_macOS_amd64.tar.gz \
	build/dist/haze_macOS_arm64.tar.gz

build/dist/haze_%.tar.gz: build/bin/haze_%
	mkdir -p build/dist/
	tar --transform 's_build/bin/__;s_haze.*_haze_' \
		--owner haze:1000 --group haze:1000 \
		-czvf $@ README.md LICENSE $<

build/bin/haze_%: $(SRCS)
	env GOOS=$(call goos,$@) GOARCH=$(call goarch,$@) go build -o $@ $(GOFLAGS) .

.PHONY: format
format:
	gofmt -s -w */*go

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	for s in $$(go list ./...); do if ! go test -failfast -v $$s; then break; fi; done | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
