package main

import (
	"context"
	"encoding/base64"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"

	"github.com/krestkrest/word-of-wisdom/api"
	"github.com/krestkrest/word-of-wisdom/internal/network/tcp"
)

var (
	address = flag.String("a", "", "Address of TCP server, required")
	level   = flag.String("l", "INFO", "log level")
	count   = flag.Int("c", 1, "requests count")
	force   = flag.Int("f", 0, "if non-zero, then the value of this flag is forced as a solution")
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	flag.Parse()

	if len(*address) == 0 {
		log.Fatal("Address of TCP server should be defined")
	}

	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Incorrect log level %s", *level)
	}
	log.SetLevel(logLevel)

	startTime := time.Now().UTC()

	ctx, cancel := context.WithCancel(context.Background())
	logger := log.WithContext(ctx)
	installTerminateHook(cancel)

	client := tcp.NewClient(*address)
	if err := client.Start(ctx); err != nil {
		log.Fatal("Failed to start TCP client", err)
	}
	logger.Debug("TCP client started")
	close := func() {
		if err := client.Close(); err != nil && !errors.Is(err, io.EOF) {
			logger.Warn("Error closing TCP connection", err)
		}
	}
	defer close()
	installTerminateHook(close)

	for i := 0; i < *count; i++ {
		logger.Debugf("Sending request #%d to server", i+1)

		reqType := uint8(api.MessageTypeRequest)
		if _, err := client.Write([]byte{reqType}); err != nil {
			logger.Fatal("Failed to send request", err)
		}

		if err := ensureReceivedMessageType(client, api.MessageTypeChallenge); err != nil {
			logger.Fatal("Error while receiving challenge from server", err)
		}

		var challenge api.MessageChallenge
		reader := msgp.NewReader(client)
		if err := challenge.DecodeMsg(reader); err != nil {
			logger.Fatal("Failed to unmarshal challenge message", err)
		}

		logger.Debugf("Received challenge from server, nonce: %s, complexity: %d",
			base64.StdEncoding.EncodeToString(challenge.Nonce[:]), challenge.Complexity)

		var solution uint64
		if *force != 0 {
			logger.Debugf("Forcing solution for challenge: %d", solution)
			solution = uint64(*force)
		} else {
			solution = challenge.FindSolution()
			logger.Debugf("Found solution for challenge: %d", solution)
		}

		logger.Debugf("Sending challenge response #%d to server", i+1)

		reqType = uint8(api.MessageTypeResponse)
		if _, err := client.Write([]byte{reqType}); err != nil {
			logger.Fatal("Failed to send challenge response", err)
		}

		response := api.NewMessageResponse(&challenge, solution)
		responseData, err := response.MarshalMsg(nil)
		if err != nil {
			logger.Fatal("Failed to marshal challenge response", err)
		}

		if _, err := client.Write(responseData); err != nil {
			logger.Fatal("Failed to send challenge response", err)
		}

		if err := ensureReceivedMessageType(client, api.MessageTypeGrant); err != nil {
			logger.Fatal("Error while receiving grant from server", err)
		}

		var grant api.MessageGrant
		reader = msgp.NewReader(client)
		if err := grant.DecodeMsg(reader); err != nil {
			logger.Fatal("Failed to unmarshal grant message", err)
		}

		if len(grant.Error) != 0 {
			logger.Errorf("Received error from server: %s", grant.Error)
		} else {
			logger.Infof("Quote: %s", grant.Quote)
		}
	}

	logger.Infof("Work completed, elapsed time: %s", time.Since(startTime))
}

func ensureReceivedMessageType(client io.Reader, expectedType api.MessageType) error {
	var messageTypeRaw [1]uint8
	if _, err := client.Read(messageTypeRaw[:]); err != nil {
		return errors.Wrap(err, "Failed to read message type")
	}
	messageType := api.MessageType(messageTypeRaw[0])
	if messageType != expectedType {
		return errors.Errorf("Incorrect message type, expected: %s, got: %s", expectedType, messageType)
	}
	return nil
}

func installTerminateHook(hook func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		hook()
	}()
}
