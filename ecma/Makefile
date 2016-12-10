.PHONY: clean test

COLF?=../colf

test: Colfer.js
	$(COLF) -b build js ../testdata/break*.colf

Colfer.js:
	$(COLF) js ../testdata/test.colf

clean:
	rm -fr build Colfer.js
