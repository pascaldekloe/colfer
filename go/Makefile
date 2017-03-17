include ../common.mk

.PHONY: test
test: gen build
	go test -v -coverprofile build/coverage -coverpkg github.com/pascaldekloe/colfer/go/gen
	go build ./build/break/...

gen: install
	$(COLF) Go ../testdata/test.colf

build: install
	mkdir -p build
	$(COLF) -b ../../../.. -p github.com/pascaldekloe/colfer/go/build/break go ../testdata/break*.colf

fuzz.zip: gen
	go get github.com/dvyukov/go-fuzz/go-fuzz-build
	$(GOPATH)/bin/go-fuzz-build -o fuzz.zip github.com/pascaldekloe/colfer/go

.PHONY: fuzz
fuzz: fuzz.zip
	rm ../testdata/corpus/seed*
	go test -run FuzzSeed

	go get github.com/dvyukov/go-fuzz/go-fuzz
	$(GOPATH)/bin/go-fuzz -bin fuzz.zip -workdir ../testdata

.PHONY: clean
clean:
	go clean .
	rm -fr gen build fuzz.zip
