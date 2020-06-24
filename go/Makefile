include ../common.mk

.PHONY: test
test: gen build
	$(GO) test -v -coverprofile build/coverage -coverpkg github.com/pascaldekloe/colfer/go/gen
	$(GO) build ./build/break/...

gen:
	$(COLF) -t ../testdata/test-go.tags Go ../testdata/test.colf

build:
	mkdir -p build
	$(COLF) -p github.com/pascaldekloe/colfer/go/build/break go ../testdata/break*.colf

fuzz.zip: gen
	$(GO) get github.com/dvyukov/go-fuzz/go-fuzz-build
	$(GOPATH)/bin/go-fuzz-build -o fuzz.zip github.com/pascaldekloe/colfer/go

.PHONY: fuzz
fuzz: fuzz.zip
	rm ../testdata/corpus/seed*
	$(GO) test -run FuzzSeed

	$(GO) get github.com/dvyukov/go-fuzz/go-fuzz
	$(GOPATH)/bin/go-fuzz -bin fuzz.zip -workdir ../testdata

.PHONY: clean
clean:
	rm -fr gen build fuzz.zip
