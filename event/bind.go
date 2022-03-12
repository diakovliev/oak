package event

import "sync/atomic"

// Q: Why do Bind / Unbind / etc not immediately take effect?
// A: For concurrent safety, most operations on a bus lock the bus. Triggers acquire a read lock on the bus,
//    as they iterate over internal bus components. Most logic within an event bus will happen from within
//    a Trigger call-- when an entity is destroyed by some collision, for example, all of its bindings should
//    be unregistered. If one were to call Unbind from within a

// Q: Why not trust users to call Bind / Unbind / etc with `go`, to allow the caller to decide when to use
//    concurrency?
// A: It is almost never correct to not call these functions with `go`, and it is a bad user experience for
//    the engine to deadlock unexpectedly because you forgot to begin some call with a goroutine.

// A Binding, returned from calls to Bind, references the details of a binding and where that binding is
// stored within a handler. The common use case for this structure would involve a system that wanted to
// keep track of its bindings for later remote unbinding. This structure can also be used to construct
// and unbind a known reference.
type Binding struct {
	Handler  Handler
	EventID  UnsafeEventID
	CallerID CallerID
	BindID   BindID
}

// Unbind unbinds the callback associated with this binding from it's own event handler. If this binding
// does not belong to its handler or has already been unbound, this will do nothing.
func (b Binding) Unbind() {
	b.Handler.Unbind(b)
}

// A BindID is a unique identifier for a binding within a bus.
type BindID int64

// UnsafeBind registers a callback function to be called whenever the provided event is triggered
// against this bus. The binding is concurrently bound, and therefore may not be immediately
// available to be triggered. When Reset is called on a Bus, all prior bindings are unbound. This
// call is 'unsafe' because UnsafeBindables use bare interface{} types.
func (bus *Bus) UnsafeBind(eventID UnsafeEventID, callerID CallerID, fn UnsafeBindable) Binding {
	bindID := BindID(atomic.AddInt64(bus.nextBindID, 1))
	go func() {
		bus.mutex.Lock()
		bus.getBindableList(eventID, callerID).storeBindable(fn, bindID)
		bus.mutex.Unlock()
	}()
	return Binding{
		Handler:  bus,
		EventID:  eventID,
		CallerID: callerID,
		BindID:   bindID,
	}
}

// PersistentBind acts like UnsafeBind, but cause Bind to be called with these inputs after a Bus is Reset, i.e.
// persisting the binding through bus resets. Unbinding this will not stop it from being rebound on the next
// Bus Reset-- ClearPersistentBindings will.
func (bus *Bus) PersistentBind(eventID UnsafeEventID, callerID CallerID, fn UnsafeBindable) Binding {
	binding := bus.UnsafeBind(eventID, callerID, fn)
	go func() {
		bus.mutex.Lock()
		bus.persistentBindings = append(bus.persistentBindings, persistentBinding{
			eventID:  eventID,
			callerID: callerID,
			fn:       fn,
		})
		bus.mutex.Unlock()
	}()
	return binding
}

// Unbind unregisters a binding from a bus concurrently. Once complete, triggers that would
// have previously caused the Bindable callback to execute will no longer do so.
func (bus *Bus) Unbind(loc Binding) {
	go func() {
		bus.mutex.Lock()
		bus.getBindableList(loc.EventID, loc.CallerID).remove(loc.BindID)
		bus.mutex.Unlock()
	}()
}

// A Bindable is a strongly typed callback function to be executed on Trigger. It must be paired
// with an event registered via RegisterEvent.
type Bindable[T any] func(CallerID, T) Response

func Bind[T any](b Handler, ev EventID[T], c CallerID, fn Bindable[T]) Binding {
	return b.UnsafeBind(ev.UnsafeEventID, c, func(c CallerID, f interface{}) Response {
		tf := f.(T)
		return fn(c, tf)
	})
}

// UnsafeBindable defines the underlying signature of all bindings.
type UnsafeBindable func(CallerID, interface{}) Response

func EmptyBinding(f func()) UnsafeBindable {
	return func(ci CallerID, i interface{}) Response {
		f()
		return NoResponse
	}
}
