.PHONY: clean test

COLF?=../colf

CC?=clang


test: clean gen gen_test
	./gen_test

gen:
	$(COLF) -b gen C ../testdata/test.colf

Colfer.o:
	$(CC) -c -std=c99 -g -O0 -pedantic -Werror gen/Colfer.c

gen_test: Colfer.o
	$(CC) -o gen_test -g Colfer.o gen_test.c

clean:
	rm -fr gen
	rm -fr *.o *_test *_test.dSYM
