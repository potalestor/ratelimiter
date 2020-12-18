package generator_test

import (
	"testing"
	"time"

	"guthub.com/potalestor/ratelimiter/pkg/generator"
)

func TestGenerator_Run(t *testing.T) {
	g := generator.NewGenerator(generator.NewCommand, 10)
	commands := g.Run()

	go func() {
		for command := range commands {
			command.Execute()
		}
	}()

	time.Sleep(time.Second)
	g.Close()
}

func ExampleGenerator_Run() {
	g := generator.NewGenerator(generator.NewCommand, 10)
	commands := g.Run()

	go func() {
		for command := range commands {
			command.Execute()
		}
	}()

	time.Sleep(time.Second)
	g.Close()
}
