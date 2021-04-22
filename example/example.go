package example

import (
	"context"
)

type Example interface {
	HasPermission(ctx context.Context, r []Role) (bool, error)
}