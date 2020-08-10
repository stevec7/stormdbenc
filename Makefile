VERSION=$(shell git describe --always --long --dirty)
.PHONY: all clean

all: test build

build:
	mkdir -p dist
	go build -v -o dist/stormdbenctest cmd/stormdbenctest/main.go

clean:
	rm dist/stormdbenctest

test:
	cd pkg/stormdbenc && go test  || (echo "Tests failed"; exit 1)

