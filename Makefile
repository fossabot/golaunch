GO=go

all: test build

build:
	$(GO) build

test:
	$(GO)	test -v ./...
