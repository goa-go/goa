GO ?= go

run:
	$(GO) build ./server/app.go
	./app.exe

testServer:
	$(GO) test ./server -v

install:
	$(GO) get ./...

.PHONY: test
test:
	$(GO) test ./test -v
