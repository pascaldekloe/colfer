include ../common.mk

CFLAGS += -O0 -g -pedantic -Werror -std=c2x

test: test.c test.h Colfer.o
	$(CC) $(CFLAGS) -o $@ test.c Colfer.o
	./test

breaktest: ../testdata/break*.colf ../template/*.txt ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ C ../testdata/break*.colf
	$(CC) $(CFLAGS) -c $@/Colfer.c
	touch $@

Colfer.c: ../testdata/seal.colf ../template/*.txt ../*.go ../cmd/colf/*.go
	$(COLF) C ../testdata/seal.colf

Colfer.h: Colfer.c

Colfer.o: Colfer.c
	$(CC) $(CFLAGS) -o $@ -c $?

.PHONY: clean
clean:
	rm -f test
	rm -f *.o
	rm -f Colfer.[ch]
	rm -fr breaktest
