include ../../common.mk

test: build
	go test -bench . -benchtime 3s

build: install
	$(COLF) -b build/gen Go ../../testdata/bench/scheme.colf
	protoc --gogofaster_out=build/gen/bench -I../../testdata/bench -I./vendor -I./vendor/github.com/gogo/protobuf/protobuf ../../testdata/bench/scheme.proto
	flatc -o build/gen -g ../../testdata/bench/scheme.fbs

.PHONY: clean
clean:
	go clean
	rm -fr build
