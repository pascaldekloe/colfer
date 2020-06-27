include common.mk

.PHONY: test
test:
	$(MAKE) -C c test
	$(MAKE) -C ecma test
	$(MAKE) -C go test
	$(MAKE) -C java test
	$(MAKE) -C rpc test
	mvn -f java/maven integration-test

.PHONY: bench
bench:
	$(MAKE) -C c/bench test
	$(MAKE) -C ecma/bench test
	$(MAKE) -C go/bench test
	$(MAKE) -C java/bench test

build:
	GOARCH=amd64 GOOS=linux $(GO) build -o build/colf-linux ./cmd/colf
	GOARCH=amd64 GOOS=darwin $(GO) build -o build/colf-darwin ./cmd/colf
	GOARCH=amd64 GOOS=openbsd $(GO) build -o build/colf-openbsd ./cmd/colf
	GOARCH=amd64 GOOS=windows $(GO) build -o build/colf.exe ./cmd/colf

.PHONY: clean
clean:
	rm -fr build

.PHONY: clean-all
clean-all: clean
	$(MAKE) -C c clean
	$(MAKE) -C c/bench clean
	$(MAKE) -C ecma clean
	$(MAKE) -C ecma/bench clean
	$(MAKE) -C go clean
	$(MAKE) -C go/bench clean
	$(MAKE) -C java clean
	$(MAKE) -C java/bench clean
	$(MAKE) -C rpc clean
	mvn -f java/maven clean
	$(GO) clean -r ./cmd/...
