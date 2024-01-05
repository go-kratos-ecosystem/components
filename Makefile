.PHONY: lint
lint:
	golangci-lint run
	@echo "Linting complete"