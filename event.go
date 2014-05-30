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
	for k, handler := range e.handlers {
		if handler != nil {
			k := k
			handler := handler
			go func() {
				err := handler(v, nil)
				if err != nil {
					log.Println("Event.publish error:", err)
					e.Lock()
					delete(e.handlers, k)
					e.Unlock()
				}
			}()
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
