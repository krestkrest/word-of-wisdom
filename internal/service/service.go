package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"

	"github.com/krestkrest/word-of-wisdom/api"
	"github.com/krestkrest/word-of-wisdom/internal/network"
	"github.com/krestkrest/word-of-wisdom/internal/storage"
)

type Service struct {
	storage    storage.Storage
	complexity uint8
}

var _ network.StreamHandler = (*Service)(nil)

func NewService(storage storage.Storage, complexity uint8) *Service {
	return &Service{storage: storage, complexity: complexity}
}

func (s *Service) HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser) {
	logger := log.WithContext(ctx).WithField("address", address)

	if err := s.handleStream(ctx, address, stream); err != nil {
		if !errors.Is(err, io.EOF) {
			logger.Warn(err)
		}

		if err := stream.Close(); err != nil {
			logger.Warn("Failed to close TCP stream:", err)
		}
	}
}

func (s *Service) handleStream(ctx context.Context, address string, stream io.ReadWriteCloser) error {
	var nonce api.Nonce
	if err := generateNonce(nonce[:]); err != nil {
		return errors.Wrap(err, "Failed to generate nonce")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var messageTypeRaw [1]uint8
		if _, err := stream.Read(messageTypeRaw[:]); err != nil {
			return errors.Wrap(err, "Failed to read message type")
		}

		messageType := api.MessageType(messageTypeRaw[0])
		switch messageType {
		case api.MessageTypeRequest:
			r := api.NewMessageChallenge(address, s.complexity, nonce)
			if err := s.writeMessage(api.MessageTypeChallenge, r, stream); err != nil {
				return err
			}
		case api.MessageTypeResponse:
			var r api.MessageResponse
			reader := msgp.NewReader(stream)
			if err := r.DecodeMsg(reader); err != nil {
				return errors.Wrap(err, "Failed to unmarshal response message")
			}
			g := s.processResponse(&r, address, nonce)
			if err := s.writeMessage(api.MessageTypeGrant, g, stream); err != nil {
				return err
			}
			if err := generateNonce(nonce[:]); err != nil {
				return errors.Wrap(err, "Failed to generate nonce")
			}
		default:
			return errors.Errorf("Invalid message type: %d", messageType)
		}
	}
}

func (s *Service) writeMessage(messageType api.MessageType, message msgp.MarshalSizer, stream io.ReadWriteCloser) error {
	data := make([]byte, 0, message.Msgsize()+1)
	data = append(data, uint8(messageType))
	data, err := message.MarshalMsg(data)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal %s message", messageType)
	}
	if _, err := stream.Write(data); err != nil {
		return errors.Wrapf(err, "Failed to write %s message to TCP", messageType)
	}
	return nil
}

func (s *Service) processResponse(r *api.MessageResponse, address string, nonce api.Nonce) *api.MessageGrant {
	if err := s.validateResponse(r, address, nonce); err != nil {
		return &api.MessageGrant{Error: err.Error()}
	}
	return &api.MessageGrant{Quote: s.storage.GetQuote()}
}

func (s *Service) validateResponse(r *api.MessageResponse, address string, nonce api.Nonce) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.Address != address {
		return errors.New("address mismatch")
	}
	if bytes.Compare(r.Nonce[:], nonce[:]) != 0 {
		return errors.New("nonce mismatch")
	}
	if err := r.CheckSolution(); err != nil {
		return err
	}
	return nil
}

func generateNonce(dst []byte) error {
	_, err := rand.Read(dst)
	return err
}
