package example

import (
	"context"
	"gobridge/example/second"
	"time"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role, u User) (bool, error)
	WhatsTheTime(ctx context.Context, time time.Time, toy second.Toy) (bool, error)
}