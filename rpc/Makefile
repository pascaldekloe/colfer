include ../common.mk

.PHONY: test
test: internal build
	$(GO) test -v -coverprofile build/coverage

.PHONY: bench
bench: internal build
	$(GO) test -run none -bench .

internal:
	$(COLF) Go

build:
	$(COLF) -b build Go ../testdata/test.colf

.PHONY: clean
clean:
	rm -fr internal build
