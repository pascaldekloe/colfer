include ../common.mk

JAVA ?= java
JAVAC ?= javac
JAVADOC ?= javadoc

.PHONY: test
test: test.class breaktest
	$(JAVA) test

breaktest: ../testdata/break*.colf ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -b $@ -x java.util.ArrayList Java ../testdata/break*.colf
	$(JAVAC) $@/*/*.java
	$(JAVADOC) -d $@ $@/*/*.java
	touch $@

gen: ../testdata/test.colf ../testdata/test-java.tags ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -t ../testdata/test-java.tags Java ../testdata/test.colf
	$(JAVAC) $@/*.java
	touch $@

gen-doc: gen
	$(JAVADOC) -d $@ $?/*.java

test.class: test.java gen
	$(JAVAC) test.java

.PHONY: clean
clean:
	rm -f *.class gen/*.class
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -fr gen gen-doc
