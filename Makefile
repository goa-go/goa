GO ?= go

run:
	$(GO) build ./server/app.go
	./app.exe

test:
	$(GO) test ./server -v
