GO ?= go
GOPATH != go env GOPATH
COLF = go run github.com/pascaldekloe/colfer/cmd/colf

.PHONY: run
run: clean test

.PHONY: install
install:
	$(GO) install github.com/pascaldekloe/colfer/cmd/...
