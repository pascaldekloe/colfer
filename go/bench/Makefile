.PHONY: bench clean

COLF?=../../colf

bench: Colfer.go FlatBuffers.go scheme.pb.go
	go test -bench . -benchtime 3s

clean:
	go clean
	rm -f Colfer.go FlatBuffers.go scheme.pb.go

Colfer.go:
	$(COLF) -b .. -f go ../../testdata/bench/scheme.colf

scheme.pb.go:
	protoc --gogofaster_out=. -I../../testdata/bench -I./vendor -I./vendor/github.com/gogo/protobuf/protobuf ../../testdata/bench/scheme.proto

FlatBuffers.go:
	flatc -o .. -g ../../testdata/bench/scheme.fbs
