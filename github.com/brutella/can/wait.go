/*
 * Copyright (C) 2020. Niklas Linz - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the LGPLv3 license, which unfortunately won't be
 * written for another century.
 *
 * You should have received a copy of the LGPLv3 license with
 * this file. If not, please write to: niklas.linz@enigmar.de
 *
 */

package can

import (
	"fmt"
	"time"
)

// A WaitResponse encapsulates the response of waiting for a frame.
type WaitResponse struct {
	Frame Frame
	Err   error
}

type waiter struct {
	id     uint32
	wait   chan WaitResponse
	bus    *Bus
	filter Handler
}

// Wait returns a channel, which receives a frame or an error, if the
// frame with the expected id didn't arrive on time.
func Wait(bus *Bus, id uint32, timeout time.Duration) <-chan WaitResponse {
	waiter := waiter{
		id:   id,
		wait: make(chan WaitResponse),
		bus:  bus,
	}

	ch := make(chan WaitResponse)

	go func() {
		select {
		case resp := <-waiter.wait:
			ch <- resp
		case <-time.After(timeout):
			err := fmt.Errorf("Timeout error waiting for %X", id)
			ch <- WaitResponse{Frame{}, err}
		}
	}()

	waiter.filter = newFilter(id, &waiter)
	bus.Subscribe(waiter.filter)

	return ch
}

func (w *waiter) Handle(frame Frame) {
	w.bus.Unsubscribe(w.filter)
	w.wait <- WaitResponse{frame, nil}
}
