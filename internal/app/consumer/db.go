package consumer

import (
	"context"
	"github.com/ozonmp/edu-solution-api/internal/model"
	"github.com/ozonmp/edu-solution-api/internal/repo"
	"sync"
	"time"

)

type Consumer interface {
	Start(ctx context.Context)
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.SolutionEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.SolutionEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.SolutionEvent) Consumer {

	wg := &sync.WaitGroup{}

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
	}
}

func (c *consumer) Start(ctx context.Context) {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						event.Type = model.Updated
						c.events <- event
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	c.wg.Wait()
}
