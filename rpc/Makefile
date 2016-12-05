.PHONY: clean test bench

test: internal
	go vet
	go fmt
	go test

bench: internal
	go test -bench .

clean:
	rm -f testcover
	rm -fr internal

internal:
	go generate

testcover:
	go test -coverprofile testcover
