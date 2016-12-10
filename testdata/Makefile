.PHONY: clean init test run-fuzz

COLF?=${GOPATH}/bin/colf

test: build
	go vet
	go fmt
	go test -v

	java -cp build/classes testdata.test

testcover-go: Colfer.go
	go test -v -coverprofile testcover-go github.com/pascaldekloe/colfer/testdata

clean:
	go clean
	go clean ./build/go/...

	rm -fr build
	rm -f Colfer.go Colfer.js O.java package-info.java
	rm -f testcover-go
	rm -f fuzz-go

init:
	go get github.com/pascaldekloe/colfer/cmd/colf

build: init Colfer.go Colfer.js O.java
	$(COLF) -b ../../../.. -p github.com/pascaldekloe/colfer/testdata/build/go/break go break*.colf
	$(COLF) -b build/java -p break java break*.colf
	$(COLF) -b build/js -p break js break*.colf

	go build
	go build ./build/go/...

	mkdir -p build/classes
	javac -d build/classes *.java build/java/break_/*/*.java
	javadoc -d build/javadoc -sourcepath build/java -subpackages . > /dev/null

Colfer.go: init
	$(COLF) -b .. go test.colf

Colfer.js: init
	$(COLF) js test.colf

O.java: init
	$(COLF) -b .. java test.colf

package-info.java: O.java

fuzz-go: build
	rm corpus/seed*
	go test -run FuzzSeed

	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build
	go-fuzz-build -o fuzz-go github.com/pascaldekloe/colfer/testdata

run-fuzz: fuzz-go
	go-fuzz -bin fuzz-go -workdir .
