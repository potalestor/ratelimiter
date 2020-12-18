package ratelimit

import (
	"log"

	"guthub.com/potalestor/ratelimiter/pkg/command"
	"guthub.com/potalestor/ratelimiter/pkg/operator"
)

// Execute with a rate limit.
func Execute(op operator.Operator, cmds <-chan command.Commander) {
	log.Println("[Info] start executing...")
	defer log.Println("[Info] cancel execute")

	for cmd := range cmds {
		if err := op.Take(); err != nil {
			// save no handeled commands to ...
			log.Printf("[Error] %s\n", err)
		} else {
			go func(cmd command.Commander, op operator.Operator) {
				cmd.Execute()
				op.Give()
			}(cmd, op)
		}
	}
}
