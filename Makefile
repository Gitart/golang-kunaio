.PHONY: all test fmt clean

export GOPATH=$(shell pwd)

all:
	$(MAKE) -C cli $@

test:
	go test -v kunaio
	$(MAKE) -C cli $@

fmt:
	find . -type f -name \*.go -exec gofmt -w '{}' ';'

clean:
	$(MAKE) -C cli $@
