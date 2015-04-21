package perfcounters

import (
	"fmt"
	"math"
	"sync"
	"time"
)

/*

RateOfCountsPerSecond32

A difference counter that shows the average number of operations completed during each second of the sample interval. Counters of this type measure time in ticks of the system clock.
Formula: (N 1 - N 0) / ((D 1 -D 0) / F), where N 1 and N 0 are performance counter readings, D 1 and D 0 are their corresponding time readings, and F represents the number of ticks per second.
Thus, the numerator represents the number of operations performed during the last sample interval, the denominator represents the number of ticks elapsed during the last sample interval, and F is
the frequency of the ticks. The value of F is factored into the equation so that the result can be displayed in seconds.
Counters of this type include System\ File Read Operations/sec.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type RateOfCountsPerSecond32 struct {
	lastTime     *time.Time
	lastCount    int32
	currentCount int32
	mu           sync.Mutex
}

func NewRateOfCountsPerSecond32() *RateOfCountsPerSecond32 {

	return &RateOfCountsPerSecond32{
		lastTime:     nil,
		lastCount:    0,
		currentCount: 0,
	}
}

func (self *RateOfCountsPerSecond32) Increment() {
	self.Add(1)
}

func (self *RateOfCountsPerSecond32) Add(value int32) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.currentCount += value

	if self.lastTime == nil {
		lastTime := time.Now()
		self.lastTime = &lastTime
	}
}

func (self *RateOfCountsPerSecond32) CalculatedValue() float64 {
	self.mu.Lock()
	defer self.mu.Unlock()

	currentTime := time.Now()

	if self.lastTime == nil {
		self.lastTime = &currentTime
		return 0
	}

	lastTime := self.lastTime
	lastCount := self.lastCount
	currentCount := self.currentCount

	diff := currentTime.Sub(*lastTime)

	calculatedValue := float64(currentCount-lastCount) / diff.Seconds()

	if math.IsNaN(calculatedValue) || math.IsInf(calculatedValue, 1) || math.IsInf(calculatedValue, -1) {
		calculatedValue = 0.0
	}

	self.lastCount = currentCount
	self.lastTime = &currentTime

	return calculatedValue
}

func (self *RateOfCountsPerSecond32) String() string {
	return fmt.Sprintf("%.3f", self.CalculatedValue())
}

/*

func main() {

    ctr := NewRateOfCountsPerSecond32()

    fmt.Println(ctr.String()) // expected: 0

    ctr.Increment()

    time.Sleep(1000 * time.Millisecond)

    fmt.Println(ctr.String()) // expected: 1

    ctr.Increment()
    ctr.Increment()

    time.Sleep(1000 * time.Millisecond)

    fmt.Println(ctr.String()) // expected: 2

    ctr.Increment()

    time.Sleep(2000 * time.Millisecond)

    fmt.Println(ctr.String()) // expected: 0.5
}

*/
