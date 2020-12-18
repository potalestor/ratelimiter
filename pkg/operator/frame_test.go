package operator_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"guthub.com/potalestor/ratelimiter/pkg/command"
	"guthub.com/potalestor/ratelimiter/pkg/generator"
	"guthub.com/potalestor/ratelimiter/pkg/operator"
)

func TestNewFrame(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	g := generator.NewGenerator(generator.NewCommand, 5)
	commands := g.Run()

	go func(ctx context.Context, commands <-chan command.Commander) {
		op := operator.NewFrame(ctx, 3, 3, time.Minute)

		for cmd := range commands {
			if err := op.Take(); err != nil {
				// save no handeled commands to ...
				fmt.Println(err)
			} else {
				go func(cmd command.Commander, op operator.Operator) {
					cmd.Execute()
					op.Give()
				}(cmd, op)
			}
		}

		fmt.Println("gourutine finished")
	}(ctx, commands)

	time.Sleep(5 * time.Second)

	cancel()

	if err := g.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)
}
