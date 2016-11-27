.PHONY: clean clean-gen build test regression bench fuzzing

test: build
	go vet ./...
	go test ./...

	colf java testdata/test.colf
	javac -d testdata/build testdata/*.java
	java -cp testdata/build testdata.test
	javadoc -d testdata/build/javadoc testdata > /dev/null

	colf -b testdata js testdata/test.colf

clean:
	go clean ./...
	rm -fr dist
	rm -f cmd/colf/colf
	rm -f testdata/*-fuzz.zip
	rm -fr testdata/build testdata/bench/build
	mkdir testdata/build testdata/bench/build

clean-gen:
	rm -f testdata/Colfer.* testdata/O.java

build:
	go generate
	go get github.com/pascaldekloe/colfer/cmd/colf

regression: build
	colf -b ../../.. -p github.com/pascaldekloe/colfer/testdata/build/break go testdata/break*.colf
	go build ./testdata/build/break/...

	colf -b testdata/build/break java testdata/break*.colf
	javac testdata/build/break/*/*.java

	colf -b testdata/build/break js testdata/break*.colf

bench: build
	go generate ./testdata/bench
	go test -bench . ./testdata/bench

	colf -b testdata/bench/build java testdata/bench/scheme.colf
	javac -d testdata/bench/build -sourcepath testdata/bench/build testdata/bench/bench.java
	java -cp testdata/bench/build testdata.bench.bench

dist: clean-gen test regression clean
	go fmt
	mkdir -p dist
	GOARCH=amd64 GOOS=linux go build -o dist/colf-linux ./cmd/colf
	GOARCH=amd64 GOOS=darwin go build -o dist/colf-darwin ./cmd/colf
	GOARCH=amd64 GOOS=windows go build -o dist/colf.exe ./cmd/colf

fuzzing:
	rm testdata/corpus/seed*
	go test -run FuzzSeed
	go-fuzz-build -o testdata/go-fuzz.zip github.com/pascaldekloe/colfer/testdata
	go-fuzz -bin testdata/go-fuzz.zip -workdir testdata
