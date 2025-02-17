.PHONY: build

COVERAGEARGS = -race -cover -covermode atomic
BUILDCOVOUT  = 
COVERDIR     := $(shell pwd)/.coverdata

build:
	go build -o build/slow

run: build
	build/slow

test: export SLOW_TESTING_GOCOVERDIR := $(COVERDIR)
test:
	@mkdir -p $(COVERDIR)
	@go test -v ./...

# The buildcov rule requires the user to set BUILDCOVOUT to the path that the output binary should
# be written to; for integration tests, this is handled in integration_test.go.
buildcov:
	@go build $(COVERAGEARGS) -o $(BUILDCOVOUT) .

testcov: export SLOW_TESTING_GOCOVERDIR := $(COVERDIR)
testcov:
	@mkdir -p $(COVERDIR)
	@go test -shuffle=on $(COVERAGEARGS) ./... -args -test.gocoverdir="$(COVERDIR)"
	@echo "=== Coverage Summary ==="
	@go tool covdata percent -i $(COVERDIR)
	@echo "=== Combining coverage data and saving to profile.cov ==="
	@go tool covdata textfmt -i $(COVERDIR) -o profile.cov
	@rm -r $(COVERDIR)
	@go tool cover -html=profile.cov -o=coverage.html
