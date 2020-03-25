GO ?= go
GOPATH != go env GOPATH
COLF = $(GOPATH)/bin/colf

.PHONY: run
run: clean test

.PHONY: install
install:
	$(GO) install github.com/pascaldekloe/colfer/cmd/...
