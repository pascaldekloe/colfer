.PHONY: clean test

CXX?=clang
CXXFLAGS=-std=c++11 -g -O0

test: clean gen_test

Colfer.o:
	$(CXX) -c $(CXXFLAGS) -pedantic -Werror gen/Colfer.cpp

gen_test: Colfer.o
	$(CXX) -o gen_test $(CXXFLAGS) Colfer.o gen_test.cpp
	./gen_test

clean:
	rm -fr *.o *_test *_test.dSYM
