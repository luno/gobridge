package backend

import (
	"context"
	"time"

	"gobridge/example/backend/second"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role, u User) (bool, error)
	WhatsTheTime(ctx context.Context, time time.Time, toy second.Toy) (bool, error)
}