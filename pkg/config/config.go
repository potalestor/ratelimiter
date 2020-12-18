package config

import (
	"errors"
	"fmt"
)

const fmtConfig = `
	Application settings:
	- maximum number of parallel commands (n): %d, 
	- maximum number of commands per minute (x): %d
`

var (
	// ErrInvalidN returns invalid N: N>0 & N<100.
	ErrInvalidN = errors.New("invalid N: N>0 & N<100")
	// ErrInvalidX returns invalid X: X>0.
	ErrInvalidX = errors.New("invalid X: X>0")
)

// Config contains the application settings.
type Config struct {
	// N is maximum number of parallel commands.
	N int
	// X is maximum number of commands per minute
	X int
}

// GetN is maximum number of parallel commands.
func (c *Config) GetN() int {
	return c.N
}

// GetX is maximum number of commands per minute.
func (c *Config) GetX() int {
	return c.X
}

// Valid returns an error in case of inconsistent data.
func (c *Config) Valid() error {
	if c.N < 0 || c.N > 100 {
		return ErrInvalidN
	}

	if c.X <= 0 {
		return ErrInvalidX
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf(fmtConfig, c.N, c.X)
}
