package rpc

import (
	"sync"
)

type EventHandler func(v interface{})

type Event struct {
	sync.Mutex
	handlers map[int]EventHandler
	seq      int
}

func (e *Event) Attach(handler EventHandler) int {
	e.Lock()
	defer e.Unlock()
	seq := e.seq
	e.handlers[seq] = handler
	e.seq++
	return seq
}

func (e *Event) Detach(handle int) {
	e.Lock()
	defer e.Unlock()
	delete(e.handlers, handle)
}

func (e *Event) publish(v interface{}) {
	e.Lock()
	defer e.Unlock()
	for _, handler := range e.handlers {
		if handler != nil {
			handler(v)
		}
	}
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	if p.event.handlers == nil {
		p.event.handlers = make(map[int]EventHandler)
	}
	return &p.event
}

func (p *EventPublisher) Publish(v interface{}) {
	p.event.publish(v)
}
