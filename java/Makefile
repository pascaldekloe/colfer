include ../common.mk

JAVA ?= java
JAVAC ?= javac
JAVADOC ?= javadoc

.PHONY: test
test: build
	$(JAVA) -cp build/classes test

gen:
	$(COLF) Java ../testdata/test.colf

build: gen
	$(COLF) -b build/java -p break -x java/util/ArrayList Java ../testdata/break*.colf

	mkdir -p build/classes
	$(JAVAC) -d build/classes test.java gen/*.java
	$(JAVAC) -d build/classes build/java/break_/*/*.java

	$(JAVADOC) -d build/javadoc build/java/break_/*/*.java > /dev/null

.PHONY: clean
clean:
	rm -fr gen build
