package eventd

import (
	"context"
	"sync"
)

// NewLifeCycle constructs new life cycle.
// Create life cycle required to end it by invokes End() method.
func NewLifeCycle(h EventHandler) *LifeCycle {
	bus := NewBus()
	eventCtx := ContextWithBus(context.Background(), bus)

	lc := &LifeCycle{
		ctx:          eventCtx,
		bus:          bus,
		eventHandler: h,
	}
	lc.listenForEvents()

	return lc
}

// LifeCycle of the app.
type LifeCycle struct {
	sync.RWMutex
	events []Event

	ctx          context.Context
	bus          *Bus
	eventHandler EventHandler
}

func (lc *LifeCycle) listenForEvents() {
	lc.bus.Subscribe(HandlerFunc(func(e Event) {
		lc.Lock()
		defer lc.Unlock()

		lc.events = append(lc.events, e)
	}))
}

// Context of the lifecycle.
func (lc *LifeCycle) Context() context.Context {
	return lc.ctx
}

// End the lifecycle.
func (lc *LifeCycle) End(err error) {
	if err == nil {
		lc.handleEvents()
	}
}

func (lc *LifeCycle) handleEvents() {
	if lc.eventHandler != nil {
		lc.RLock()
		defer lc.RUnlock()

		lc.eventHandler.HandleEvent(lc.events)
	}
}

// EventHandler handles the events.
type EventHandler interface {
	HandleEvent([]Event)
}

// EventHandlerFunc is the adapter of EventHandler.
type EventHandlerFunc func([]Event)

// HandleEvent calls f(s)
func (f EventHandlerFunc) HandleEvent(events []Event) {
	f(events)
}

type nopEventHandler struct {
}

func (h nopEventHandler) HandleEvent([]Event) {}

// NopEventHandler implements no-operation EventHandler.
func NopEventHandler() EventHandler {
	return nopEventHandler{}
}
