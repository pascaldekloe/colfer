include ../../common.mk

JAVA ?= java
JAVAC ?= javac

.PHONY: test
bench.txt: bench.class
	$(JAVA) bench > $@
	echo >> $@
	$(JAVA) -version 2>> $@

bench: ../../testdata/bench/scheme.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) Java ../../testdata/bench/scheme.colf
	touch $@

bench.class: bench
	$(JAVAC) bench.java

.PHONY: clean
clean:
	rm -f bench.txt bench.class
	rm -fr bench
