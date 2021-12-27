SHELL := /usr/bin/env bash

.PHONY: all
all: \
	clang-format \
	buf-lint \
	buf-generate

include .tools/buf/rules.mk
include .tools/clang-format/rules.mk

.PHONY: .build/protoc-gen-typescript-aip
.build/protoc-gen-typescript-aip:
	$(info [$@] building binary...)
	@cd ../.. && go build -o $(abspath $@) .

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto schemas...)
	@$(buf) lint

.PHONY: clang-format
clang-format: proto_files = $(shell find . -type f -name '*.proto')
clang-format: $(clang_format)
	$(info [$@] formatting proto files...)
	@$(clang_format) \
		-i \
		--style='{BasedOnStyle: Google, ColumnLimit: 0, Language: Proto}' \
		$(proto_files)

protoc_plugins := \
	.build/protoc-gen-typescript-aip

.PHONY: buf-generate
buf-generate: $(buf) $(protoc_plugins)
	$(info [$@] generating protobuf stubs...)
	@$(buf) generate \
		--template buf.gen.yaml \
		--path einride
