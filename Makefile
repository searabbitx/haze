.PHONY: format
format:
	gofmt -w */*go

.PHONY: test
test:
	go test ./... -v
