package perfcounters

import (
	"fmt"
	"sync"
	"time"
)

/*

CountPerTimeInterval32

An average counter designed to monitor the average length of a queue to a resource over time. It shows the difference between the queue lengths observed during the last two sample
intervals divided by the duration of the interval. This type of counter is typically used to track the number of items that are queued or waiting.
Formula: (N 1 - N 0) / (D 1 - D 0), where the numerator represents the number of items in the queue and the denominator represents the time elapsed during the last sample interval.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type CountPerTimeInterval32 struct {
	lastCount    int32
	lastTime     *time.Time
	currentCount int32
	mu           sync.Mutex
}

func NewCountPerTimeInterval32() *CountPerTimeInterval32 {

	return &CountPerTimeInterval32{
		lastTime:     nil,
		lastCount:    0,
		currentCount: 0,
	}
}

func (self *CountPerTimeInterval32) Increment() {
	self.Add(1)
}

func (self *CountPerTimeInterval32) Add(value int32) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.currentCount += value

	if self.lastTime == nil {
		now := time.Now()
		self.lastTime = &now
	}
}

func (self *CountPerTimeInterval32) CalculatedValue() float64 {
	self.mu.Lock()
	defer self.mu.Unlock()

	currentTime := time.Now()

	if self.lastTime == nil {
		self.lastTime = &currentTime
		return 0
	}

	lastTime := *self.lastTime
	lastCount := self.lastCount
	currentCount := self.currentCount

	calculatedValue := float64(int64(currentCount-lastCount) / (currentTime.Sub(lastTime).Nanoseconds() / 1e6))

	self.lastTime = &currentTime
	self.lastCount = currentCount

	return calculatedValue
}

func (self *CountPerTimeInterval32) String() string {
	return fmt.Sprintf("%.3f", self.CalculatedValue())
}
