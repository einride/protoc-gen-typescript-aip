SHELL := /usr/bin/env bash

all: \
	go-lint \
	go-review \
	example/proto \
	eslint \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/commitlint/rules.mk
include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/semantic-release/rules.mk
include tools/eslint/rules.mk

.PHONY: example/proto
example/proto:
	@make -C $@

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@go test -count 1 -cover -race ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy -v

.PHONY: $(eslint)
eslint: $(eslint)
	$(info [$@] linting typescript files...)
	$(eslint) --config $(eslint_cwd)/.eslintrc.js --quiet "example/proto/gen/typescript/**/*.ts"
