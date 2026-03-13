package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type UserEvent struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

type UserDeleter interface {
	DeleteUser(id int) error
}

func StartConsumer(brokersStr string, userService UserDeleter) {
	brokers := strings.Split(brokersStr, ",")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "fth-delete-group",
		Topic:    "user-events",
		MaxWait:  1 * time.Second,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	slog.Info("Kafka Consumer started", "brokers", brokersStr, "topic", "user-events")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			slog.Error("Failed to read message from Kafka", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var event UserEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			slog.Error("Failed to unmarshal event", "error", err)
			continue
		}

		if event.Action == "deleted" {

			err := userService.DeleteUser(event.UserID)
			if err != nil {
				slog.Error("Failed to delete user in local DB", "userID", event.UserID, "error", err)
			}
		}
	}
}
