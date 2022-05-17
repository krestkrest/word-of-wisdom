package network

import (
	"context"
	"io"
)

//go:generate mockgen -source=network.go -destination=network_mock.go -package=network StreamHandler

// StreamHandler interface provides callback method to process data stream
type StreamHandler interface {
	HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser)
}
