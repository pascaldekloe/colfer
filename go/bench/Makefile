include ../../common.mk

FLATC?=flatc
PROTOC?=protoc

.PHONY: test
test: build
	go test -bench .

.PHONY: stat
stat: build
	go test -bench . -count 5 > build/bench.log
	benchstat build/bench.log
	@echo
	@go version
	@$(PROTOC) --version
	@$(FLATC) --version

build: install
	$(COLF) -b build/gen Go ../../testdata/bench/scheme.colf
	$(PROTOC) --gogofaster_out=build/gen/bench -I../../testdata/bench -I./vendor -I./vendor/github.com/gogo/protobuf/protobuf ../../testdata/bench/scheme.proto
	$(FLATC) -o build/gen -g ../../testdata/bench/scheme.fbs

.PHONY: clean
clean:
	go clean
	rm -fr build
