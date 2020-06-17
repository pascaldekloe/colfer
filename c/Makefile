include ../common.mk

CC ?= clang
CFLAGS = -std=c11 -g -O0
CFLAGS_STRICT = $(CFLAGS) -pedantic -Werror

.PHONY: test
test: build
	build/gen_test

build: gen
	mkdir -p build
	$(CC) -o build/Colfer.o -c $(CFLAGS_STRICT) gen/Colfer.c
	$(CC) -o build/gen_test $(CFLAGS) build/Colfer.o gen_test.c

gen:
	$(COLF) -b gen C ../testdata/test.colf

.PHONY: clean
clean:
	rm -fr gen build
