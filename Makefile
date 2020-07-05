PROJECTNAME=$(shell basename "$(PWD)")
MAKEFLAGS += --silent

.PHONY: all
all: help

.PHONY: install
## install: Install deps
install:
	@echo "	Installing deps..."
	@go mod download

.PHONY: gen
## gen: Generate wire injectors
gen:
	@echo "	Generating..."
	@find . -type f -name 'wire_gen.go' -delete
	@cd internal/store && wire && cd ../..
	@cd internal/app && wire && cd ../..

.PHONY: dev
## dev: Run bailiff in dev mode (live reload on change)
dev:
	@echo " Running dev mode..."
	@watchexec -e go -r "clear && make start >/dev/null"

.PHONY: start
## start: Run bailiff in normal mode
start:
	@echo "	Running "$(PROJECTNAME)"..."
	@go run cmd/bailiff/main.go

.PHONY: test
## test: Run unit tests
test:
	@echo "	Running tests..."
	@go test ./...

.PHONY: test-watch
## test-watch: Run unit tests in watch mode (rerun on change)
test-watch:
	@echo " Running tests in watch mode..."
	go test ./...
	watchexec -e go -r "go test ./..."

.PHONY: example
## example: Run example in docker-compose
example:
	@echo " Running example..."
	docker-compose -f examples/compose/docker-compose.yml up --build

.PHONY: help
help: Makefile
	@echo
	@echo "	Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
