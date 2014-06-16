package rpc

import (
	"log"
	"sync"
)

type EventHandler func(v interface{}, err error) error

type Event struct {
	sync.Mutex
	handlers map[int]EventHandler
	seq      int
}

func (e *Event) Attach(handler EventHandler) int {
	e.Lock()
	defer e.Unlock()
	e.seq++
	e.handlers[e.seq] = handler
	return e.seq
}

func (e *Event) Detach(handle int) {
	e.Lock()
	defer e.Unlock()
	delete(e.handlers, handle)
}

func (e *Event) publish(v interface{}) {
	e.Lock()
	defer e.Unlock()
	for k, handler := range e.handlers {
		if handler != nil {
			err := handler(v, nil)
			if err != nil {
				log.Println("Event.publish error:", err)
				delete(e.handlers, k)
			}
		}
	}
}

type EventPublisher struct {
	event Event
}

func (p *EventPublisher) Event() *Event {
	if p.event.handlers == nil {
		p.event.handlers = make(map[int]EventHandler)
		p.event.seq = 1024
	}
	return &p.event
}

func (p *EventPublisher) Publish(v interface{}) {
	p.event.publish(v)
}
