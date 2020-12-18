package operator

import (
	"context"
	"log"
	"time"
)

type state int

const (
	initState state = iota
	workState
	waitState
	abortState
	doneState
)
const fmtStat = `[INFO ] wait action: sleep %v executed %d command`

type actionFunc func() (nextState state, err error)

// Frame with semaphore.
type Frame struct {
	Semaphore
	// ctx is context close instance.
	ctx context.Context
	// lifetime of sliding window.
	lifetime time.Duration
	// max for the lifetime.
	max int
	// counter value for current sliding window.
	counter int
	// deadline for current sliding window.
	deadline time.Time
	// current state is initState.
	current state
	// actions
	actions []actionFunc
}

// NewFrame returns new instance.
func NewFrame(ctx context.Context, size, max int, lifetime time.Duration) *Frame {
	s := &Frame{
		Semaphore{
			bufchan: make(chan struct{}, size),
		},
		ctx,
		lifetime,
		max,
		0,
		time.Time{},
		initState,
		nil,
	}
	s.actions = []actionFunc{
		s.initAction,
		s.workAction,
		s.waitAction,
		s.abortAction,
		s.doneAction,
	}

	return s
}

// Take writes struct to channel with sliding window algorithm.
func (s *Frame) Take() error {
	var err error

	st := s.current

	for {
		st, err = s.actions[st]()
		if err != nil {
			break
		}

		if st == doneState {
			s.current = workState

			return nil
		}
	}

	s.current = st

	return err
}

// flags ...

func (s *Frame) isDeadline() bool {
	return time.Since(s.deadline) > 0
}

func (s *Frame) isMax() bool {
	return s.counter >= s.max
}

// actions ...

func (s *Frame) initAction() (state, error) {
	log.Printf("[INFO ] init action: executed %d command", s.counter)
	s.deadline = time.Now().Add(s.lifetime)
	s.counter = 0

	return workState, nil
}

func (s *Frame) workAction() (state, error) {
	if s.isDeadline() {
		return initState, nil
	}

	if s.isMax() {
		return waitState, nil
	}

	s.counter++

	return doneState, s.Semaphore.Take()
}

func (s *Frame) waitAction() (state, error) {
	remain := time.Until(s.deadline)

	log.Printf(fmtStat, remain, s.counter)

	for {
		select {
		case <-s.ctx.Done():
			s.Close()

			return abortState, context.Canceled

		case <-time.After(remain):
			return initState, nil
		}
	}
}

func (s *Frame) abortAction() (state, error) {
	return abortState, context.Canceled
}

func (s *Frame) doneAction() (state, error) {
	return workState, nil
}
