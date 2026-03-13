package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/segmentio/kafka-go"
)

type UserEvent struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

type UserService interface {
	DeleteUser(id int) error
}

type UserConsumer struct {
	reader      *kafka.Reader
	userService UserService
}

func NewUserConsumer(brokers string, userService UserService, group string, topic string) *UserConsumer {
	return &UserConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  strings.Split(brokers, ","),
			GroupID:  group,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		userService: userService,
	}
}

func (c *UserConsumer) Start() {
	defer c.reader.Close()

	for {
		m, err := c.reader.FetchMessage(context.Background())
		if err != nil {
			if context.Background().Err() != nil {
				return
			}
			slog.Error("fetch error", "err", err)
			continue
		}

		var event UserEvent
		if err := json.Unmarshal(m.Value, &event); err == nil {
			if event.Action == "delete" {
				if err := c.userService.DeleteUser(event.UserID); err != nil {
					slog.Error("delete failed", "user_id", event.UserID, "err", err)
					continue
				}
			}
		} else {
			slog.Error("unmarshal error", "err", err)
		}

		if err := c.reader.CommitMessages(context.Background(), m); err != nil {
			slog.Error("commit error", "err", err)
		}
	}
}
