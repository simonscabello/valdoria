BINARY  := valdoria
PKG     := ./cmd/game
BIN_DIR := bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all run dev test build build-linux build-windows build-all fmt vet clean

all: build

## run: executa o jogo em modo release
run:
	go run $(PKG)

## dev: executa o jogo em modo de desenvolvimento
dev:
	VALDORIA_DEV=1 go run $(PKG)

## test: roda todos os testes
test:
	go test ./...

## build: compila para o sistema atual em bin/
build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(PKG)

## build-linux: compila o binário Linux (amd64)
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-linux-amd64 $(PKG)

## build-windows: compila o executável Windows (amd64, sem janela de console)
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS) -H windowsgui" -o $(BIN_DIR)/$(BINARY)-windows-amd64.exe $(PKG)

## build-all: compila para Linux e Windows
build-all: build-linux build-windows

## fmt: formata o código
fmt:
	gofmt -w .

## vet: análise estática
vet:
	go vet ./...

## clean: remove os artefatos de build
clean:
	rm -rf $(BIN_DIR)
