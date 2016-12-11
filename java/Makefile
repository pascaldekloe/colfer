.PHONY: bench clean test

COLF?=../colf

test: build
	java -cp build/classes test

bench: build
	java -cp build/classes bench

clean:
	rm -fr build gen

build: gen
	$(COLF) -b build/java java ../testdata/bench/scheme.colf
	$(COLF) -b build/java -p break java ../testdata/break*.colf

	mkdir -p build/classes
	javac -d build/classes test.java gen/*.java
	javac -d build/classes bench.java build/java/bench/*.java
	javac -d build/classes build/java/break_/*/*.java

	javadoc -d build/javadoc -sourcepath build/java -subpackages . > /dev/null

gen:
	$(COLF) java ../testdata/test.colf
