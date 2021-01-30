include ../common.mk

NODE ?= node
NPM ?= npm

.PHONY: test
test: test.js Colfer.js node_modules/.bin/qunit breaktest
	node_modules/.bin/qunit test.js

breaktest: ../testdata/break*.colf ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ JavaScript ../testdata/break*.colf
	$(NODE) --check $@/*.js
	touch $@

Colfer.js: ../testdata/test.colf ../*.go ../cmd/colf/*.go
	$(COLF) JavaScript ../testdata/test.colf
	$(NODE) --check $@

node_modules/.bin/qunit:
	$(NPM) install qunit

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -f Colfer.js
	rm -fr node_modules
