/*
	signal.go
	Purpose: An internal pubsub system to listen and emit events.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/06  v1.0.0 Evan Chen   Initial release
*/

// Package signal is an internal pubsub system for internal events.
package signal

import (
	"sync"
)

type Event struct {
	Event string
	Args  []any
}

// Emit emits args to the given topic,
// the returned channel will be closed after every listener has received the event.
//
// The returned channel will not be closed if any of the listeners is busy
// and the event is still in the queue.
// Emit and On is sync but does not wait for actions to finish.
func Emit(topic string, args ...any) (sended <-chan struct{}) {
	done := make(chan struct{})

	c, ok := store[topic]
	wait := sync.WaitGroup{}
	if ok {
		c.lock.RLock()
		for i := range c.listeners {
			wait.Add(1)
			go func(index int) {
				defer func() {
					wait.Done()
					recover()
				}()
				c.listeners[index] <- Event{
					Event: topic,
					Args:  args,
				}
			}(i)
		}
		c.lock.RUnlock()
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	return done
}

// Handle is like On, but it handles the channel internally
func Handle(topic string, F func(e Event)) (close func()) {
	ch := On(topic)
	go func() {
		for evt := range ch {
			F(evt)
		}
	}()
	return func() { Off(topic, ch) }
}

// On returns a channel for receiving Events when a given topic is emitted.
func On(topic string) <-chan Event {
	store_lock.Lock()
	c, ok := store[topic]
	if !ok {
		c = &channel{
			listeners: make([]chan Event, 0),
		}
		store[topic] = c
	}
	store_lock.Unlock()

	c.lock.Lock()
	listener := make(chan Event)
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()

	return listener
}

// Off stops a given channel on the topic from receiving Events
func Off(topic string, channels ...<-chan Event) {
	store_lock.Lock()
	c, ok := store[topic]
	store_lock.Unlock()
	if !ok {
		return
	}

	c.lock.Lock()
	for i := range channels {
		for j := range c.listeners {
			if channels[i] == c.listeners[j] {
				close(c.listeners[j])
				c.listeners = remove(c.listeners, j)
				break
			}

		}
	}
	c.lock.Unlock()
}

// Remove closes every listener listening on the given topic
func Remove(topic string) {
	store_lock.Lock()
	c, ok := store[topic]
	store_lock.Unlock()
	if !ok {
		return
	}
	c.lock.Lock()
	for i := range c.listeners {
		close(c.listeners[i])
	}
	c.listeners = make([]chan Event, 0)
	c.lock.Unlock()
}

type channel struct {
	listeners []chan Event
	lock      sync.RWMutex
}

var store = map[string]*channel{}
var store_lock sync.Mutex

func remove(s []chan Event, i int) []chan Event {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
