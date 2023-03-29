include ../../common.mk

CC = afl-clang
CFLAGS ?= -O2

FUZZ ?= afl-fuzz

.PHONY: test
test: fuzz
	mkdir -p seed
	cp ../../testdata/corpus/seed* seed
	$(FUZZ) -i seed -o findings ./fuzz

Colfer.h Colfer.c &: ../../testdata/test.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) -f -s 1024 C ../../testdata/test.colf

Colfer.o: Colfer.h Colfer.c
	$(CC) $(CFLAGS) -o $@ -c -std=c11 Colfer.c

fuzz: fuzz.c Colfer.o
	$(CC) $(CFLAGS) -o $@ -std=c11 fuzz.c Colfer.o

.PHONY: clean
clean:
	rm -f fuzz *.o
	rm -f Colfer.?

clean-all: clean
	rm -fr seed findings
