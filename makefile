BINARY_NAME=proto
BIN_TUI=tui-app

build-proto:
	go build -o ${BINARY_NAME} ./cmd/tui-proto

run-proto:
	go build -o ${BINARY_NAME} ./cmd/tui-proto
	./${BINARY_NAME}

clean-proto:
	go clean
	rm ${BINARY_NAME}
	sudo rm -r ~/.journi


build:
	go fmt ./cmd/tui
	go fmt ./internal/*.go
	go fmt ./internal/helpers
	go fmt ./internal/models
	go build -o ${BIN_TUI} ./cmd/tui

run:
	go fmt ./cmd/tui
	go fmt ./internal/*.go
	go fmt ./internal/helpers
	go fmt ./internal/models
	go build -o ${BIN_TUI} ./cmd/tui
	./${BIN_TUI}

clean:
	go clean
	rm ${BIN_TUI}
	sudo rm -r ~/.journi
