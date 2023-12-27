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

seal: ../testdata/seal.colf ../testdata/seal-java.tags ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -t ../testdata/seal-java.tags Java ../testdata/seal.colf
	$(JAVAC) $@/*.java
	touch $@

seal-doc: seal
	$(JAVADOC) -d $@ $?/*.java

test.class: test.java seal
	$(JAVAC) test.java

.PHONY: clean
clean:
	rm -f *.class
	rm -fr seal seal-doc
	rm -fr breaktest
