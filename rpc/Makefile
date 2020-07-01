include ../common.mk

.PHONY: test
test: internal gen
	$(GO) test -v

.PHONY: bench
bench: internal build
	$(GO) test -run match-none -bench .

internal: *.colf ../*.go ../cmd/colf/*.go
	$(COLF) Go
	touch $@

gen: ../testdata/test.colf ../*.go ../cmd/colf/*.go
	$(COLF) Go ../testdata/test.colf
	touch $@

.PHONY: clean
clean:
	rm -fr gen

.PHONY: clean-all
clean-all: clean
	rm -fr internal
