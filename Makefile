help: ## Display all make commands with comments.
	@grep -h -E '^[a-zA-Z0-9_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: install_errcheck ## Run it in local after taking the clone.

test: ## Run tests
	go test -v -race $(go list ./... | grep -v vendor)

format:
	./build/format.sh

pre_commit: ## Run it before committing changes.
	go mod tidy
	go mod vendor
	go vet ./...
	go fmt ./...
	errcheck ./...

install_errcheck:
	go get github.com/kisielk/errcheck

##---------------------------------------------
