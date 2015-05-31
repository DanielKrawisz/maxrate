package maxrate

import (
	"time"
	//"container/list"
)

// TstNew returns a new new MaxRate that uses the given list.
/*func TstNew(rate float64, interval float64, list *list.List) *MaxRate {
	return &MaxRate {
		maxRate: rate,
		interval: interval,
		list: list,
	}
}*/

// TstAddEvent takes a list, an amount transferred, and a time and pushes
// an event struct, which is unexported, to the list.
// It works like Transfer except without actually waiting.
func (m *MaxRate) TstTransfer(size float64, time time.Time) {
	m.transferred += size
	m.list.PushBack(&event{transferred: size, time: time})
}
