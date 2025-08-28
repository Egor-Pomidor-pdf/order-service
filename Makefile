# Сборка
build:
	docker-compose up --build

# Запуск
up:
	docker-compose up -d

# Остановить все сервисы
down:
	docker-compose down

# Показать логи
logs:
	docker-compose logs -f

# Тесты с покрытием
test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

# Показать покрытие тестов
test.coverage:
	go tool cover -func=cover.out | grep "total"

test-e2e:
	docker-compose up -d 
	sleep 10                
	go test -v ./tests/ -run Test_OrderFlow -timeout 30s 

# отправка тестового сообщения в кафку
send-order:	
	go run './cmd/producer-test/producer.go'
