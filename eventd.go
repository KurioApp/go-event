package eventd

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type contextKey int

const (
	keyBus contextKey = iota
)

// Event is the domain event.
type Event struct {
	name        string
	body        interface{}
	occuredTime time.Time
}

// Name of the event.
func (e Event) Name() string {
	return e.name
}

// Body of the event.
func (e Event) Body() interface{} {
	return e.body
}

// OccuredTime of the event.
func (e Event) OccuredTime() time.Time {
	return e.occuredTime
}

// NewEvent construtcs new event.
func NewEvent(name string, body interface{}) Event {
	return Event{
		name:        name,
		body:        body,
		occuredTime: time.Now(),
	}
}

// Publisher publish event.
type Publisher interface {
	Publish(name string, body interface{})
}

// Handler handles the event.
type Handler interface {
	Handle(Event)
}

// HandlerFunc is adapter for Handler.
type HandlerFunc func(Event)

// Handle invokes f(e)
func (f HandlerFunc) Handle(e Event) {
	f(e)
}

// HandleOnly specific type of event.
func HandleOnly(h Handler, names ...string) Handler {
	return HandlerFunc(func(e Event) {
		for _, name := range names {
			if e.Name() == name {
				h.Handle(e)
				break
			}
		}
	})
}

// ContextWithBus constructs new context with bus.
func ContextWithBus(ctx context.Context, bus *Bus) context.Context {
	return context.WithValue(ctx, keyBus, bus)
}

// PublisherFromContext return the publisher.
func PublisherFromContext(ctx context.Context) (Publisher, bool) {
	pub, ok := ctx.Value(keyBus).(Publisher)
	return pub, ok
}

// Bus of the event.
type Bus struct {
	sync.RWMutex
	handlers []Handler
}

// NewBus constructs bus.
func NewBus() *Bus {
	return &Bus{}
}

// Subscribe to events.
// It will return true if successfully subscribed, means the handler haven't subscribed before.
func (b *Bus) Subscribe(h Handler) bool {
	b.Lock()
	defer b.Unlock()

	for _, v := range b.handlers {
		if v == h {
			return false
		}
	}

	b.handlers = append(b.handlers, h)
	return true
}

// Unsubscribe handler from the bus.
// It will return false if the handler not subscribed.
func (b *Bus) Unsubscribe(h Handler) bool {
	b.Lock()
	defer b.Unlock()

	for i, v := range b.handlers {
		if v == h {
			b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
			return true
		}
	}

	return false
}

// Publish an event.
func (b *Bus) Publish(name string, body interface{}) {
	e := NewEvent(name, body)

	b.RLock()
	defer b.RUnlock()
	for _, h := range b.handlers {
		h.Handle(e)
	}
}

// PublishEvent publish the event, event name will be defered from the type name.
func PublishEvent(p Publisher, body interface{}) {
	name := reflect.TypeOf(body).Name()
	p.Publish(name, body)
}
