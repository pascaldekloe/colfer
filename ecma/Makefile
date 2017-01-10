include ../common.mk

.PHONY: test
test: gen build node_modules
	node run.js

build: install
	$(COLF) -b build JavaScript ../testdata/break*.colf

gen: install
	$(COLF) -b gen JavaScript ../testdata/test.colf

node_modules:
	npm install qunit

.PHONY: clean
clean:
	rm -fr gen build
