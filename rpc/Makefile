.PHONY: clean test bench

test: internal
	go vet
	go fmt
	go test -v

bench: internal
	go test -bench .

clean:
	rm -f testcover
	rm -fr internal

internal:
	go generate -x

testcover:
	go test -v -coverprofile testcover
