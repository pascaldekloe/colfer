clean:
	go clean ./...
	rm -f cmd/colf/colf
	rm -f testdata/*.class testdata/bench/*.class

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
