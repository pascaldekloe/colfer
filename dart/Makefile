include ../common.mk

DART ?= dart

.PHONY: test
test: colfer_gen.dart colfer_gen_test.dart breaktest
	$(DART) test .

colfer_gen.dart: ../testdata/test.colf ../*.go ../cmd/colf/*.go
	$(COLF) Dart ../testdata/test.colf
	$(DART) analyze .

breaktest: ../testdata/break*.colf ../*.go ../cmd/colf/*.go
	$(COLF) -b $@ Dart ../testdata/break*.colf
	$(DART) analyze ./$@
	touch $@

.PHONY: clean
clean:
	rm -fr breaktest

.PHONY: clean-all
clean-all: clean
	rm -f colfer_gen.dart
