package perfcounters

import (
	"strconv"
	"sync/atomic"
)

/*

NumberOfItems64

An instantaneous counter that shows the most recently observed value. Used, for example, to maintain a simple count of items or operations.
Formula: None. Does not display an average, but shows the raw data as it is collected.
Counters of this type include Memory\Available Bytes.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type NumberOfItems64 struct {
	count int64
}

func NewNumberOfItems64() *NumberOfItems64 {
	return &NumberOfItems64{
		count: 0,
	}
}

func (self *NumberOfItems64) Increment() {
	atomic.AddInt64(&(self.count), 1)
}

func (self *NumberOfItems64) Add(count int64) {
	atomic.AddInt64(&(self.count), count)
}

func (self *NumberOfItems64) String() string {
	value := atomic.LoadInt64(&self.count)
	return strconv.FormatInt(value, 10)
}
