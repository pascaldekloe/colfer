include ../common.mk

CC=clang
CFLAGS=-std=c99 -g -O0 -pedantic -Werror

.PHONY: test
test: build
	build/gen_test

build: gen
	mkdir -p build
	$(CC) -o build/Colfer.o -c $(CFLAGS) gen/Colfer.c
	$(CC) -o build/gen_test $(CFLAGS) build/Colfer.o gen_test.c

gen: install
	$(COLF) -b gen C ../testdata/test.colf

.PHONY: clean
clean:
	rm -fr gen build
