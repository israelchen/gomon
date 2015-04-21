package telemetry

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

type TestHandler struct {
	startedHandler func(t *Telemetry)
	endedHandler   func(t *Telemetry)
}

func (handler *TestHandler) Started(t *Telemetry) {

	if handler.startedHandler != nil {
		handler.startedHandler(t)
	}
}

func (handler *TestHandler) Ended(t *Telemetry) {

	if handler.endedHandler != nil {
		handler.endedHandler(t)
	}
}

func TestNewTelemetry(t *testing.T) {

	ctx := NewTelemetry(context.Background(), "test.telemetry")
	defer ctx.Close()

	if ctx.Name() != "test.telemetry" {
		t.Error("Name is different than expected.")
	}

	if c, ok := FromContext(ctx); ok {

		if c != ctx {
			t.Error("FromContext result is different than expected.")
		}
	} else {
		t.Error("FromContext was not returned.")
	}
}

func TestElapsedWorksCorrectly(t *testing.T) {

	ctx := NewTelemetry(context.Background(), "test.telemetry")

	time.Sleep(50 * time.Millisecond)

	// manually close telemetry
	ctx.Close()

	if ctx.EndTime().Sub(*ctx.StartTime()).Nanoseconds()/1e6 < 50 {
		t.Error("Elapsed is less than expected.")
	}
}

func TestNestedAttachesToBase(t *testing.T) {

	base := NewTelemetry(context.Background(), "test.telemetry.base")
	nested := NewTelemetry(base, "test.telemetry.nested")

	if len(base.children) < 1 {
		t.Fatal("nested did not attach itself to base.")
	}

	if base.children[0] != nested {
		t.Fatal("nested did not attach itself to base.")
	}
}

func TestCloseNestedDoesNotCloseBase(t *testing.T) {

	var baseCloseCalled bool

	baseCloseHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			baseCloseCalled = true
		},
	}

	base := NewTelemetry(context.Background(), "test.telemetry.base", baseCloseHandler)

	var nestedCloseCalled bool

	nestedCloseHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			nestedCloseCalled = true
		},
	}

	nested := NewTelemetry(base, "test.telemetry.nested", nestedCloseHandler)

	// manually close nested. Only nested handlers should be invoked, base remains intact.
	nested.Close()

	if nestedCloseCalled == false {
		t.Fatal("nested close was not called.")
	}

	if baseCloseCalled == true {
		t.Fatal("base close was called.")
	}
}

func TestCloseBaseAlsoClosesNested(t *testing.T) {

	var baseCloseCalled bool

	baseCloseHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			baseCloseCalled = true
		},
	}

	base := NewTelemetry(context.Background(), "test.telemetry.base", baseCloseHandler)

	var nestedCloseCalled bool

	nestedCloseHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			nestedCloseCalled = true
		},
	}

	nested := NewTelemetry(base, "test.telemetry.nested", nestedCloseHandler)

	nested = nested

	// manually close base. Both base and nested's handlers should be invoked.
	base.Close()

	if baseCloseCalled == false {
		t.Fatal("base close was called.")
	}

	if nestedCloseCalled == false {
		t.Fatal("nested close was not called.")
	}
}

func TestTelemetryCanOnlyBeClosedOnce(t *testing.T) {

	var wasCloseCalled bool

	closeHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			wasCloseCalled = true
		},
	}

	ctx := NewTelemetry(context.Background(), "test.telemetry.base", closeHandler)

	// manually close telemetry.
	ctx.Close()

	if wasCloseCalled == false {
		t.Fatal("close was not called.")
	}

	wasCloseCalled = false

	// close again
	ctx.Close()

	if wasCloseCalled == true {
		t.Fatal("close was called again.")
	}
}

func TestClosureCapturesBooleanByReference(t *testing.T) {

	var wasCalled bool

	testHandler := &TestHandler{

		endedHandler: func(tel *Telemetry) {
			wasCalled = true
		},
	}

	testHandler.endedHandler(nil)

	if wasCalled == false {
		t.Fatal("was not set the first time.")
	}

	wasCalled = false

	testHandler.endedHandler(nil)

	if wasCalled == false {
		t.Fatal("was not set the second time.")
	}
}
