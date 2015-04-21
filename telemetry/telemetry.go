package telemetry

import (
	"github.com/israelchen/gomon/util"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type Handler interface {
	Started(telemetry *Telemetry)
	Ended(telemetry *Telemetry)
}

type Telemetry struct {
	context.Context
	name      string
	startTime *time.Time
	endTime   *time.Time
	err       error
	result    interface{}
	data      map[interface{}]interface{}
	mu        sync.RWMutex
	children  []*Telemetry
	handlers  []Handler
	closed    bool
}

var telemetryKey int = 0

func NewTelemetry(parent context.Context, name string, handlers ...Handler) (t *Telemetry) {

	util.Require(parent != nil, "telemetry: parent cannot be nil.")
	util.Require(len(name) > 0, "telemetry: name cannot be empty.")

	startTime := time.Now()

	t = &Telemetry{
		Context:   parent,
		name:      name,
		startTime: &startTime,
		endTime:   nil,
		handlers:  handlers,
		data:      make(map[interface{}]interface{}),
		children:  nil,
	}

	// search up the context chain for a parent telemetry
	parentTelemetry := parent.Value(telemetryKey)

	if parentTelemetry != nil {
		// attach ourselves to parent telemetry
		parentTelemetry.(*Telemetry).attach(t)
	}

	for _, handler := range t.handlers {
		// invoke start handler method
		handler.Started(t)
	}

	return t
}

func FromContext(ctx context.Context) (*Telemetry, bool) {

	util.Require(ctx != nil, "telemetry: ctx cannot be nil.")

	telemetry, ok := ctx.Value(telemetryKey).(*Telemetry)
	return telemetry, ok
}

func (self *Telemetry) attach(child *Telemetry) {

	self.mu.Lock()
	defer self.mu.Unlock()

	// attach child to telemetry
	self.children = append(self.children, child)
}

func (self *Telemetry) Close() {

	endTime := time.Now()
	self.endTime = &endTime

	self.mu.Lock()
	defer self.mu.Unlock()

	if self.closed {
		return
	}

	for _, child := range self.children {
		child.Close()
	}

	for _, handler := range self.handlers {
		handler.Ended(self)
	}

	self.closed = true
}

func (self *Telemetry) RecordValue(key interface{}, value interface{}) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.data[key] = value
}

func (self *Telemetry) Value(key interface{}) interface{} {

	if key == telemetryKey {
		return self
	}

	self.mu.RLock()
	defer self.mu.RUnlock()

	value := self.data[key]

	if value != nil {
		return value
	}

	return self.Context.Value(key)
}

func (self *Telemetry) Keys() []interface{} {
	self.mu.RLock()
	defer self.mu.RUnlock()

	keys := make([]interface{}, len(self.data))

	for k := range self.data {
		keys = append(keys, k)
	}

	return keys
}

func (self *Telemetry) Name() string {
	return self.name
}

func (self *Telemetry) Result() interface{} {
	return self.result
}

func (self *Telemetry) SetResult(result interface{}) {
	self.result = result
}

func (self *Telemetry) Error() error {
	return self.err
}

func (self *Telemetry) SetError(err error) {
	self.err = err
}

func (self *Telemetry) StartTime() *time.Time {
	return self.startTime
}

func (self *Telemetry) EndTime() *time.Time {
	return self.endTime
}

func (self *Telemetry) Children() []*Telemetry {
	return self.children
}
