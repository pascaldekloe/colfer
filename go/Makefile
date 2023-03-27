include ../common.mk

.PHONY: test
test: Colfer.go gen_test.go breaktest
	$(GO) test -v

Colfer.go: ../testdata/test.colf ../testdata/test-go.tags ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -t ../testdata/test-go.tags Go ../testdata/test.colf
	mv gen/Colfer.go .
	rmdir gen

breaktest: ../testdata/break*.colf ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -p github.com/pascaldekloe/colfer/go/$@ go ../testdata/break*.colf
	$(GO) build ./$@/...
	touch $@

.PHONY: fuzz
fuzz: test
	$(GO) test -fuzz O

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -f Colfer.go
