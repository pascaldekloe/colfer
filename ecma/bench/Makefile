include ../../common.mk

NODE ?= node

bench.txt: Colfer.js bench.js benchmark.js lodash.js
	$(NODE) bench.js > $@
	echo >> $@
	echo "Node version: \c" >> $@
	$(NODE) --version >> $@

Colfer.js: ../../testdata/bench/scheme.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) -s 1024 JS ../../testdata/bench/scheme.colf

.PHONY: clean
clean:
	rm -f bench.txt Colfer.js
