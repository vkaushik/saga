help: ## Display all make commands with comments.
	@grep -h -E '^[a-zA-Z0-9_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"go test -v -race $(go list ./... | grep -v vendor)"

format:
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"./build/format.sh"

clean: ## Clean the docker residues.
	docker-compose -f ./build/docker-compose.yml down -v --rmi all --remove-orphans

test_mocks: ## Generate mocks for tests.
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"./build/mocks/generate_mocks.sh"

build_docker_images: ## Rebuild base image after any changes to Dockerfile.
	docker-compose -f ./build/docker-compose.yml build

get_running_container:
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"tail -f /dev/null"

##---------------------------------------------
## Git Hooks

pre_commit: ## Run it before committing changes.
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"go mod tidy && go mod vendor && go vet ./... && go fmt ./... && errcheck ./..."

pre_push: test ## Run it before pushing changes.

install_git_hooks: ## Install pre-commit and pre-push git hooks
	if [ -f ./.git/hooks/pre-commit ]; then mv ./.git/hooks/pre-commit ./git/hooks/old-pre-commit; fi
	if [ -f ./.git/hooks/pre-push]; then mv ./.git/hooks/pre-push ./git/hooks/old-pre-push; fi
	ln -s ./build/git-hooks/pre-commit ./.git/hooks/pre-commit
	ln -s ./build/git-hooks/pre-push ./.git/hooks/pre-push
	chmod +x ./.git/hooks/pre-commit ./.git/hooks/pre-push

uninstall_git_hooks: ## Uninstall pre-commit and pre-push git hooks
	rm ./.git/hooks/pre-commit
	rm ./.git/hooks/pre-push
	if [ -f ./.git/hooks/old-pre-commit ]; then mv ./.git/hooks/old-pre-commit ./git/hooks/old-pre-commit; fi
	if [ -f ./.git/hooks/old-pre-push ]; then mv ./.git/hooks/old-pre-push ./git/hooks/old-pre-push; fi

##---------------------------------------------
## Others

dependency_graph:
	docker-compose -f ./build/docker-compose.yml run gontainer bash -c \
	"godepgraph github.com/vkaushik/saga | dot -Tpng -o dependency_graph.png"