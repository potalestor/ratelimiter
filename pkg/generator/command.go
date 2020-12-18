package generator

import (
	"log"
	"math/rand"
	"time"

	"guthub.com/potalestor/ratelimiter/pkg/command"
)

// maximumCommandExecution is Maximum command execution time.
const maximumCommandExecution = time.Second

// Command is Commander implementation.
type Command struct {
	id int
}

// NewCommand returns new instance.
func NewCommand(id int) command.Commander {
	return &Command{
		id: id,
	}
}

// Execute runs command.
func (c *Command) Execute() {
	start := time.Now()
	time.Sleep(c.getDuration())
	elapsed := time.Since(start)
	log.Printf("[INFO ] The command %d was executed for %v", c.id, elapsed)
}

// getDuration returns random duration less 1.
func (c *Command) getDuration() time.Duration {
	return time.Duration(rand.Int63n(int64(maximumCommandExecution)))
}
