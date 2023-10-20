package backend

import (
	"context"
	"gobridge/example/backend/second"
	"time"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role, u User, inventoryUpdate map[int64]bool) (bool, error)
	WhatsTheTime(ctx context.Context, date time.Time, toy second.Toy) (bool, error)
}
