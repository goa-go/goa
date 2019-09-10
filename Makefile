GO ?= go

run:
	$(GO) build ./server/app.go
	./app.exe

install:
	$(GO) get ./...

.PHONY: test
test:
	$(GO) test ./... -v

test_cover:
	$(GO) test -race -coverprofile=coverage.txt -covermode=atomic ./...

fmt:
	$(GO) fmt ./...
