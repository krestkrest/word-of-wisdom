package api

import (
	"time"

	"github.com/pkg/errors"
)

const (
	maxAddressLen   = 45
	maxCalcDuration = time.Second * 10
)

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
	panic("")
}
