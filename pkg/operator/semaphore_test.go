package operator_test

import (
	"fmt"
	"testing"
	"time"

	"guthub.com/potalestor/ratelimiter/pkg/command"
	"guthub.com/potalestor/ratelimiter/pkg/generator"
	"guthub.com/potalestor/ratelimiter/pkg/operator"
)

func TestNewSemaphore(t *testing.T) {
	g := generator.NewGenerator(generator.NewCommand, 5)
	commands := g.Run()

	go func(commands <-chan command.Commander) {
		op := operator.NewSemaphore(3)
		defer op.Close()

		for cmd := range commands {
			if err := op.Take(); err != nil {
				t.Error(err)
			}

			go func(cmd command.Commander, op operator.Operator) {
				cmd.Execute()
				op.Give()
			}(cmd, op)
		}

		fmt.Println("gourutine finished")
	}(commands)

	time.Sleep(3 * time.Second)

	if err := g.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)
}
