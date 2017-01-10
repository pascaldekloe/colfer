include ../common.mk

.PHONY: test
test: internal build
	go test -v -coverprofile build/coverage

.PHONY: test
bench: internal build
	go test -run none -bench .

internal: install
	$(COLF) Go

build: install
	$(COLF) -b build Go ../testdata/test.colf

.PHONY: clean
clean:
	rm -fr internal build
