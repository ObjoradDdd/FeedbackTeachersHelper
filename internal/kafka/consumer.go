package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

type UserService interface {
	DeleteUser(ctx context.Context, id int) error
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

func (c *ConsumerManager) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, c.consumersCount)
	internalCtx, stopAllConsumers := context.WithCancel(ctx)
	defer stopAllConsumers()

	for _, cons := range c.consumers {
		wg.Add(1)
		go func(c consumer) {
			defer wg.Done()
			if err := c.startConsumer(internalCtx); err != nil {
				errChan <- err
			}
		}(cons)
	}

	select {
	case err := <-errChan:
		stopAllConsumers()
		wg.Wait()
		return err
	case <-ctx.Done():
		stopAllConsumers()
		wg.Wait()
		return nil
	}
}

func (c *consumer) startConsumer(ctx context.Context) error {
	defer c.reader.Close()

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Error("fetch error", "err", err)
			return fmt.Errorf("kafka fetch error: %w", err)
		}

		var event UserEvent
		if err := json.Unmarshal(m.Value, &event); err == nil {
			if err := c.handleEvent(ctx, event); err != nil {
				slog.Error("handle event error", "err", err)
				return fmt.Errorf("handle event error: %w", err)
			}
			slog.Info("processed event", "user_id", event.UserID, "action", event.Action)
		} else {
			slog.Error("unmarshal error", "err", err)
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			slog.Error("commit error", "err", err)
			return fmt.Errorf("kafka commit error: %w", err)
		}
	}
}

func (c *consumer) handleEvent(ctx context.Context, event UserEvent) error {
	switch event.Action {
	case "delete":
		if err := c.Delete(ctx, event.UserID); err != nil {
			return err
		}

	default:
		slog.Warn("unknown action", "action", event.Action)
	}
	return nil
}
