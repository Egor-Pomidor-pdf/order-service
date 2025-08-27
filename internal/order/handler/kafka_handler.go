package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Egor-Pomidor-pdf/order-service/internal/order"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
)




var validate = validator.New()

func (h *OrderHandler) HandleMessage(message []byte, offset kafka.Offset) error {
	var order order.Order

	if err := json.Unmarshal(message, &order); err != nil {
        return fmt.Errorf("INVALID_JSON: %w", err)
    }
	if err := validate.Struct(order); err != nil {
		return fmt.Errorf("VALIDATION_ERROR: %w", err)
	}

	slog.Info("processing order from Kafka", "uid", order.OrderUID, "offset", offset)
	if err := h.service.ProcessOrder(context.Background(), order); err != nil {
        return fmt.Errorf("DATABASE_ERROR: %w", err) 
    }
	return nil
}
