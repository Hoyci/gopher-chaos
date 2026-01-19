package interceptors

import (
	"context"
	"time"

	"github.com/hoyci/gopher-chaos/pkg/chaos"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	Chaos *chaos.Chaos
}

func NewInterceptor(c *chaos.Chaos) *Interceptor {
	return &Interceptor{
		Chaos: c,
	}
}

func (I *Interceptor) inject(ctx context.Context) error {
	if !I.Chaos.ShouldInjectChaos() {
		return nil
	}

	latency := I.Chaos.CalculateLatency()
	timer := time.NewTimer(latency)
	defer timer.Stop()

	select {
	case <-timer.C:
		return status.Error(
			I.Chaos.Config.Error,
			"chaos injected error",
		)

	case <-ctx.Done():
		return ctx.Err()
	}
}
