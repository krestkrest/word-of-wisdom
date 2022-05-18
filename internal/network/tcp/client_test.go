//go:build integration

package tcp_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krestkrest/word-of-wisdom/internal/network"
	"github.com/krestkrest/word-of-wisdom/internal/network/tcp"
)

func TestClient_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := tcp.NewServer("127.0.0.1:3434", network.NewMockStreamHandler(ctrl))
	go func(t *testing.T) {
		assert.NoError(t, s.Start(ctx))
	}(t)

	time.Sleep(time.Second)

	c := tcp.NewClient("127.0.0.1:3434")
	require.NoError(t, c.Start(ctx))
	assert.NoError(t, c.Close())
}
