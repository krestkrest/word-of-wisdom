package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"math/bits"
	"time"

	"github.com/pkg/errors"
)

const (
	maxAddressLen   = 45
	maxCalcDuration = time.Second * 10
)

func (m *MessageChallenge) Check(b *bytes.Buffer, solution uint64) bool {
	if b == nil {
		b = &bytes.Buffer{}
	}
	m.concatData(b, solution)
	hash := sha1.Sum(b.Bytes())
	index := len(hash) - 1
	zeroBits := m.Complexity
	for ; index >= 0 && zeroBits >= 8; index-- {
		if hash[index] != 0 {
			return false
		}
		zeroBits -= 8
	}
	if zeroBits > 0 {
		return bits.TrailingZeros8(hash[index]) >= int(zeroBits)
	}
	return true
}

func (m *MessageChallenge) concatData(b *bytes.Buffer, solution uint64) {
	b.WriteString(m.Address)
	b.Write(m.Nonce[:])

	i := make([]byte, 8)
	binary.LittleEndian.PutUint64(i, uint64(m.UnixTime))
	b.Write(i)

	b.Write([]byte{m.Complexity})

	binary.LittleEndian.PutUint64(i, solution)
	b.Write(i)
}

func (m *MessageChallenge) FindSolution() uint64 {
	var b bytes.Buffer
	var solution uint64
	for {
		if m.Check(&b, solution) {
			return solution
		}
		b.Reset()
		solution += 1
	}
}

func (m *MessageResponse) Validate() error {
	if len(m.Address) > maxAddressLen {
		return errors.New("address length is greater than max limit for IPv6")
	}
	t := time.Unix(0, m.UnixTime).UTC()
	elapsed := time.Since(t)
	if elapsed < 0 || elapsed > maxCalcDuration {
		return errors.New("incorrect time")
	}
	return nil
}

func (m *MessageResponse) CheckSolution() error {
	if !m.Check(nil, m.Solution) {
		return errors.New("incorrect solution")
	}
	return nil
}
