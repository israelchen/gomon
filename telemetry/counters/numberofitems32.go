package perfcounters

import (
	"strconv"
	"sync/atomic"
)

/*

NumberOfItems32

An instantaneous counter that shows the most recently observed value. Used, for example, to maintain a simple count of items or operations.
Formula: None. Does not display an average, but shows the raw data as it is collected.
Counters of this type include Memory\Available Bytes.

[[source: https://msdn.microsoft.com/en-us/library/system.diagnostics.performancecountertype(v=vs.90).aspx]]

*/

type NumberOfItems32 struct {
	count int32
}

func NewNumberOfItems32() *NumberOfItems32 {
	return &NumberOfItems32{
		count: 0,
	}
}

func (self *NumberOfItems32) Increment() {
	atomic.AddInt32(&(self.count), 1)
}

func (self *NumberOfItems32) Add(count int32) {
	atomic.AddInt32(&(self.count), count)
}

func (self *NumberOfItems32) String() string {
	return strconv.Itoa(int(self.count))
}
