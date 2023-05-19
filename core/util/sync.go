/*
	sync.go
	Purpose: sync utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

import (
	"fmt"
	"sync"
	"time"
)

type MutexMap struct {
	mut     sync.RWMutex           // handle concurrent access of chanMap
	chanMap map[string](chan bool) // dynamic mutexes map
}

func NewMutextMap() *MutexMap {
	return &MutexMap{
		mut:     sync.RWMutex{},
		chanMap: make(map[string](chan bool)),
	}
}

// Acquire a lock, with timeout
func (mm *MutexMap) Lock(id string, timeout ...time.Duration) error {
	var wait time.Duration = 100 * time.Millisecond
	if len(timeout) > 0 {
		wait = timeout[0]
	}

	// get global lock to read from map and get a channel
	mm.mut.Lock()
	if _, ok := mm.chanMap[id]; !ok {
		mm.chanMap[id] = make(chan bool, 1)
	}
	ch := mm.chanMap[id]
	mm.mut.Unlock()

	// try to write to buffered channel, with timeout
	select {
	case ch <- true:
		return nil
	case <-time.After(wait):
		return fmt.Errorf("timeout waiting for %v", id)
	}
}

// release lock
func (mm *MutexMap) Release(id string) {
	mm.mut.Lock()
	ch := mm.chanMap[id]
	mm.mut.Unlock()
	<-ch
}

func (mm *MutexMap) Clean() {
	mm.mut.Lock()
	for k, v := range mm.chanMap {
		select {
		case v <- true:
			close(v)
			delete(mm.chanMap, k)
		default:
		}
	}
	mm.mut.Unlock()
}
