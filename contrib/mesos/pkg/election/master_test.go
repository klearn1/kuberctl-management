/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package election

import (
	"testing"
	"time"
)

type slowService struct {
	t  *testing.T
	on bool
	// We explicitly have no lock to prove that
	// Start and Stop are not called concurrently.
	changes chan<- bool
	done    <-chan struct{}
}

func (s *slowService) Validate(d, c Master) {
	// noop
}

func (s *slowService) Start() {
	select {
	case <-s.done:
		return // avoid writing to closed changes chan
	default:
	}
	if s.on {
		s.t.Errorf("started already on service")
	}
	time.Sleep(2 * time.Millisecond)
	s.on = true
	s.changes <- true
}

func (s *slowService) Stop() {
	select {
	case <-s.done:
		return // avoid writing to closed changes chan
	default:
	}
	if !s.on {
		s.t.Errorf("stopped already off service")
	}
	time.Sleep(2 * time.Millisecond)
	s.on = false
	s.changes <- false
}

func TestNotify(t *testing.T) {
	m := NewFake()
	changes := make(chan bool, 1500)
	done := make(chan struct{})
	s := &slowService{t: t, changes: changes, done: done}
	defer close(done)

	// change master to "notme" such that the initial m.Elect call inside Notify
	// will trigger an obversable event. We will wait for it to make sure the
	// Notify loop will see those master changes triggered by the go routine below.
	m.ChangeMaster(Master("me"))
	temporaryWatch := m.mux.Watch()
	ch := temporaryWatch.ResultChan()

	go Notify(m, "", "me", s, done)

	// wait for the event triggered by the initial m.Elect of Notify. Then drain
	// the channel to not block anything.
	<-ch
	temporaryWatch.Stop()
	for _ = range ch {
	}

	for i := 0; i < 500; i++ {
		for _, key := range []string{"me", "notme", "alsonotme"} {
			m.ChangeMaster(Master(key))
		}
	}

	// elections that don't include "me" don't trigger change events
	const want = 1000
	timeout := time.After(100 * time.Millisecond)
	got := 0

outer:
	for ; got < want; got++ {
		select {
		case <-changes:
		case <-timeout:
			break outer
		}
	}

	if got != want {
		t.Errorf("unexpected number of changes: got, want: %v, %v", got, want)
	}
}
