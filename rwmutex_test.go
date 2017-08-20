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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNamedRWMutex(t *testing.T) {
	t.Run("two named locks", func(t *testing.T) {
		n := NewNamedRWMutex()
		n.Lock("a")
		n.Lock("b")
		assert.Equal(t, 2, len(n.mutexMap), "unexpected number of locks")
		assert.NotNil(t, n.mutexMap["a"], "lock 'a' not acquired")
		assert.NotNil(t, n.mutexMap["b"], "lock 'b' not acquired")
		n.Unlock("a")
		n.Unlock("b")
		assert.Equal(t, 0, len(n.mutexMap), "locks not released")
		assert.Nil(t, n.mutexMap["a"], "lock 'a' not released")
		assert.Nil(t, n.mutexMap["b"], "lock 'b' not released")
	})

	t.Run("lock", func(t *testing.T) {
		n := NewNamedRWMutex()
		x := 0
		n.Lock("a")
		go func() {
			time.Sleep(time.Millisecond * 100)
			x++
			n.Unlock("a")
		}()
		n.Lock("a")
		x *= 10
		n.Unlock("a")
		assert.Equal(t, 10, x, "incorrect result")
	})

	t.Run("two named rlocks", func(t *testing.T) {
		n := NewNamedRWMutex()
		n.RLock("a")
		n.RLock("b")
		assert.Equal(t, 2, len(n.mutexMap), "unexpected number of locks")
		assert.NotNil(t, n.mutexMap["a"], "lock 'a' not acquired")
		assert.NotNil(t, n.mutexMap["b"], "lock 'b' not acquired")
		n.RUnlock("a")
		n.RUnlock("b")
		assert.Equal(t, 0, len(n.mutexMap), "locks not released")
		assert.Nil(t, n.mutexMap["a"], "lock 'a' not released")
		assert.Nil(t, n.mutexMap["b"], "lock 'b' not released")
	})

	t.Run("rlock after lock", func(t *testing.T) {
		n := NewNamedRWMutex()
		x := 0
		n.Lock("a")
		go func() {
			time.Sleep(time.Millisecond * 100)
			x++
			n.Unlock("a")
		}()
		n.RLock("a")
		assert.Equal(t, 1, x, "incorrect result")
		n.RUnlock("a")
	})

	t.Run("lock after rlock", func(t *testing.T) {
		n := NewNamedRWMutex()
		x := 0
		n.RLock("a")
		go func() {
			time.Sleep(time.Millisecond * 100)
			assert.Equal(t, 0, x, "incorrect result")
			n.RUnlock("a")
		}()
		n.Lock("a")
		x++
		n.Unlock("a")
	})

	t.Run("multiple rlocks", func(t *testing.T) {
		n := NewNamedRWMutex()
		x := 0
		go func() {
			n.RLock("a")
			time.Sleep(time.Millisecond * 100)
			assert.Equal(t, 0, x, "incorrect result")
			n.RUnlock("a")
		}()
		go func() {
			n.RLock("a")
			time.Sleep(time.Millisecond * 100)
			assert.Equal(t, 0, x, "incorrect result")
			n.RUnlock("a")
		}()
		go func() {
			time.Sleep(time.Millisecond * 50)
			n.Lock("a")
			x++
			n.Unlock("a")
		}()
	})
}
