.PHONY: clean test

COLF?=../colf

test: build node_modules
	node run.js

build: Colfer.js
	$(COLF) -b build js ../testdata/break*.colf

Colfer.js:
	$(COLF) js ../testdata/test.colf

node_modules:
	npm install qunit

clean:
	rm -fr build Colfer.js
