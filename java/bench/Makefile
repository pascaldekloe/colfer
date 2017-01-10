include ../../common.mk

.PHONY: test
test: build
	java -cp build/classes bench

build: install
	$(COLF) -b build/java Java ../../testdata/bench/scheme.colf

	mkdir -p build/classes
	javac -d build/classes bench.java build/java/bench/*.java

.PHONY: clean
clean:
	rm -fr build
