/*
Package lifecycle deals with application lifecycle.

Concept

On Domain-Driven Design there is application service concept. Application service act as Use Case or service that deals directly with the user code or Port (term in Hexagonal architecture).
Application service maintain consistency and transaction also handled here.

When the domain event raised, the changes on the aggreate hasn't really applied to the repository and there is possibility operation on the repository faild that makes the domain-event should not applied.
This os why lifecycle created. It the event should only apply on the end of the lifecycle of the app transaction.

	func DoSomething() (err error) {
		lc := lifecycle.New(handler)
		defer lc.End(err)

		// some domain logic that raise events
		// ...
	}

Implement the handler:
	handler := lifecycle.HandleFunc(func(s []eventd.Event) {
		// TODO handle events here
	})

Any error found when lc.End(err) will ignore any of the captured events.
*/
package lifecycle
