BUILD_DATE          := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT			:=$(shell git rev-parse "HEAD^{commit}" 2>/dev/null)
GIT_VERSION         := $(shell git describe --tags --match='v*' --abbrev=14 "$(GIT_COMMIT)^{commit}" 2>/dev/null)

.PHONY: all
all: style ut format vet install

.PHONY: style
style:
	@echo ">> checking code style"
	@! gofmt -d $$(find ./ -path ./vendor -prune -o -name '*.go' -print) | grep '^'

.PHONY: ut
ut:
	@echo ">> running all tests"
	@go test -cover -count=1 -gcflags=-l -race

.PHONY: format
format:
	@echo ">> formatting code"
	@go fmt

.PHONY: vet
vet:
	@echo ">> vetting code"
	@go vet -all

.PHONY: install
install:
	@echo ">> install zkcmd"
	@go install -ldflags "\
	-X 'github.com/benzimu/zkcmd/common/version.buildDate=$(BUILD_DATE)' \
	-X 'github.com/benzimu/zkcmd/common/version.version=$(GIT_VERSION)' \
	-X 'github.com/benzimu/zkcmd/common/version.gitCommit=$(GIT_COMMIT)'\
	" zkcmd.go

.PHONY: uninstall
uninstall:
	@rm -rf $(GOPATH)/bin/zkcmd
