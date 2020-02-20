VERSION_CLI := $(shell head ./VERSION_CLI)
VERSION_SERV := $(shell head ./VERSION_SERV)

NAME_CLI := lightf

DIR_CLI := ./cmd/lightf-cli/lightf-cli.go
DIR_SERV := ./cmd/lightf-serv/lightf-serv.go

.PHONY: server
server:
	go run $(DIR_SERV)

.PHONY: cli
cli:
	go build -o $(HOME)/go/bin/$(NAME_CLI) -v $(DIR_CLI)

.PHONY: install
install:
	go install $(DIR_CLI)
