include ../common.mk

.PHONY: test
test: build
	java -cp build/classes test

gen: install
	$(COLF) Java ../testdata/test.colf

build: gen install
	$(COLF) -b build/java -p break -x java/util/ArrayList Java ../testdata/break*.colf

	mkdir -p build/classes
	javac -d build/classes test.java gen/*.java
	javac -d build/classes build/java/break_/*/*.java

	javadoc -d build/javadoc build/java/break_/*/*.java > /dev/null

.PHONY: clean
clean:
	rm -fr gen build
