include ../common.mk

DART ?= dart

.PHONY: test
test: gen/Colfer.dart gen_test.dart breaktest
	$(DART) analyze .
	$(DART) test .

gen/Colfer.dart: ../testdata/test.colf ../*.go ../cmd/colf/*.go
	$(COLF) Dart ../testdata/test.colf

breaktest: ../testdata/break*.colf ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ Dart ../testdata/break*.colf
	$(DART) analyze ./$@
	touch $@

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -fr gen
