package chaos

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type ChaosConfigLatency struct {
	Min time.Duration
	Max time.Duration
}

func (l ChaosConfigLatency) validate() error {
	if l.Min < 0 || l.Max < 0 {
		return fmt.Errorf(
			"%w: latency values must be non-negative",
			ErrInvalidLatency,
		)
	}

	if l.Min > l.Max {
		return fmt.Errorf(
			"%w: latency min (%s) must be <= max (%s)",
			ErrInvalidLatency,
			l.Min,
			l.Max,
		)
	}

	return nil
}

func (c *Chaos) CalculateLatency() time.Duration {
	min := c.Config.Latency.Min
	max := c.Config.Latency.Max

	if max <= min {
		return min
	}

	r := c.pool.Get().(*rand.Rand)
	randomValue := r.Int64N(int64(c.Config.Latency.Max - c.Config.Latency.Min))
	c.pool.Put(r)

	return c.Config.Latency.Min + time.Duration(randomValue)
}
