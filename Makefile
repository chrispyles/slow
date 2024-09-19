.PHONY: build

build:
	go build -o build/slow

run: build
	build/slow

test:
	go test -v ./...

testcov:
	go test -v -coverprofile=coverage.out ./...
	go tool covdata textfmt -i .coverdata -o coverage2.out
	cat coverage2.out | tail -n +2 >> coverage.out
	rm -r coverage2.out .coverdata
	go tool cover -html=coverage.out -o=coverage.html
