package example

import "context"

type Example interface {
	Name(ctx context.Context) (string, error)
}