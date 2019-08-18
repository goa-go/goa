GO ?= go

run:
	$(GO) build ./server/app.go
	./app.exe

install:
	$(GO) get ./...

.PHONY: test
test:
	$(GO) test ./test/... -v

test_cover:
	$(GO) test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./,./responser ./test/ ./test/body
