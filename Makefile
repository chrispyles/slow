.PHONY: build

build:
	go build -o build/slow

run: build
	build/slow

test:
	go test -v ./...

testcov:
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o=coverage.html
