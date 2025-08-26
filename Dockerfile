FROM golang:1.24.6

WORKDIR /app

RUN apt-get update && apt-get install -y librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service ./cmd/app

WORKDIR /app

CMD ["./order-service"]
