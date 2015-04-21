package perfcounters

import (
	"fmt"
	"sync"
	"time"
)

/*

AverageTimer32

An average counter that measures the time it takes, on average, to complete a process or operation. Counters of this type display a ratio of the total elapsed time of the sample interval to the
number of processes or operations completed during that time. This counter type measures time in ticks of the system clock.
Formula: ((N 1 -N 0)/F)/(B 1 -B 0), where N 1 and N 0 are performance counter readings, B 1 and B 0 are their corresponding AverageBase values, and F is the number of ticks per second.
The value of F is factored into the equation so that the result can be displayed in seconds. Thus, the numerator represents the numbers of ticks counted during the last sample interval, F
represents the frequency of the ticks, and the denominator represents the number of operations completed during the last sample interval.

Counters of this type include PhysicalDisk\ Avg. Disk sec/Transfer.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type AverageTimer32 struct {
	lastTime    time.Time
	lastBase    int32
	currentTime time.Time
	currentBase int32
	mu          sync.Mutex
}

func NewAverageTimer32() *AverageTimer32 {

	now := time.Now()

	return &AverageTimer32{
		lastTime:    now,
		lastBase:    0,
		currentTime: now,
		currentBase: 0,
	}
}

func (self *AverageTimer32) Add(duration time.Duration) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.currentTime = self.currentTime.Add(duration)
	self.currentBase += 1
}

func (self *AverageTimer32) CalculatedValue() float64 {
	self.mu.Lock()
	defer self.mu.Unlock()

	lastTime := self.lastTime
	currentTime := self.currentTime
	lastBase := self.lastBase
	currentBase := self.currentBase

	if currentBase == 0 || currentBase-lastBase == 0 {
		return 0
	}

	calculatedValue := float64((currentTime.Sub(lastTime).Nanoseconds() / 1e6) / int64(currentBase-lastBase))

	self.lastTime = currentTime
	self.lastBase = currentBase

	return calculatedValue
}

func (self *AverageTimer32) String() string {
	return fmt.Sprintf("%.3f", self.CalculatedValue())
}
