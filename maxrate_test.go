package maxrate_test

import (
	"math"
	"testing"
	"time"
	//"container/list"

	"github.com/DanielKrawisz/maxrate"
)

type MockEvent struct {
	transferred float64
	time        time.Time
}

type WaitTest struct {
	testWait float64
	expWait  float64
}

const ϵ = .00001

func closeEnough(n, m float64) bool {
	return math.Abs(n-m) < ϵ
}

// Tests cases for an empty list.
func TestMaxRate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		maxRate  float64
		interval float64
		initList []MockEvent
		expAvg   float64
		wait     []WaitTest
	}{
		{ // A test case for lists that should be empty.
			maxRate:  1,
			interval: 1,
			initList: []MockEvent{
				{
					transferred: .7,
					time:        now.Add(-2 * time.Minute),
				},
			},
			expAvg: 0,
			wait: []WaitTest{
				{
					testWait: .5,
					expWait:  0,
				},
				{
					testWait: 1.5,
					expWait:  0,
				},
			},
		},
		{
			maxRate:  2,
			interval: 5,
			initList: []MockEvent{
				{
					transferred: 3,
					time:        now.Add(-8 * time.Minute),
				},
				{
					transferred: 14,
					time:        now.Add(-6 * time.Minute),
				},
				{
					transferred: 4,
					time:        now.Add(-3 * time.Minute),
				},
			},
			expAvg: 1.2,
			wait: []WaitTest{
				{
					testWait: 1,
					expWait:  0,
				},
				{
					testWait: 6,
					expWait:  1,
				},
				{
					testWait: 16,
					expWait:  3,
				},
			},
		},
	}

	for n, test := range tests {
		if n != 1 {
			continue
		}

		max := maxrate.New(test.maxRate, test.interval)

		for _, elem := range test.initList {
			max.TstTransfer(elem.transferred, elem.time)
		}

		avg := max.AverageRate()
		if !closeEnough(avg, test.expAvg) {
			t.Errorf("Incorrect average on trial %d: expected %f, got %f.", n, test.expAvg, avg)
		}

		for m, wait := range test.wait {
			testWait := float64(max.WaitTime(wait.testWait)) / float64(time.Minute)

			if !closeEnough(testWait, wait.expWait) {
				t.Errorf("Incorrect wait time on trial %d, %d: expected %f, got %f.", n, m, wait.expWait, testWait)
			}
		}
	}
}

//
func TestTransfer(t *testing.T) {

}
