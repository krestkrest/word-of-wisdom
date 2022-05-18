package network

import (
	"context"
	"io"
)

//go:generate mockgen -source=network.go -destination=network_mock.go -package=network StreamHandler

// StreamHandler provides callback method to process data stream
type StreamHandler interface {
	HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser)
}

// Server initializes server that receives incoming requests
type Server interface {
	Start(ctx context.Context) error
}
