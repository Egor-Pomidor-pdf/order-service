# Запустить все сервисы
up:
	docker-compose up --build

# Остановить все сервисы
down:
	docker-compose down

# Показать логи
logs:
	docker-compose logs -f