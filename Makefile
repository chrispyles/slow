.PHONY: build

COVERAGEARGS     = -race -cover -covermode atomic
BUILDCOVOUT      = 
COVERDIR         := $(shell pwd)/.coverdata
IGNORECOVPATTERN = /\/testing\//

build:
	go build -o build/slow

run: build
	build/slow


test:
	@rm -rf $(COVERDIR)
	@mkdir $(COVERDIR)
	@go test -shuffle=off -v ./...
	@rm -r $(COVERDIR)

# The build_integration_test and buildcov rules require the user to set BUILDCOVOUT to the path that
# the output binary should be written to; for integration tests, this is handled in
# integration_test.go.

build_integration_test:
	@go build -o $(BUILDCOVOUT) .

buildcov:
	@go build $(COVERAGEARGS) -o $(BUILDCOVOUT) .

testcov: export SLOW_TESTING_GOCOVERDIR := $(COVERDIR)
testcov:
	@rm -rf $(COVERDIR) profile.cov coverage.html
	@mkdir $(COVERDIR)
	@go test -shuffle=on $(COVERAGEARGS) ./... -args -test.gocoverdir="$(COVERDIR)"
	@echo "=== Coverage Summary ==="
	@go tool covdata percent -i $(COVERDIR)
	@echo "=== Combining coverage data and saving to profile.cov ==="
	@go tool covdata textfmt -i $(COVERDIR) -o profile.cov
	@rm -r $(COVERDIR)
	@echo "=== Filtering out ignored files from coverage data ==="
	@if [[ $$(uname -s) == "Darwin" ]]; then sed -i '' '$(IGNORECOVPATTERN)d' profile.cov; else sed -i '$(IGNORECOVPATTERN)d' profile.cov; fi
	@echo "=== Generating coverage HTML report ==="
	@go tool cover -html=profile.cov -o=coverage.html
	@echo "=== Done ==="
