include ../common.mk

CFLAGS += -O0 -pedantic -Werror

.PHONY: test
test: gen_test breaktest
	./gen_test

breaktest: ../testdata/break*.colf ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ C ../testdata/break*.colf
	$(CC) $(CFLAGS) -c -std=c11 $@/Colfer.c
	touch $@

Colfer.h Colfer.c &: ../testdata/test.colf ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) C ../testdata/test.colf

Colfer.o: Colfer.h Colfer.c
	$(CC) $(CFLAGS) -o $@ -c -std=c11 Colfer.c

gen_test: gen_test.h gen_test.c Colfer.o
	$(CC) $(CFLAGS) -o $@ -std=c11 gen_test.c Colfer.o

.PHONY: clean
clean:
	rm -fr gen_test *.o breaktest

.PHONY: clean-all
clean-all: clean
	rm -f Colfer.?
