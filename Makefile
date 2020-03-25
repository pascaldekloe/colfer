include common.mk

MAKE ?= make
COLF = $(GOPATH)/bin/colf

.PHONY: dist
dist: clean test build
	$(GO) fmt ./...
	$(GO) vet ./...

.PHONY: test
test: install
	$(MAKE) -C c
	$(MAKE) -C ecma
	$(MAKE) -C go
	$(MAKE) -C java
	$(MAKE) -C rpc
	# Fails on Travis CI: mvn -f java/maven integration-test

.PHONY: bench
bench: install
	$(MAKE) -C c/bench
	$(MAKE) -C go/bench
	$(MAKE) -C java/bench

build:
	GOARCH=amd64 GOOS=linux $(GO) build -o build/colf-linux ./cmd/colf
	GOARCH=amd64 GOOS=darwin $(GO) build -o build/colf-darwin ./cmd/colf
	GOARCH=amd64 GOOS=openbsd $(GO) build -o build/colf-openbsd ./cmd/colf
	GOARCH=amd64 GOOS=windows $(GO) build -o build/colf.exe ./cmd/colf

.PHONY: clean
clean:
	$(GO) clean -i ./cmd/...
	rm -fr build
	# Fails on Travis CI:  mvn -f java/maven clean
