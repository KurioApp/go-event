/*
Package eventd provide the event-driven API.
	bus := NewBus()
	eventd.PublishEvent(UserCreated{
		ID: "foo",
		Name: "Foo",
	})

Subscribe to capture the published events.
	bus.Subscribe(HandlerFunc(func(e Event) {
		// we got event here
	}))

To use the API on the Domain-Driven Design concept, use lifecycle.Lifecycle.
	lc := lifecycle.New(handler)
	lc.End(err)
*/
package eventd
