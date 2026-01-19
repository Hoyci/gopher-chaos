package chaos

import (
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

type Chaos struct {
	Config ChaosConfig
	pool   *sync.Pool
	logger *log.Logger
}

type Option func(*Chaos)

func NewChaos(cfg ChaosConfig, opts ...Option) *Chaos {
	c := &Chaos{
		Config: cfg,
		logger: log.Default(),
		pool: &sync.Pool{
			New: func() any {
				seed := uint64(time.Now().UnixNano())
				return rand.New(rand.NewPCG(seed, seed^0xdeadbeef))
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithLogger(l *log.Logger) Option {
	return func(c *Chaos) {
		if l != nil {
			c.logger = l
		}
	}
}

func (c *Chaos) ShouldInjectChaos() bool {
	r := c.pool.Get().(*rand.Rand)
	should := r.Float64() < c.Config.Probability
	c.pool.Put(r)

	if should {
		c.logger.Printf(
			"chaos injected (probability=%.2f)",
			c.Config.Probability,
		)
	}

	return should
}
