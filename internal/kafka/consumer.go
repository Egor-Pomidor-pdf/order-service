package kafka

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Egor-Pomidor-pdf/order-service/internal/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeout = 10000
	timeout        = 5000
)

type Handler interface {
	HandleMessage(message []byte, offset kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(handler Handler, cfg config.KafkaConfig) (*Consumer, error) {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.Brokers, ","),
		"group.id":                 cfg.GroupID,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       false,
		"auto.offset.reset":        "latest",
	}

	c, err := kafka.NewConsumer(kafkaConfig)

	if err != nil {
		return nil, err
	}

	if err := c.Subscribe(cfg.Topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil

}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(timeout)

		if err != nil {
			if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrTimedOut {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			slog.Error("Error reading message", "error", err)
		}
		if kafkaMsg == nil {
			continue
		}

		if err := c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
			if strings.Contains(err.Error(), "VALIDATION_ERROR:") ||
				strings.Contains(err.Error(), "INVALID_JSON:") {
				slog.Error("SKIPPING_INVALID_MESSAGE",
					"error", err,
					"raw_message", string(kafkaMsg.Value),
				)
				c.consumer.StoreMessage(kafkaMsg)
				c.consumer.Commit()
				continue
			} else {
				slog.Error("DATABASE_ERROR - WILL RETRY", "error", err)
				continue
			}
		}

		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			slog.Error("Error storing message offset", "error", err)
			continue
		}

		if _, err := c.consumer.Commit(); err != nil {
			slog.Error("Error committing offset", "error", err)
		}

		slog.Info("Message processed successfully", "offset", kafkaMsg.TopicPartition.Offset)

	}
}

func (c *Consumer) Stop() error {
	c.stop = true

	if _, err := c.consumer.Commit(); err != nil {
		slog.Error("failed to commit offsets on stop", "error", err)
	}

	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}
