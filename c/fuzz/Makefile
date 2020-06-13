CC ?= afl-clang
CFLAGS = -std=c11 -g -O0

FUZZ ?= afl-fuzz

.PHONY: test
test: build
	mkdir -p build/seed
	cp ../../testdata/corpus/seed* build/seed
	$(FUZZ) -i build/seed -o findings build/fuzz

build:
	mkdir build
	$(CC) -o build/fuzz-Colfer.o -c $(CFLAGS) ../gen/Colfer.c
	$(CC) -o build/fuzz $(CLFAGS) build/fuzz-Colfer.o fuzz.c

.PHONY: clean
clean:
	rm -fr build
