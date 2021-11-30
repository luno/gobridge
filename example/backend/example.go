package backend

import (
	"context"
	"time"

	"gobridge/example/backend/second"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role, u User, inventoryUpdate map[int64]bool) (bool, error)
	WhatsTheTime(ctx context.Context, date time.Time, toy second.Toy) (bool, error)
}