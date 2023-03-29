include ../../common.mk

target: src/main/resources/amd64/colf-linux src/main/resources/amd64/colf-darwin src/main/resources/amd64/colf-openbsd src/main/resources/amd64/colf.exe
	mvn package

bundle.zip: target
	mvn source:jar javadoc:jar
	cp target/colfer-maven-plugin* .
	cp pom.xml `ls colfer-maven-plugin-*-javadoc.jar | sed 's/-javadoc.jar/.pom/'`
	for f in `ls colfer-maven-plugin*.{jar,pom}`; do gpg -ab $$f; done
	zip bundle colfer-maven-plugin*

src/main/resources/amd64/colf-linux: ../../*.go ../../cmd/colf/*.go
	GOARCH=amd64 GOOS=linux $(GO) build -o src/main/resources/amd64/colf-linux ../../cmd/colf
	touch $@

src/main/resources/amd64/colf-darwin: ../../*.go ../../cmd/colf/*.go
	GOARCH=amd64 GOOS=darwin $(GO) build -o src/main/resources/amd64/colf-darwin ../../cmd/colf
	touch $@

src/main/resources/amd64/colf-openbsd: ../../*.go ../../cmd/colf/*.go
	GOARCH=amd64 GOOS=openbsd $(GO) build -o src/main/resources/amd64/colf-openbsd ../../cmd/colf
	touch $@

src/main/resources/amd64/colf.exe: ../../*.go ../../cmd/colf/*.go
	GOARCH=amd64 GOOS=windows $(GO) build -o src/main/resources/amd64/colf.exe ../../cmd/colf
	touch $@

.PHONY: clean
clean:
	rm -fr target
	rm -f colfer-maven-plugin* bundle.zip

.PHONY: clean-all
clean-all: clean
	rm -fr src/main/resources/amd64/*

.PHONY: reset
reset:
	git checkout src/main/resources/amd64
