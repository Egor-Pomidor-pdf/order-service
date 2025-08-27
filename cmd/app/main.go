package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/Egor-Pomidor-pdf/order-service/internal/config"
	"github.com/Egor-Pomidor-pdf/order-service/internal/db"
	"github.com/Egor-Pomidor-pdf/order-service/internal/kafka"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/handler"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/repository"
	"github.com/Egor-Pomidor-pdf/order-service/internal/order/service"
	"github.com/Egor-Pomidor-pdf/order-service/internal/server"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
    // Initialize environment variables
    _ = godotenv.Load() 
	cfg := config.MustLoad()
    setupLogger(cfg.Env) 

	slog.Info("starting app", slog.String("env", cfg.Env))
	slog.Debug("debug messages are enabled")


    // Подключение к базе данных
	db, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

    // Инициализация репозиториев, сервисов и HTTP сервера
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
    orderHandler := handler.NewOrderHandler(orderService)
    srv := server.NewServer(cfg.Server, orderService)

    // Запуск HTTP сервера в отдельной горутине
    go func() {
        slog.Info("starting HTTP server", slog.String("address", cfg.Server.Address))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("failed to start server", slog.String("error", err.Error()))
            os.Exit(1)
        }
    }()
    
    // Инициализация и запуск Kafka consumer
    consumer, err := kafka.NewConsumer(orderHandler, cfg.Kafka)
    if err != nil {
        slog.Error("failed to create Kafka consumer", slog.String("error", err.Error()))
        return
    }
    go consumer.Start()

    // Ожидание сигнала завершения (например, SIGINT или SIGTERM)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    slog.Info("received shutdown signal, shutting down...")

    // Остановка Kafka consumer
    if err := consumer.Stop(); err != nil {
        slog.Error("failed to stop Kafka consumer", slog.String("error", err.Error()))
    }

    // Завершение работы HTTP сервера
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("failed to shutdown server", slog.String("error", err.Error()))
    }
    defer cancel()

    // Закрытие подключения к базе данных
    if err := db.Close();err != nil {
        slog.Error("failed to close database", slog.String("error", err.Error()))
    }
    
    slog.Info("shutting down application")
}

// setupLogger configures the global logger based on the environment
func setupLogger(env string) {
    var handler slog.Handler
    switch env {
    case envLocal:
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    case envDev:
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    case envProd:
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    }
    slog.SetDefault(slog.New(handler))  // ← Устанавливаем global logger
}
