include ../common.mk

.PHONY: test
test: gen build node_modules
	node_modules/.bin/qunit test.js

build:
	$(COLF) -b build JavaScript ../testdata/break*.colf

gen:
	$(COLF) -b gen JavaScript ../testdata/test.colf

node_modules:
	npm install qunit

.PHONY: clean
clean:
	rm -fr gen build
