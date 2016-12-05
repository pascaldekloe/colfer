.PHONY: clean bench build test

build:
	go get github.com/pascaldekloe/colfer/cmd/...

clean:
	go clean -i -r ./cmd/...
	rm -fr dist

	make -C rpc clean
	make -C testdata clean
	make -C testdata/bench clean

test: build
	make -C testdata test
	make -C rpc test

bench:
	make -C testdata/bench bench
	make -C rpc bench

dist: clean test
	go fmt
	mkdir -p dist
	GOARCH=amd64 GOOS=linux go build -o dist/colf-linux ./cmd/colf
	GOARCH=amd64 GOOS=darwin go build -o dist/colf-darwin ./cmd/colf
	GOARCH=amd64 GOOS=windows go build -o dist/colf.exe ./cmd/colf
