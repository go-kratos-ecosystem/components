.PHONY: lint
lint:
	golangci-lint run
	@echo "✅ Linting completed"

.PHONY: fix
fix:
	golangci-lint run --fix
	@echo "✅ Lint fixing completed"