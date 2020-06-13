include ../../common.mk

FLATC ?= flatc
PROTOC ?= protoc

.PHONY: test
test: build
	$(GO) test -bench .

stat.txt: build
	$(GO) test -bench . -count 5 > build/bench.log
	benchstat build/bench.log > stat.txt
	@echo >> stat.txt
	@$(GO) version >> stat.txt
	@$(PROTOC) --version >> stat.txt
	@$(FLATC) --version >> stat.txt

build:
	$(COLF) -b build/gen Go ../../testdata/bench/scheme.colf
	$(PROTOC) --gogofaster_out=build/gen/bench -I../../testdata/bench ../../testdata/bench/scheme.proto
	$(FLATC) -o build/gen -g ../../testdata/bench/scheme.fbs

.PHONY: clean
clean:
	rm -fr build stat.txt
