.PHONY: bench clean

COLF?=${GOPATH}/bin/colf

bench: build
	go test -bench .

	java -cp build/classes testdata.bench.bench

build:
	go generate -x

	$(COLF) -b build/java java
	mkdir -p build/classes
	javac -d build/classes -sourcepath build/java *.java

clean:
	rm -fr build
