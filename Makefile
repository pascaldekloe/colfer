build:
	go generate

test: build
	go vet ./...
	go test
