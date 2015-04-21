package telemetry

import (
	"expvar"
	"github.com/israelchen/gomon/perfcounters"
	"github.com/israelchen/gomon/util"
)

type PerfHandler struct {
	publishedMap    *expvar.Map
	totalCalls      *perfcounters.NumberOfItems32
	successfulCalls *perfcounters.NumberOfItems32
	failedCalls     *perfcounters.NumberOfItems32
	callsPerSec     *perfcounters.CountPerTimeInterval32
	callLatency     *perfcounters.AverageTimer32
}

func NewPerfHandler(telemetryName string) *PerfHandler {
	util.Require(len(telemetryName) > 0, "telemetry: telemetryName cannot be empty.")

	m := expvar.NewMap(telemetryName)

	handler := &PerfHandler{
		publishedMap:    m,
		totalCalls:      perfcounters.NewNumberOfItems32(),
		successfulCalls: perfcounters.NewNumberOfItems32(),
		failedCalls:     perfcounters.NewNumberOfItems32(),
		callsPerSec:     perfcounters.NewCountPerTimeInterval32(),
		callLatency:     perfcounters.NewAverageTimer32(),
	}

	m.Set("totalCalls", handler.totalCalls)
	m.Set("successfulCalls", handler.successfulCalls)
	m.Set("failedCalls", handler.failedCalls)
	m.Set("callsPerSec", handler.callsPerSec)
	m.Set("callLatency", handler.callLatency)

	return handler
}

func (self *PerfHandler) Started(t *Telemetry) {
	util.Require(t != nil, "telemetry: telemetry cannot be nil.")

	self.totalCalls.Increment()
	self.callsPerSec.Increment()
}

func (self *PerfHandler) Ended(t *Telemetry) {
	util.Require(t != nil, "telemetry: telemetry cannot be nil.")
	util.Require(t.EndTime() != nil, "telemetry: endTime cannot be nil. This handler should be invoked on telemetry end operation only.")

	if t.Error() == nil {
		self.successfulCalls.Increment()
	} else {
		self.failedCalls.Increment()
	}

	elapsed := t.StartTime().Sub(*t.EndTime())
	self.callLatency.Add(elapsed)
}
