.PHONY: lint
lint:
	golangci-lint run
	@echo "Linting complete"

.PHONY: fix
fix:
	golangci-lint run --fix
	@echo "Lint fixing complete"