help: ## Display all make commands with comments.
	@grep -h -E '^[a-zA-Z0-9_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"go test -v -race $(go list ./... | grep -v vendor)"

format:
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"./build/format.sh"

pre_commit: ## Run it before committing changes.
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"go mod tidy && go mod vendor && go vet ./... && go fmt ./... && errcheck ./..."

clean: ## Clean the docker residues.
	docker-compose -f ./build/docker-compose.yml down -v --rmi all --remove-orphans

test_mocks: ## Generate mocks for tests.
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"./build/mocks/generate_mocks.sh"

rebuild_docker_image:
	docker-compose -f ./build/docker-compose.yml build

##---------------------------------------------
