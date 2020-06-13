include ../../common.mk

NODE ?= node

.PHONY: test
test: build
	$(NODE) bench.js
	$(NODE) --version

build:
	$(COLF) -b build -s 1024 js ../../testdata/bench/scheme.colf

.PHONY: clean
clean:
	rm -fr build
