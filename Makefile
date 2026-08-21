.PHONY: help build install uninstall clean test tidy fmt verify verify-help path doctor run setup

GOPATH_BIN := $(shell go env GOPATH)/bin
BINARY     := bin/vibeit
INSTALLED  := $(GOPATH_BIN)/vibeit
export PATH := $(GOPATH_BIN):$(PATH)

.DEFAULT_GOAL := help

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*## "}; \
	/^[a-zA-Z0-9_.-]+:.*## / {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}' \
	$(MAKEFILE_LIST)
	@echo ""
	@echo "  Quick start:  source ./install.sh"
	@echo "  GOPATH bin:   $(GOPATH_BIN)"

setup: install ## Alias — one-step install (same as install)

install: ## Build + install + PATH (~/.zshrc). Run: source ./install.sh
	@chmod +x install.sh scripts/env.sh scripts/verify.sh scripts/install-path.sh 2>/dev/null || chmod +x install.sh
	@./install.sh

build: ## Build ./bin/vibeit only
	go build -o $(BINARY) ./cmd/vibeit

uninstall: ## Remove vibeit from $$(go env GOPATH)/bin
	@if [ -f "$(INSTALLED)" ]; then \
		rm -f "$(INSTALLED)"; \
		echo "Removed $(INSTALLED)"; \
	else \
		echo "Not installed: $(INSTALLED)"; \
	fi

clean: ## Remove ./bin and Go build cache
	rm -rf bin/
	go clean ./...

test: ## Run go test ./...
	go test ./...

tidy: ## go mod tidy
	go mod tidy

fmt: ## gofmt all Go sources
	gofmt -w cmd internal templates

path: ## Print export PATH=... for your shell
	@echo 'export PATH="$(GOPATH_BIN):$$PATH"'

doctor: ## Check go, node, npm, flutter, vibeit on PATH
	@echo "== vibeit doctor =="
	@echo "GOPATH bin: $(GOPATH_BIN)"
	@for tool in go node npm flutter; do \
		if command -v $$tool >/dev/null 2>&1; then \
			printf "  OK  %-8s %s\n" "$$tool" "$$(command -v $$tool)"; \
		else \
			printf "  --  %-8s not found\n" "$$tool"; \
		fi; \
	done
	@if command -v vibeit >/dev/null 2>&1; then \
		printf "  OK  %-8s %s\n" "vibeit" "$$(command -v vibeit)"; \
	elif [ -x "./$(BINARY)" ]; then \
		printf "  OK  %-8s ./$(BINARY) (run: source ./install.sh)\n" "vibeit"; \
	else \
		printf "  --  %-8s not installed — run: source ./install.sh\n" "vibeit"; \
	fi

run: build ## Run ./bin/vibeit (ARGS="create demo --backend")
	./$(BINARY) $(ARGS)

verify: ## Smoke test in $$TMPDIR (auto-deleted; does NOT touch cwd)
	@chmod +x scripts/env.sh scripts/verify.sh install.sh
	@$(MAKE) build
	@./scripts/env.sh 2>/dev/null || true
	@command -v vibeit >/dev/null 2>&1 || export PATH="$(GOPATH_BIN):$$PATH"; \
	 chmod +x scripts/verify.sh; ./scripts/verify.sh

verify-help: build ## Print vibeit CLI help
	./$(BINARY) --help
	./$(BINARY) create --help
	./$(BINARY) add --help
	./$(BINARY) skills update --help
