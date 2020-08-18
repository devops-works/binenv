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

build: linux

linux: fmt lint $(BIN) ; $(info $(M) building static executable for Linux……) @ ## Build program binary
	$Q env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-amd64-$(VERSION)

darwin: fmt lint $(BIN) ; $(info $(M) building static executable for MacOS…) @ ## Build program binary
	$Q env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-darwin-$(VERSION)

windows: fmt lint clean $(BIN) ; $(info $(M) building static executable for Windoze…) @ ## Build program binary
	$Q env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-tags release -a \
		-ldflags '-w -extldflags "-static" -X github.com/devops-works/binenv/cmd.Version=$(VERSION) -X github.com/devops-works/binenv/cmd.BuildDate=$(DATE)' \
		-o $(BIN)/$(PACKAGE)-win-$(VERSION)

release: windows darwin linux ; $(info $(M) stripping release executable for Linux…) @ ## Build program binary
	$Q strip $(BIN)/$(PACKAGE)-amd64-$(VERSION)
	$Q (cd bin && sha256sum * > SHA256SUMS)
	$Q cp $(BIN)/$(PACKAGE)-amd64-$(VERSION) $(BIN)/$(PACKAGE)
	$Q gzip $(BIN)/$(PACKAGE)-amd64-$(VERSION)
	$Q gzip $(BIN)/$(PACKAGE)-darwin-$(VERSION)
	$Q gzip $(BIN)/$(PACKAGE)-win-$(VERSION)
	$Q $(BIN)/$(PACKAGE) version


# don't even think about making a joke with this target name
prepush: outdated ; $(info $(M) execute CI linters…) @ ## execute linting tests so we should not fail liting in CI
	$Q $(GO) vet ./...
	$Q docker run  -v $(pwd)/README.md:/tmp/README.md pipelinecomponents/markdownlint:latest mdl --style all -r ~MD034,~MD013 /tmp/README.md

# Tools

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(REPOSITORY)…)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(REPOSITORY) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: REPOSITORY=golang.org/x/lint/golint

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: REPOSITORY=github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: REPOSITORY=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: REPOSITORY=github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: REPOSITORY=github.com/tebeka/go2xunit

GOMODOUTDATED = $(BIN)/go-mod-outdated
$(BIN)/go-mod-outdated: REPOSITORY=github.com/psampaz/go-mod-outdated

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt lint ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint | $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests with xUnit output
	$Q mkdir -p test
	$Q 2>&1 $(GO) test -timeout 20s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
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
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: outdated
outdated: | $(GOMODOUTDATED) ; $(info $(M) running go-mod-outdated…) @ ## Run go-mod-outdated
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
