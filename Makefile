.PHONY: install-dev
install-dev:
	@echo "Installing..."
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

.PHONY: fmt
fmt:
	@echo "Formatting..."
	gci write -s standard -s default -s "prefix(github.com/cugu/uberfx)" .
	gofumpt -l -w .
	@echo "Done."

.PHONY: lint
lint:
	@echo "Linting..."
	golangci-lint run
	@echo "Done."
