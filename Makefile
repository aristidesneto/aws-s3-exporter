.PHONY: build run test clean lint help

# Variáveis
APP_NAME = aws-exporter
GO = go
BIN_DIR = bin

help:
	@echo "Comandos disponíveis:"
	@echo "  make build    - Compila o projeto"
	@echo "  make run      - Executa a aplicação"
	@echo "  make test     - Executa os testes"
	@echo "  make clean    - Remove os binários gerados"
	@echo "  make lint     - Executa o linter (golangci-lint)"
	@echo "  make fmt      - Formata o código fonte"

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) ./cmd/aws-exporter

run: build
	./$(BIN_DIR)/$(APP_NAME)

test:
	$(GO) test -v ./...

clean:
	rm -rf $(BIN_DIR)
	$(GO) clean

lint:
	@if ! command -v ~/go/bin/golangci-lint &> /dev/null; then \
		echo "golangci-lint não encontrado. Instalando..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	~/go/bin/golangci-lint run

fmt:
	$(GO) fmt ./...

# Comandos para desenvolvimento
dev: fmt lint build run
