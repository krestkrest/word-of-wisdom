//go:build integration

package tcp_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/krestkrest/word-of-wisdom/internal/network"
	"github.com/krestkrest/word-of-wisdom/internal/network/tcp"
)

func TestServer_Start_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := tcp.NewServer("127.0.0.1:0", network.NewMockStreamHandler(ctrl))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.NoError(t, s.Start(ctx))
}
