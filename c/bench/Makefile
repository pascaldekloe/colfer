include ../../common.mk

.PHONY: test
test: build
	build/bench 10000000

build: install
	$(COLF) -b build/gen -f -s 1024 c ../../testdata/bench/scheme.colf
	protoc --cpp_out=build/gen -I../../testdata/bench ../../testdata/bench/scheme.proto
	flatc --cpp -o build/gen ../../testdata/bench/scheme.fbs

	$(CC) $(CFLAGS) -o build/Colfer.o -c -std=c99 -O3 build/gen/Colfer.c
	$(CC) $(CFLAGS) -o build/bench -std=c++11 -O3 -lstdc++ bench.cpp build/Colfer.o

.PHONY: clean
clean:
	rm -fr build
