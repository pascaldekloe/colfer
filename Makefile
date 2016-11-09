clean:
	go clean ./...
	rm -f cmd/colf/colf
	rm -f testdata/*-fuzz.zip
	rm -fr testdata/build testdata/bench/build
	mkdir testdata/build testdata/bench/build

clean-gen:
	rm -f testdata/Colfer.* testdata/O.java

build:
	go generate
	go get github.com/pascaldekloe/colfer/cmd/colf

test: build
	go vet ./...
	go test ./...

	colf java testdata/test.colf
	javac -d testdata/build testdata/*.java
	java -cp testdata/build testdata.test

	colf -b testdata js testdata/test.colf

regression: build
	colf -b ../../.. -p github.com/pascaldekloe/colfer/testdata/build/break go testdata/break*.colf
	go test ./testdata/build/break/...

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

fuzzing:
	rm testdata/corpus/seed*
	go test -run FuzzSeed
	go-fuzz-build -o testdata/go-fuzz.zip github.com/pascaldekloe/colfer/testdata
	go-fuzz -bin testdata/go-fuzz.zip -workdir testdata
