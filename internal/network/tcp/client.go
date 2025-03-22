package tcp

import (
	"log/slog"
)

type Client struct {
	logger *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{
		logger: logger,
	}
}
