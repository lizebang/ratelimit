/*
 * Revision History:
 *     Initial: 2018/07/12        Li Zebang
 */

package ratelimit

import (
	"sync"
	"time"
)

// Limiter -
type Limiter interface {
	Take(n int) int
	ResetInterval(interval int64)
	ResetQuantum(quantum int64)
}

type limiter struct {
	mux      *sync.Mutex
	current  int64
	interval int64
	capacity int64
	quantum  int64
}

// NewLimiter -
func NewLimiter(interval, quantum, initial int64) Limiter {
	return &limiter{
		mux:      &sync.Mutex{},
		current:  time.Now().Unix(),
		interval: interval,
		capacity: initial,
		quantum:  quantum,
	}
}

func (l *limiter) adjust(t time.Time) {
	l.capacity += (t.Unix() - l.current) / l.interval * (l.quantum / 3600)
	if l.capacity > l.quantum {
		l.capacity = l.quantum
	}
	l.current = t.Unix()
}

// Take -
func (l *limiter) Take(n int) int {
	l.mux.Lock()
	defer l.mux.Unlock()
	if n <= 0 {
		return 0
	}
	l.adjust(time.Now())
	if l.capacity < int64(n) {
		n = int(l.capacity)
		l.capacity = 0
		return n
	}
	l.capacity -= int64(n)
	return n
}

// ResetInterval -
func (l *limiter) ResetInterval(interval int64) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.interval = interval
}

// ResetQuantum -
func (l *limiter) ResetQuantum(quantum int64) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.quantum = quantum
}
