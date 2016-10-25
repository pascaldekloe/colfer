clean:
	go clean ./...
	rm -f cmd/colf/colf
	rm -f testdata/*.class testdata/bench/*.class
	rm -f testdata/*-fuzz.zip

build:
	go generate
	go get github.com/pascaldekloe/colfer/cmd/colf
	javac testdata/*.java testdata/bench/*.java

test: build
	go vet ./...
	go test
	java -cp . testdata.test

bench:
	go test -run none -bench .
	java -cp . testdata.bench.bench

dist: test clean
	go fmt

fuzzing:
	rm testdata/corpus/seed*
	go test -run FuzzSeed
	go-fuzz-build -o testdata/go-fuzz.zip github.com/pascaldekloe/colfer/testdata
	go-fuzz -bin testdata/go-fuzz.zip -workdir testdata
