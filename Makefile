.PHONY: clean deploy

deploy: colf
	make -C c clean test
	make -C ecma clean test
	make -C go clean test
	make -C java clean test
	make -C rpc clean test

colf:
	go build ./cmd/colf

install:
	go get ./cmd/...

clean:
	go clean -i ./cmd/...
	rm -fr colf dist */build

dist: clean deploy
	go fmt ./...
	go vet ./...

	mkdir -p dist
	GOARCH=amd64 GOOS=linux go build -o dist/colf-linux ./cmd/colf
	GOARCH=amd64 GOOS=darwin go build -o dist/colf-darwin ./cmd/colf
	GOARCH=amd64 GOOS=windows go build -o dist/colf.exe ./cmd/colf
