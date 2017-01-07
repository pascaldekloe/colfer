.PHONY: clean test

COLF?=../colf

CC?=clang
CFLAGS=-std=c99 -g -O0 -pedantic -Werror


test: clean gen gen_test
	./gen_test

gen:
	$(COLF) -b gen C ../testdata/test.colf

Colfer.o:
	$(CC) $(CFLAGS) -c  gen/Colfer.c

gen_test: Colfer.o
	$(CC) $(CFLAGS) -o gen_test Colfer.o gen_test.c

clean:
	rm -fr gen
	rm -fr *.o *_test *_test.dSYM
