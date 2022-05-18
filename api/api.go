package api

import (
	"time"
)

//go:generate msgp

const (
	nonceSize = 20
)

type Nonce [nonceSize]byte

// MessageType is the type of messages exchanged between client and server
type MessageType uint8

const (
	MessageTypeRequest   MessageType = iota // request challenge from server
	MessageTypeChallenge                    // challenge data returned from server
	MessageTypeResponse                     // response from client with solved challenge
	MessageTypeGrant                        // requested information granted from server
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeRequest:
		return "request"
	case MessageTypeChallenge:
		return "challenge"
	case MessageTypeResponse:
		return "response"
	case MessageTypeGrant:
		return "grant"
	default:
		return "unknown"
	}
}

type MessageChallenge struct {
	Address    string `msg:"a"` // client's address
	Nonce      Nonce  `msg:"n"` // random string
	UnixTime   int64  `msg:"t"` // current time of server in unix nano
	Complexity uint8  `msg:"c"` // number of bits to be zeroed
}

type MessageResponse struct {
	MessageChallenge
	Solution uint64 `msg:"s"` // additional number applied to hash
}

type MessageGrant struct {
	Quote string `msg:"q"` // quote is filled if the solution is correct
	Error string `msg:"e"` // error is filled if any errors occurred
}

func NewMessageChallenge(address string, complexity uint8, nonce Nonce) *MessageChallenge {
	result := &MessageChallenge{
		Address:    address,
		UnixTime:   time.Now().UTC().UnixNano(),
		Complexity: complexity,
	}
	copy(result.Nonce[:], nonce[:])
	return result
}

func NewMessageResponse(challenge *MessageChallenge, solution uint64) *MessageResponse {
	return &MessageResponse{
		MessageChallenge: *challenge,
		Solution:         solution,
	}
}
