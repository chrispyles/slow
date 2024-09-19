.PHONY: build

build:
	go build -o build/slow

run: build
	build/slow

test:
	go test -v ./...

testcov:
	go test -v -race -covermode atomic -coverprofile=profile.cov ./...
	go tool covdata textfmt -i .coverdata -o coverage2.out
	cat coverage2.out | tail -n +2 >> profile.cov
	rm -r coverage2.out .coverdata
	go tool cover -html=profile.cov -o=coverage.html
