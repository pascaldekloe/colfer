include ../../common.mk

FLATC ?= flatc
PROTOC ?= protoc

CFLAGS += -O3
CXXFLAGS += -O3

bench.txt: bench
	./bench > $@
	@echo >> $@
	@$(FLATC) --version >> $@
	@$(PROTOC) --version >> $@
	@$(CC) --version >> $@

Colfer.h Colfer.c &: ../../testdata/bench/scheme.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) -f -s 1024 C ../../testdata/bench/scheme.colf

# ProtoBuf
scheme.pb.cc: ../../testdata/bench/scheme.proto
	$(PROTOC) --cpp_out=. -I../../testdata/bench ../../testdata/bench/scheme.proto

# FlatBuffers
scheme_generated.h: ../../testdata/bench/scheme.fbs
	$(FLATC) --cpp -o . ../../testdata/bench/scheme.fbs

Colfer.o: Colfer.h Colfer.c
	$(CC) $(CFLAGS) -o $@ -c -std=c14 Colfer.c

scheme.pb.o: scheme.pb.cc
	$(CXX) $(CXXFLAGS) -o $@ -c -std=c++14 scheme.pb.cc

bench: bench.cpp Colfer.o scheme.pb.o scheme_generated.h
	$(CXX) $(CXXFLAGS) -o $@ -lstdc++ -lprotobuf -lbenchmark -std=c++14 bench.cpp Colfer.o scheme.pb.o

.PHONY: clean
clean:
	rm -f bench bench.txt
	rm -f Colfer.* scheme.pb.* scheme_generated.*
