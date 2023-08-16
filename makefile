BINARY_NAME=proto
BIN_TUI=tuirnal

build-proto:
	go build -o ${BINARY_NAME} ./cmd/tui-proto

run-proto:
	go build -o ${BINARY_NAME} ./cmd/tui-proto
	./${BINARY_NAME}

clean-proto:
	go clean
	rm ${BINARY_NAME}
	sudo rm -r ~/.tuirnal


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
	sudo rm ~/.tuirnal/tuirnal.db

clean-r:
	go clean
	rm ${BIN_TUI}
	sudo rm -r ~/.tuirnal
