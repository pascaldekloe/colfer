include ../common.mk

DART ?= dart

.PHONY: test
test: gen/Colfer.dart gen_test.dart breaktest
	$(DART) test .
	$(DART) analyze --fatal-infos .

gen/Colfer.dart: ../testdata/test.colf ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) Dart ../testdata/test.colf

breaktest: ../testdata/break*.colf ../*.go ../template/*.txt ../cmd/colf/*.go
	$(COLF) -b $@ Dart ../testdata/break*.colf
	touch $@

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -fr gen
