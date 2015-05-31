package maxrate

import (
	"container/list"
	"time"
)

// event is a struct representing an event in the max rate calculator
// that transfers a particular amount at a particular time.
type event struct {
	transferred float64   // Amount transferred in this event.
	time        time.Time // The time of this event.
}

// MaxRate keeps track of a series of transfer events that took place over a
// past interval of time and waits Times are given in minutes.
type MaxRate struct {
	maxRate     float64    // The maximum rate.
	interval    float64    // The duration over which the rate is calculated.
	transferred float64    // Amount transferred over the given duration.
	list        *list.List // A list of events.
}

// removeExpired removes expired events from the max rate calculator.
func (m *MaxRate) removeExpired() {
	now := time.Now()
	expireTime := now.Add(time.Duration(-m.interval * float64(time.Minute)))

	//var remove *event = nil
	elem := m.list.Front()
	for elem != nil {
		e, _ := elem.Value.(*event)
		last := elem
		elem = elem.Next()

		if e.time.After(expireTime) {
			return
		}

		var transferred float64
		interval := float64(now.Sub(e.time)) / float64(time.Minute)
		if e.transferred/interval > m.maxRate {
			transferred = m.maxRate * interval
			m.list.PushBack(&event{transferred: e.transferred - transferred, time: now})
		} else {
			transferred = e.transferred
		}
		m.transferred -= transferred

		m.list.Remove(last)
	}
}

// WaitTime returns the duration of time to wait before transfering an amount
// of the given size.
func (m *MaxRate) WaitTime(size float64) time.Duration {
	m.removeExpired()
	maxTransfer := m.maxRate * m.interval

	if m.transferred+size < maxTransfer {
		return 0
	}

	var data float64
	dataToFill := m.transferred + size - maxTransfer

	if dataToFill > m.transferred {
		data = m.transferred
	} else {
		data = dataToFill
	}

	return time.Duration(float64(time.Minute) * data / m.maxRate)
}

// Transfer adds an event transfering amount x. The function
// waits to ensure that the maximum rate is not violated. 
func (m *MaxRate) Transfer(size float64) {
	wait := m.WaitTime(size)
	if wait > 0 {
		time.Sleep(wait)
	}
	m.transferred += size
	m.list.PushBack(&event{transferred: size, time: time.Now()})
}

// AverageRate returns the average rate of transfer over the duration interval.
func (m *MaxRate) AverageRate() float64 {
	m.removeExpired()
	return m.transferred / m.interval
}

// New returns a new MaxRate.
func New(rate float64, interval float64) *MaxRate {
	return &MaxRate{
		maxRate:  rate,
		interval: interval,
		list:     list.New(),
	}
}
