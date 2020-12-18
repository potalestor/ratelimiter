package generator_test

import (
	"testing"

	"guthub.com/potalestor/ratelimiter/pkg/generator"
)

func TestCommand_Execute(t *testing.T) {
	c := &generator.Command{}

	for i := 0; i < 3; i++ {
		c.Execute()
	}
}
