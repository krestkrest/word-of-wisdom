package tcp

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var once sync.Once

const (
	dialTimeout = time.Second * 3
)

type Client struct {
	address    string
	connection net.Conn
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

var _ io.ReadWriteCloser = (*Client)(nil)

func (c *Client) Start(ctx context.Context) error {
	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: keepAlivePeriod,
	}

	conn, err := dialer.DialContext(ctx, "tcp", c.address)
	if err != nil {
		return errors.Wrapf(err, "failed to dial to address %s", c.address)
	}
	setupConnection(ctx, conn.(*net.TCPConn))

	c.connection = conn
	return nil
}

func (c *Client) Read(p []byte) (n int, err error) {
	if c.connection == nil {
		return 0, errors.New("TCP client is not started")
	}
	return c.connection.Read(p)
}

func (c *Client) Write(p []byte) (n int, err error) {
	if c.connection == nil {
		return 0, errors.New("TCP client is not started")
	}
	return c.connection.Write(p)
}

func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}
	var err error
	once.Do(func() {
		err = c.connection.Close()
	})
	return err
}
