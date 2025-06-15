.PHONY: build run test docker-build docker-run clean help

build:
	go build -o stress-test ./cmd/stress

run:
	go run ./cmd/stress

docker-build:
	docker build -t stress-test .

http-up:
	docker-compose up -d test-server

http-down:
	docker-compose down

test-local:
	docker-compose run --rm stress-test --url=http://test-server/get --requests=50 --concurrency=5

help:
	@echo "Comandos disponíveis:"
	@echo "  build          - Compila a aplicação"
	@echo "  run            - Executa a aplicação localmente"
	@echo "  docker-build   - Constrói a imagem Docker"
	@echo "  docker-run-example - Executa um exemplo via Docker"
	@echo "  http-up        - Sobe servidor de teste local"
	@echo "  http-down      - Para o servidor de teste local"
	@echo "  test-local     - Executa teste contra servidor local"
	@echo "  help           - Mostra esta ajuda"
