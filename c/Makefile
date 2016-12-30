.PHONY: clean

CC?=clang
CFLAGS=-std=c99 -g -O0 -pedantic -Werror

pass: clean gen_test
	./gen_test && touch pass

Colfer.o:
	$(CC) $(CFLAGS) -c  gen/Colfer.c

gen_test: Colfer.o
	$(CC) $(CFLAGS) -o gen_test Colfer.o gen_test.c

clean:
	rm -fr *.o *_test *_test.dSYM pass
