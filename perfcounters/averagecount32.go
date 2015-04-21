package perfcounters

import (
	"fmt"
	"sync"
)

/*

AverageCount32

An average counter that shows how many items are processed, on average, during an operation. Counters of this type display a ratio of the items processed to the number of operations
completed. The ratio is calculated by comparing the number of items processed during the last interval to the number of operations completed during the last interval.
Formula: (N 1 -N 0)/(B 1 -B 0), where N 1 and N 0 are performance counter readings, and the B 1 and B 0 are their corresponding AverageBase values. Thus, the numerator represents the
numbers of items processed during the sample interval, and the denominator represents the number of operations completed during the sample interval.

Counters of this type include PhysicalDisk\ Avg. Disk Bytes/Transfer.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type AverageCount32 struct {
	lastCount int32
	lastBase  int32
	count     int32
	base      int32
	mu        sync.Mutex
}

func NewAverageCount32() *AverageCount32 {
	return &AverageCount32{}
}

func (self *AverageCount32) Increment() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.count += 1
	self.base += 1
}

func (self *AverageCount32) Add(value int32) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.count += value
	self.base += 1
}

func (self *AverageCount32) CalculatedValue() float32 {
	self.mu.Lock()
	defer self.mu.Unlock()

	count := self.count
	base := self.base

	lastCount := self.lastCount
	lastBase := self.lastBase

	if base == 0 {
		return 0
	}

	if base-lastBase == 0 {
		return 0
	}

	calculatedValue := float32((count - lastCount) / (base - lastBase))

	self.lastCount = count
	self.lastBase = base

	return calculatedValue
}

func (self *AverageCount32) String() string {
	return fmt.Sprintf("%.3f", self.CalculatedValue())
}

/*

func main() {

    counter := NewAverageCount32()

    fmt.Println(counter.String()) // should display 0.00

    counter.Increment() // one operation, 1 item

    fmt.Println(counter.String()) // should display 1.00

    counter.Add(9) // one operation, 9 item

    fmt.Println(counter.String()) // should display 9.00
    fmt.Println(counter.String()) // should display 0.00

    counter.Add(4) // one operation, 4 items
    counter.Add(8) // one operation, 8 items

    fmt.Println(counter.String()) // should display 6.00
}

*/
