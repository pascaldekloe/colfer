include ../common.mk

.PHONY: test
test: internal build
	$(GO) test -v -coverprofile build/coverage

.PHONY: test
bench: internal build
	$(GO) test -run none -bench .

internal: install
	$(COLF) Go

build: install
	$(COLF) -b build Go ../testdata/test.colf

.PHONY: clean
clean:
	rm -fr internal build
