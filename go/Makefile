.PHONY: clean test testcover fuzz

COLF?=../colf

test: gen
	go vet
	go fmt
	go test -v

	$(COLF) -b ../../../.. -p github.com/pascaldekloe/colfer/go/build/break go ../testdata/break*.colf
	go build ./build/break/...

testcover: gen
	go test -v -coverprofile testcover -coverpkg github.com/pascaldekloe/colfer/go/gen

clean:
	go clean ./...

	rm -fr build fuzz.zip gen testcover

gen:
	$(COLF) go ../testdata/test.colf

fuzz.zip: gen
	go get github.com/dvyukov/go-fuzz/go-fuzz-build
	go-fuzz-build -o fuzz.zip github.com/pascaldekloe/colfer/go

fuzz: fuzz.zip
	rm ../testdata/corpus/seed*
	go test -run FuzzSeed

	go get github.com/dvyukov/go-fuzz/go-fuzz
	go-fuzz -bin fuzz.zip -workdir ../testdata
