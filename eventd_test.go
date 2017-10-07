package eventd_test

import (
	"context"
	"testing"

	"github.com/KurioApp/go-eventd"
)

func TestPublisherFromContext(t *testing.T) {
	bus := eventd.NewBus()
	ctx := eventd.ContextWithBus(context.Background(), bus)

	_, ok := eventd.PublisherFromContext(ctx)
	if !ok {
		t.Error("no publisher in context")
	}
}

func TestBus_Publish(t *testing.T) {
	bus := eventd.NewBus()

	var events []eventd.Event
	bus.Subscribe(eventd.HandlerFunc(func(evt eventd.Event) {
		events = append(events, evt)
	}))

	msg := "Hello World!"
	bus.Publish("Greet", msg)
	if got, want := len(events), 1; got != want {
		t.Fatal("got:", got, "want:", want)
	}

	if got, want := events[0].Name(), "Greet"; got != want {
		t.Error("got:", got, "want:", want)
	}

	if got, want := events[0].Body(), msg; got != want {
		t.Error("got:", got, "want:", want)
	}
}

func TestPublishBody(t *testing.T) {
	bus := eventd.NewBus()

	var events []eventd.Event
	bus.Subscribe(eventd.HandlerFunc(func(evt eventd.Event) {
		events = append(events, evt)
	}))

	msg := "Hello World!"
	eventd.PublishEvent(bus, msg)
	if got, want := len(events), 1; got != want {
		t.Fatal("got:", got, "want:", want)
	}

	if got, want := events[0].Name(), "string"; got != want {
		t.Error("got:", got, "want:", want)
	}

	if got, want := events[0].Body(), msg; got != want {
		t.Error("got:", got, "want:", want)
	}
}
