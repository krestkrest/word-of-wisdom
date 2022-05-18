package tcp

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/krestkrest/word-of-wisdom/internal/network"
)

const (
	keepAlivePeriod = 10 * time.Second
)

type Server struct {
	address string
	handler network.StreamHandler
}

var _ network.Server = (*Server)(nil)

func NewServer(address string, handler network.StreamHandler) *Server {
	return &Server{
		address: address,
		handler: handler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := log.WithContext(ctx)
	logger.Info("TCP server starting...")

	l := net.ListenConfig{KeepAlive: keepAlivePeriod}
	listener, err := l.Listen(ctx, "tcp", s.address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen TCP on address %s", s.address)
	}
	logger.Infof("Started TCP server on address %s", listener.Addr().String())

	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			logger.Warn("Error while closing TCP server: ", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Debug("Connection closed, quitting accept loop")
				return nil
			}

			return err
		}

		remoteAddress := conn.RemoteAddr().String()

		tcpConn, ok := conn.(*net.TCPConn)
		if ok {
			setupConnection(ctx, tcpConn)
		} else {
			logger.Warn("Failed to setup TCP connection options for address ", remoteAddress)
		}

		logger.Debug("Accepted new connection from ", remoteAddress)
		go s.handler.HandleStream(ctx, remoteAddress, conn)
	}
}

func setupConnection(ctx context.Context, conn *net.TCPConn) {
	logger := log.WithContext(ctx).WithField("address", conn.RemoteAddr())

	if err := conn.SetNoDelay(true); err != nil {
		logger.Warn("Failed to set connection no delay: ", err)
	}
	if err := conn.SetKeepAlivePeriod(keepAlivePeriod); err != nil {
		logger.Warn("Failed to set keep alive period: ", err)
	}
	if err := conn.SetKeepAlive(true); err != nil {
		logger.Warn("Failed to set keep alive: ", err)
	}
}
