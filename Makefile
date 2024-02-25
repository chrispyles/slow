.PHONY: build

build:
	go build -o build/slow

run: build
	build/slow

test:
	go test -v ./...
