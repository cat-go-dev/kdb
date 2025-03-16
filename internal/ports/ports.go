package ports

import "context"

type Database interface {
	Execute(ctx context.Context, commandStr string) (*Result, error)
}

type Result struct {
	Msg string
}
