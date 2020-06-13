GO ?= go
GOPATH != $(GO) env GOPATH
COLF = $(GO) run github.com/pascaldekloe/colfer/cmd/colf

MAKE ?= make
