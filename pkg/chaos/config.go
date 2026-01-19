package chaos

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidProbability = errors.New("invalid probability")
	ErrInvalidLatency     = errors.New("invalid latency")
	ErrInvalidErrorCode   = errors.New("invalid grpc error code")
)

type ChaosConfig struct {
	Probability float64
	Latency     ChaosConfigLatency
	Error       codes.Code
}

func (c ChaosConfig) Validate() error {
	if c.Probability < 0 || c.Probability > 1 {
		return fmt.Errorf(
			"%w: probability must be between 0 and 1, got %f",
			ErrInvalidProbability,
			c.Probability,
		)
	}

	if err := c.Latency.validate(); err != nil {
		return err
	}

	if c.Error < codes.OK || c.Error > codes.Unauthenticated {
		return fmt.Errorf(
			"%w: invalid grpc code %d",
			ErrInvalidErrorCode,
			c.Error,
		)
	}

	return nil
}
