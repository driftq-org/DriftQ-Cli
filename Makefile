.PHONY: build install test
build:
	go build ./...
install:
	go install ./cmd/driftq
test:
	go test ./...
