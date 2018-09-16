.PHONY: help
.DEFAULT_GOAL := help

deps: ## install dependency with deps
	cd src/ && dep ensure

build: ## build interpreter
	go build -o interpreter ./src/cmd/main.go

run: ## run interpreter
	go run ./src/cmd/main.go

test: ## test with gotest
	gotest -v ./...

watch: ## running tests watching files with looper
	looper

docker-build: ## create docker image
	docker build -t interpreter .

docker-run: ## run interpreter with docker
	docker run -e LOG_LEVEL=debug -v $(PWD):/go/src/github.com/g-hyoga/writing-interpreter-in-go interpreter go run ./src/cmd/main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
