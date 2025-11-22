.PHONY: up down restart logs build lint

up: ## Запустить проект в Docker
	docker-compose up -d

down: ## Остановить проект
	docker-compose down

restart: ## Перезапустить проект
	docker-compose restart

logs: ## Показать логи
	docker-compose logs -f

build: ## Пересобрать образ
	docker-compose build

rebuild: build up ## Пересобрать и запустить

lint: ## Проверить код линтером
	golangci-lint run ./...

