.PHONY: clean test

COLF?=../colf

CC=clang
CFLAGS=-std=c99 -g -O0 -pedantic -Werror


test: clean gen gen_test
	./gen_test

gen:
	$(COLF) -b gen C ../testdata/test.colf

Colfer.o:
	$(CC) -c $(CFLAGS) gen/Colfer.c

gen_test: Colfer.o
	$(CC) -o gen_test $(CFLAGS) Colfer.o gen_test.c

clean:
	rm -fr gen
	rm -fr *.o *_test *_test.dSYM
