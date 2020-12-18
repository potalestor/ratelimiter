package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"guthub.com/potalestor/ratelimiter/pkg/config"
	"guthub.com/potalestor/ratelimiter/pkg/generator"
	"guthub.com/potalestor/ratelimiter/pkg/logger"
	"guthub.com/potalestor/ratelimiter/pkg/operator"
	"guthub.com/potalestor/ratelimiter/pkg/ratelimit"
)

const (
	nCommand     = "n"
	nDescription = "maximum number of parallel commands"
	nDefault     = 10

	xCommand     = "x"
	xDescription = "maximum number of commands per minute"
	xDefault     = 100

	logExt = ".log"

	generatorPoolSize = 1
)

func createLogfile() string {
	file := os.Args[0]

	return strings.TrimSuffix(file, filepath.Ext(file)) + logExt
}

func main() {
	// initialize logger
	logger.Initialize(createLogfile())

	// create config & validate input data
	n := flag.Int(nCommand, nDefault, nDescription)
	x := flag.Int(xCommand, xDefault, xDescription)

	flag.Parse()

	cfg := &config.Config{N: *n, X: *x}
	if err := cfg.Valid(); err != nil {
		log.Fatal(err)
	}

	// create & run test command generator
	g := generator.NewGenerator(generator.NewCommand, generatorPoolSize)
	cmds := g.Run()

	// create frame for ratelimit function
	ctx, cancel := context.WithCancel(context.Background())
	f := operator.NewFrame(ctx, cfg.N, cfg.X, time.Minute)

	// wait for interrupt signal to gracefully shutdown the app
	// kill -2 is syscall.SIGINT
	// kill (no param) default send syscall.SIGTERM
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// signal the frame to shutdown
		cancel()
		// signal the generator  to shutdown
		if err := g.Close(); err != nil {
			log.Printf("[ERROR] %v\n", err)
		}
	}()

	// run ratelimiter
	ratelimit.Execute(f, cmds)
}
