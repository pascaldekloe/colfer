include ../common.mk

NODE ?= npx ts-node
TSC ?= npx tsc
NPM ?= npm

.PHONY: test
test: test.ts Colfer.ts node_modules/.bin/qunit breaktest
	node_modules/.bin/qunit --require ts-node/register test.ts

breaktest: ../testdata/break*.colf ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ ts ../testdata/break*.colf
	$(TSC)
	touch $@

Colfer.ts: ../testdata/test.colf ../*.go ../cmd/colf/*.go
	$(COLF) ts ../testdata/test.colf
	$(TSC)

node_modules/.bin/qunit:
	$(NPM) install

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -f Colfer.ts
	rm -fr node_modules
