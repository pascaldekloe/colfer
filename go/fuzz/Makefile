include ../../common.mk

.PHONY: test
test: fuzz.zip go-fuzz
	./go-fuzz -bin fuzz.zip -workdir ../../testdata

Colfer.go: ../../testdata/test.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) Go ../../testdata/test.colf
	mv gen/Colfer.go .
	rmdir gen

go-fuzz:
	$(GO) build -o $@ github.com/dvyukov/go-fuzz/go-fuzz

go-fuzz-build:
	$(GO) build -o $@ github.com/dvyukov/go-fuzz/go-fuzz-build

fuzz.zip: Colfer.go go-fuzz-build
	./go-fuzz-build -o $@ .

.PHONY: clean
clean:
	rm -f fuzz.zip Colfer.go

.PHONY: clean-all
clean-all: clean
	rm -f go-fuzz go-fuzz-build
