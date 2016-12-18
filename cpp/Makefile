.PHONY: clean test

CXX?=clang
CXXFLAGS?=-std=c++11 -pedantic -Werror

test: clean a.out
	./a.out

a.out:
	$(CXX) $(CXXFLAGS) test.cpp gen/Colfer.cpp

clean:
	rm -f a.out
