package kafka

import (
	"context"
	"log/slog"
)

type UserEvent struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

func (c *consumer) Delete(ctx context.Context, userID int) error {
	if err := c.userService.DeleteUser(ctx, userID); err != nil {
		slog.Error("delete failed", "user_id", userID, "err", err)
		return err
	}
	return nil
}
