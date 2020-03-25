include ../../common.mk

FLATC?=flatc
PROTOC?=protoc

.PHONY: test
test: build
	$(GO) test -bench .

.PHONY: stat
stat: build
	$(GO) test -bench . -count 5 > build/bench.log
	benchstat build/bench.log
	@echo
	@$(GO) version
	@$(PROTOC) --version
	@$(FLATC) --version

build: install
	$(COLF) -b build/gen Go ../../testdata/bench/scheme.colf
	$(PROTOC) --gogofaster_out=build/gen/bench -I../../testdata/bench ../../testdata/bench/scheme.proto
	$(FLATC) -o build/gen -g ../../testdata/bench/scheme.fbs

.PHONY: clean
clean:
	$(GO) clean
	rm -fr build
