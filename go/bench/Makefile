include ../../common.mk

FLATC ?= flatc
PROTOC ?= protoc

BENCHN ?= 10

bench.txt: Colfer.go scheme.pb.go FlatBuffers.go
	$(GO) test -bench . > $@
	@echo >> $@
	@$(GO) version >> $@
	@$(PROTOC) --version >> $@
	@$(FLATC) --version >> $@

new.txt: Colfer.go scheme.pb.go FlatBuffers.go
	$(GO) test -bench /colfer -count $(BENCHN) > $@

old.txt:
	@[ -f new.txt ] || echo '# Run "make new.txt" to produce a reference first.'
	cp new.txt old.txt

.PHONY: diff
diff: old.txt new.txt
	benchstat old.txt new.txt

Colfer.go: ../../testdata/bench/scheme.colf ../../*.go ../../cmd/colf/*.go
	$(COLF) -b .. -s 1024 Go ../../testdata/bench/scheme.colf

scheme.pb.go: ../../testdata/bench/scheme.proto
	$(PROTOC) --gogofaster_out=. -I../../testdata/bench ../../testdata/bench/scheme.proto

FlatBuffers.go: ../../testdata/bench/scheme.fbs
	$(FLATC) -o .. -g ../../testdata/bench/scheme.fbs

.PHONY: clean
clean:
	rm -fr bench.txt new.txt old.txt
	rm -fr Colfer.go scheme.pb.go FlatBuffers.go
