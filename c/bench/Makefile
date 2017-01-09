.PHONY: clean run

COLF?=../../colf

run: clean bench
	./bench 20000000

bench: gen Colfer.o
	g++ -o bench -std=c++11 -O3 bench.cpp Colfer.o

gen:
	$(COLF) -b gen -f c ../../testdata/bench/scheme.colf
	protoc --cpp_out=gen -I../../testdata/bench ../../testdata/bench/scheme.proto
	flatc -c -o gen ../../testdata/bench/scheme.fbs

Colfer.o: gen
	gcc -c -O3 gen/Colfer.c

clean:
	rm -f *.o bench
	rm -fr gen
