.PHONY: bench clean test

COLF?=../colf

test: internal
	go vet
	go fmt

	$(COLF) -b build go ../testdata/test.colf

	go test -v

testcover:
	go test -v -coverprofile testcover

bench: test
	go test -bench .

clean:
	rm -fr internal testcover

internal:
	$(COLF) go
