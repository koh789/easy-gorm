.PHONY: help
help: ## help you can see the list of tasks with `make help`.
	@echo "------- タスク一覧 ------"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36mmake %-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install-all
install-all:
	make install-goimports
	make install-golangci-lint

.PHONY: install-goimports ## install goimports
install-goimports:
	which goimports || go get golang.org/x/tools/cmd/goimports

.PHONY: install-golangci-lint
install-golangci-lint: ## install golangci-lint
	which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.52.1

.PHONY: test-all
test-all: ## static analysis & testing
	make lint
	make test

.PHONY: lint
lint: ## golangci-lint
	make install-golangci-lint
	golangci-lint cache clean
	golangci-lint run

.PHONY: test
test: ## test
	@echo "+ go test..."
	go clean ./... && go test -v ./...
	@echo "+ go test clear."
