package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

type state uint32

const (
	stateConfiguring state = iota
	stateOperational
	stateErrored
)

type PrepareState struct {
	state *atomic.Uint32
	cond  *sync.Cond
	mu    *sync.Mutex
}

func NewPrepareState() *PrepareState {
	var s atomic.Uint32
	s.Store(uint32(stateConfiguring))
	var mu sync.Mutex
	return &PrepareState{
		state: &s,
		cond:  sync.NewCond(&mu),
		mu:    &mu,
	}
}

func (s *PrepareState) SetOperational() {
	s.state.Store(uint32(stateOperational))
	s.cond.Broadcast()
}

func (s *PrepareState) SetErrored() {
	s.state.Store(uint32(stateErrored))

	s.cond.Broadcast()
}

func (s *PrepareState) WaitUntilOperational() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		switch state(s.state.Load()) {
		case stateOperational:
			return nil
		case stateErrored:
			return errors.New("internal error")
		case stateConfiguring:
			s.cond.Wait()
		}
	}
}
