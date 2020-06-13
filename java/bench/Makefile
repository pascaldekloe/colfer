include ../../common.mk

JAVA ?= java
JAVAC ?= javac

.PHONY: test
test: build
	$(JAVA) -cp build/classes bench

build:
	$(COLF) -b build/java Java ../../testdata/bench/scheme.colf

	mkdir -p build/classes
	$(JAVAC) -d build/classes bench.java build/java/bench/*.java

.PHONY: clean
clean:
	rm -fr build
