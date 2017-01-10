GOPATH?=../../..
COLF?=$(GOPATH)/bin/colf

.PHONY: run
run: clean test

.PHONY: install
install:
	go install github.com/pascaldekloe/colfer/cmd/...
