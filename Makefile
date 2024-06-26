PACKAGE  = binenv
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
BIN      = $(CURDIR)/bin

GO      = go
GODOC   = godoc
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

export GO111MODULE=on

.PHONY: all


all: fmt lint $(BIN) ; $(info $(M) building executable…) @ ## Build program binary
	$Q $(GO) build \
		-tags release \
		-ldflags '-X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)

cache-all: cache validate distributions

build: linux

linux: fmt lint $(BIN) ; $(info $(M) building static executable for Linux……) @ ## Build program binary
	$Q env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-linux-amd64
	$Q env GOOS=linux GOARCH=386 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-linux-386

freebsd: fmt lint $(BIN) ; $(info $(M) building static executable for FreeBSD……) @ ## Build program binary
	$Q env GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-freebsd-amd64
	$Q env GOOS=freebsd GOARCH=386 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-freebsd-386

darwin: fmt lint $(BIN) ; $(info $(M) building static executable for MacOS…) @ ## Build program binary
	$Q env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-darwin-amd64

windows: fmt lint clean $(BIN) ; $(info $(M) building static executable for Windoze…) @ ## Build program binary
	$Q env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-win-amd64

release: windows darwin linux freebsd ; $(info $(M) stripping release executable for Linux…) @ ## Build program binary
	$Q strip $(BIN)/$(PACKAGE)-linux-amd64
	$Q strip $(BIN)/$(PACKAGE)-freebsd-amd64
	$Q (cd bin && sha256sum * > SHA256SUMS.txt)
	$Q cp $(BIN)/$(PACKAGE)-linux-amd64 $(BIN)/$(PACKAGE)
	# $Q gzip $(BIN)/$(PACKAGE)-linux-amd64
	# $Q gzip $(BIN)/$(PACKAGE)-linux-386
	# $Q gzip $(BIN)/$(PACKAGE)-darwin-amd64
	# $Q gzip $(BIN)/$(PACKAGE)-win-amd64
	$Q $(BIN)/$(PACKAGE) version

goreleaser-test: fmt lint clean ; $(info $(M) goreleaser dry-run…) @ ## Build program binary
	goreleaser release --snapshot --skip=publish --clean

goreleaser: fmt lint clean test; $(info $(M) create a release with goreleaser…) @ ## Build program binary
	goreleaser releaser --clean

prepush: outdated ; $(info $(M) execute CI linters…) @ ## execute linting tests so we should not fail liting in CI
	$Q $(GO) vet ./...
	$Q docker run  -v $(pwd)/README.md:/tmp/README.md pipelinecomponents/markdownlint:latest mdl --style all -r ~MD034,~MD013 /tmp/README.md

distributions: $(BIN) ; $(info $(M) creating DISTRIBUTIONS.md…) @ ## builds DISTRIBUTIONS.md file from distributions.yaml
	$Q ./bin/binenv search -w | sed "s,\x1B\[[0-9;]*[a-zA-Z],,g" | awk -F',' '{ print "- ["$$1"]("$$2"): "$$3","$$4","$$5","$$6","$$7","$$8}' | sed -e 's/,*$$//' | tr -d '"' > DISTRIBUTIONS.md 

validate: bin ; $(info $(M) validating cache against distributions…) @ ## validates cache against distributions
	$Q ./scripts/validate.sh code

cache: bin ; $(info $(M) building distribution cache…) @ ## builds distribution cache
	$Q ./scripts/buildcache.sh

e2e: bin ;  $(info $(M) runs end2end integration tests (very long)…) @ ## installs all supported distribution in a thowaway container
	$Q docker build . -t binenv-e2e
	$Q echo starting docker
	$Q docker run -ti --name binenv-e2e -e GITHUB_TOKEN --rm -v $(pwd)/distributions/distributions.yaml:/home/binenv/.config/binenv/distributions.yaml binenv-e2e

$(BIN):
	@mkdir -p $@

# Tools

GOLINT = $(GO) run golang.org/x/lint/golint@latest
GOCOVMERGE = $(GO) run github.com/wadey/gocovmerge@latest
GOCOV = $(GO) run github.com/axw/gocov/gocov@latest
GOCOVXML = $(GO) run github.com/AlekSi/gocov-xml@latest
GO2XUNIT = $(GO) run github.com/tebeka/go2xunit@latest
GOMODOUTDATED = $(GO) run github.com/psampaz/go-mod-outdated@latest

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test tests: fmt lint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint | $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q mkdir -p test
	$Q 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)/coverage
	$Q for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$Q $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: outdated
outdated: ; $(info $(M) running go-mod-outdated…) @ ## Run go-mod-outdated
	$Q $(GO) list -u -m -json all 2>/dev/null | $(GOMODOUTDATED) -update
	$Q $(GO) list -u -m -json all 2>/dev/null | $(GOMODOUTDATED) -update -direct

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt ./...

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
