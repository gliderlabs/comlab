package events

import "sync"

var defaultEmitter = &Emitter{}

type Event interface {
	EventName() string
}

type Listener struct {
	EventName string
	Once      bool
	Filter    func(Event) bool
	Handler   func(Event)
}

type Emitter struct {
	sync.Mutex
	listeners []*Listener
}

func Listen(listener *Listener) { defaultEmitter.Listen(listener) }
func (e *Emitter) Listen(listener *Listener) {
	if listener.Handler == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	e.listeners = append(e.listeners, listener)
}

func Unlisten(listener *Listener) { defaultEmitter.Unlisten(listener) }
func (e *Emitter) Unlisten(listener *Listener) {
	e.Lock()
	defer e.Unlock()
	var i int
	var l *Listener
	var found bool
	for i, l = range e.listeners {
		if l == listener {
			found = true
			break
		}
	}
	if found {
		e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
	}
}

func Emit(event Event) { defaultEmitter.Emit(event) }
func (e *Emitter) Emit(event Event) {
	e.Lock()
	defer e.Unlock()
	eventName := event.EventName()
	var keep []*Listener
	for _, listener := range e.listeners {
		if listener.EventName == eventName || listener.EventName == "" {
			if listener.Filter != nil {
				if ok := listener.Filter(event); !ok {
					keep = append(keep, listener)
					continue
				}
			}
			listener.Handler(event)
			if !listener.Once {
				keep = append(keep, listener)
			}
		} else {
			keep = append(keep, listener)
		}
	}
	e.listeners = keep
}
