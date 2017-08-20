// MIT License
//
// Copyright (c) 2017 Sigurd Spieckermann
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package nsync

import (
	"sync"
)

// NamedRWMutex implements a map of named read-write mutexes. A read-write mutex
// is referenced by a name, so mutexes can be created dynamically.
type NamedRWMutex struct {
	mutex    sync.Mutex
	mutexMap map[string]*rwMutexEntry
}

type rwMutexEntry struct {
	mutex sync.RWMutex
	num   int
}

// NewNamedRWMutex creates a new map of named read-write mutexes.
func NewNamedRWMutex() *NamedRWMutex {
	return &NamedRWMutex{mutexMap: make(map[string]*rwMutexEntry)}
}

// Lock acquires a write lock for name. It blocks until the lock has been
// acquired.
func (n *NamedRWMutex) Lock(name string) {
	n.mutex.Lock()
	e, ok := n.mutexMap[name]
	if !ok {
		e = &rwMutexEntry{}
		n.mutexMap[name] = e
	}
	e.num++
	n.mutex.Unlock()
	e.mutex.Lock()
}

// Unlock releases the write lock for name. If no write lock exists for name,
// the function panics.
func (n *NamedRWMutex) Unlock(name string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	e, ok := n.mutexMap[name]
	if !ok {
		panic("no named mutex acquired: " + name)
	}
	e.mutex.Unlock()
	e.num--
	if e.num == 0 {
		delete(n.mutexMap, name)
	}
}

// RLock acquires a read lock for name. It blocks until the lock has been
// acquired.
func (n *NamedRWMutex) RLock(name string) {
	n.mutex.Lock()
	e, ok := n.mutexMap[name]
	if !ok {
		e = &rwMutexEntry{}
		n.mutexMap[name] = e
	}
	e.num++
	n.mutex.Unlock()
	e.mutex.RLock()
}

// RUnlock releases the read lock for name. If no write lock exists for name,
// the function panics.
func (n *NamedRWMutex) RUnlock(name string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	e, ok := n.mutexMap[name]
	if !ok {
		panic("no named mutex acquired: " + name)
	}
	e.mutex.RUnlock()
	e.num--
	if e.num == 0 {
		delete(n.mutexMap, name)
	}
}
