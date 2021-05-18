include common.mk

.PHONY: test
test:
	$(GO) test -v

	$(MAKE) -C c test
	$(MAKE) -C dart test
	$(MAKE) -C ecma test
	$(MAKE) -C go test
	$(MAKE) -C java test
	$(MAKE) -C java/maven target
	$(MAKE) -C rpc test

.PHONY: clean
clean:
	$(GO) clean -r ./cmd/...
	$(MAKE) -C c clean
	$(MAKE) -C c/bench clean
	$(MAKE) -C c/fuzz clean
	$(MAKE) -C dart clean
	$(MAKE) -C ecma clean
	$(MAKE) -C ecma/bench clean
	$(MAKE) -C go clean
	$(MAKE) -C go/bench clean
	$(MAKE) -C go/fuzz clean
	$(MAKE) -C java clean
	$(MAKE) -C java/bench clean
	$(MAKE) -C java/maven clean
	$(MAKE) -C rpc clean
