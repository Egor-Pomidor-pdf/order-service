# Order Service

Микросервис для управления заказами.

## 🚀 Быстрый старт

1.  **Скопируйте файл окружения:**
    ```bash
    cp .env.example .env
    ```

2.  **Отредактируйте `.env` файл** (при необходимости):
    ```bash
    # Установите свои значения для переменных
    POSTGRES_PASSWORD=your_strong_password_here
    # Остальные переменные можно оставить по умолчанию
    ```

3.  **Запустите приложение:**
    ```bash
    docker-compose up --build
    ```

4.  **Приложение будет доступно по адресу:** http://localhost:8080

## 📦 API endpoints

- `GET /order/<order_uid> ` - получить заказ

## 🛠 Технологии

- Go
- PostgreSQL
- Kafka
- Docker
- Docker Compose

## 🐛 Troubleshooting

Если возникают проблемы с подключением к БД:
```bash
docker-compose down -v
docker-compose up --build