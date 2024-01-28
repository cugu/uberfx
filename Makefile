.PHONY: test_plugin
test_plugin:
	@echo "Building test plugin..."
	rm -rf ./example/wasi
	mkdir -p ./example/wasi
	GOOS=wasip1 GOARCH=wasm go build -o ./example/wasi/server.wasm ./example/minimal/
	@echo "Done."

.PHONY: fmt
fmt:
	@echo "Formatting..."
	go fmt ./...
	gci write -s standard -s default -s "prefix(github.com/cugu/uberfx)" .
	@echo "Done."

.PHONY: lint
lint:
	@echo "Linting..."
	golangci-lint run
	@echo "Done."
