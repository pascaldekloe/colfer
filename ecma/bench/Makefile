include ../../common.mk

NODE ?= node

bench.txt: Colfer.js bench.js ../node_modules/benchmark
	$(NODE) bench.js > $@
	echo >> $@
	echo "Node version: \c" >> $@
	$(NODE) --version >> $@

Colfer.js: ../../testdata/bench/scheme.colf ../../*.go
	$(COLF) -s 1024 JS ../../testdata/bench/scheme.colf

../node_modules/benchmark:
	npm install benchmark

.PHONY: clean
clean:
	rm -f bench.txt Colfer.js
