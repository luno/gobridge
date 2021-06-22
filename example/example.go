package example

import (
	"context"
	"time"

	"gobridge/example/second"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role, u User) (bool, error)
	WhatsTheTime(ctx context.Context, time time.Time, toy second.Toy) (bool, error)
}