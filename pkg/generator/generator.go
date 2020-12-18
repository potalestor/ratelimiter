package generator

import (
	"context"
	"log"
	"sync"

	"guthub.com/potalestor/ratelimiter/pkg/command"
)

// Factory is factory for Commander.
type Factory func(i int) command.Commander

// Generator is used to generate fake commands.
type Generator struct {
	sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	comchan  chan command.Commander
	factory  Factory
	poolSize int
}

// NewGenerator returns new instance.
func NewGenerator(factory Factory, poolSize int) *Generator {
	ctx, cancel := context.WithCancel(context.Background())

	return &Generator{
		ctx:      ctx,
		cancel:   cancel,
		comchan:  make(chan command.Commander),
		factory:  factory,
		poolSize: poolSize,
	}
}

// Run and returns reading Commander channel.
func (g *Generator) Run() <-chan command.Commander {
	g.Add(g.poolSize)

	for i := 0; i < g.poolSize; i++ {
		go func(i int, g *Generator) {
			defer g.Done()
			log.Printf("[INFO ] The generator %d is running", i)

			for {
				select {
				case <-g.ctx.Done():
					log.Printf("[INFO ] The generator %d has finished working", i)

					return
				default:
					g.comchan <- g.factory(i)
				}
			}
		}(i, g)
	}

	return g.comchan
}

// Close and stop generator.
func (g *Generator) Close() error {
	g.cancel()
	g.Wait()
	close(g.comchan)

	return nil
}
