package main

import (
	"log"
	"sync"
)

type Broadcast struct {
	evnt          chan *Event
	subscribers   map[chan *Event]struct{}
	subscribersMu sync.RWMutex
}

func NewBroadcast(eventChan chan *Event) *Broadcast {
	b := &Broadcast{
		evnt:        eventChan,
		subscribers: make(map[chan *Event]struct{}),
	}

	go b.broadcaster()

	return b
}

func (b *Broadcast) SendEvent(event *Event) {
	const OP = "Broadcast.SendEvent"
	log.Print(OP)

	b.evnt <- event
}

func (b *Broadcast) broadcaster() {
	const OP = "Broadcast.Broadcaster"
	log.Print(OP)

	for event := range b.evnt {
		b.subscribersMu.RLock()

		for subChan := range b.subscribers {
			select {
			case subChan <- event:
				log.Print(event.String())
			default:
				log.Print("пропущена запись: ", event, "для канала ", subChan)
			}
		}

		b.subscribersMu.RUnlock()
	}
}

func (b *Broadcast) Subscribe() chan *Event {
	const OP = "Broadcast.Subscribe"
	log.Print(OP)

	ch := make(chan *Event)

	b.subscribersMu.Lock()
	b.subscribers[ch] = struct{}{}
	b.subscribersMu.Unlock()

	return ch
}

func (b *Broadcast) Unsubscribe(ch chan *Event) {
	const OP = "Broadcast.Unsubscribe"
	log.Print(OP)

	b.subscribersMu.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.subscribersMu.Unlock()
}
