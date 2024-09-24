.PHONY: init
init:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/bufbuild/buf/cmd/buf@latest

.PHONY: lint
lint:
	golangci-lint run
	@echo "✅ Linting completed"

.PHONY: fix
fix:
	golangci-lint run --fix
	@echo "✅ Lint fixing completed"

.PHONY: buf-gen
	buf generate
	@echo "✅ Buf generation completed"

.PHONY: test
test:
	go test ./... -race -cover
	@echo "✅ Testing completed"

.PHONY: fmt
fmt:
	gofmt -w -e "vendor" .
	@echo "✅ Formatting completed"

.PHONY: fumpt
fumpt:
	gofumpt -w -e "vendor" .
	@echo "✅ Formatting completed"

.PHONY: nilaway-install
nilaway-install:
	go install go.uber.org/nilaway/cmd/nilaway@latest

.PHONY: nilaway
nilaway:
	nilaway ./...
	@echo "✅ Nilaway completed"