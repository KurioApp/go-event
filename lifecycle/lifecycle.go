package lifecycle

import (
	"context"
	"sync"

	"github.com/KurioApp/go-eventd"
)

// New constructs new life cycle.
// Create life cycle required to end it by invokes End() method.
func New(h EventHandler) *LifeCycle {
	bus := eventd.NewBus()
	eventCtx := eventd.ContextWithBus(context.Background(), bus)

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
	events []eventd.Event

	ctx          context.Context
	bus          *eventd.Bus
	eventHandler EventHandler
}

func (lc *LifeCycle) listenForEvents() {
	lc.bus.Subscribe(eventd.HandlerFunc(func(e eventd.Event) {
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
	HandleEvent([]eventd.Event)
}

// EventHandlerFunc is the adapter of EventHandler.
type EventHandlerFunc func([]eventd.Event)

// HandleEvent calls f(s)
func (f EventHandlerFunc) HandleEvent(events []eventd.Event) {
	f(events)
}

type nopEventHandler struct {
}

func (h nopEventHandler) HandleEvent([]eventd.Event) {}

// NopEventHandler implements no-operation EventHandler.
func NopEventHandler() EventHandler {
	return nopEventHandler{}
}
