package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

type UserEvent struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

type UserService interface {
	DeleteUser(id int) error
}

type ConsumerManager struct {
	consumersCount int
	consumers      []consumer
}

type consumer struct {
	userService UserService
	reader      *kafka.Reader
}

func NewConsumerManager(brokers string, userService UserService, group string, topic string, consumersCount int) *ConsumerManager {
	newConsumers := make([]consumer, consumersCount)
	for i := 0; i < consumersCount; i++ {
		newConsumers[i] = consumer{
			userService: userService,
			reader: kafka.NewReader(kafka.ReaderConfig{
				Brokers: strings.Split(brokers, ","),
				GroupID: group,
				Topic:   topic,
			}),
		}
	}
	return &ConsumerManager{
		consumers:      newConsumers,
		consumersCount: consumersCount,
	}
}

func (c *ConsumerManager) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for _, consumer := range c.consumers {
		wg.Add(1)
		go consumer.startConsumer(ctx, wg)
	}
	wg.Wait()
}

func (c *consumer) startConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer c.reader.Close()
	defer wg.Done()

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			slog.Error("fetch error", "err", err)
			continue
		}

		var event UserEvent
		if err := json.Unmarshal(m.Value, &event); err == nil {
			c.handleEvent(event)
			slog.Info("processed event", "user_id", event.UserID, "action", event.Action)
		} else {
			slog.Error("unmarshal error", "err", err)
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			slog.Error("commit error", "err", err)
		}
	}
}

func (c *consumer) handleEvent(event UserEvent) {
	switch event.Action {
	case "delete":
		c.Delete(event.UserID)
	default:
		slog.Warn("unknown action", "action", event.Action)
	}
}

func (c *consumer) Delete(userID int) {
	if err := c.userService.DeleteUser(userID); err != nil {
		slog.Error("delete failed", "user_id", userID, "err", err)
	}
}
