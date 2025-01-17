package main

import (
	"golang.org/x/net/context"
	"log"
	"sync"
)

type Broadcast struct {
	evnt          chan *Event
	subscribers   map[chan *Event]struct{}
	subscribersMu sync.RWMutex
}

func NewBroadcast(ctx context.Context, eventChan chan *Event) *Broadcast {
	b := &Broadcast{
		evnt:        eventChan,
		subscribers: make(map[chan *Event]struct{}),
	}

	go b.broadcaster(ctx)

	return b
}

func (b *Broadcast) SendEvent(event *Event) {
	const OP = "Broadcast.SendEvent"

	b.evnt <- event
}

func (b *Broadcast) broadcaster(ctx context.Context) {
	const OP = "Broadcast.Broadcaster"

	for event := range b.evnt {
		b.subscribersMu.RLock()

		for subChan := range b.subscribers {
			select {
			case subChan <- event:
				//log.Print(event.String())
			case <-ctx.Done():
				return
			default:
				log.Print("пропущена запись: ", event, "для канала ", subChan)
			}
		}

		b.subscribersMu.RUnlock()
	}
}

func (b *Broadcast) Subscribe() chan *Event {
	const OP = "Broadcast.Subscribe"

	ch := make(chan *Event)

	b.subscribersMu.Lock()
	b.subscribers[ch] = struct{}{}
	b.subscribersMu.Unlock()

	return ch
}

func (b *Broadcast) Unsubscribe(ch chan *Event) {
	const OP = "Broadcast.Unsubscribe"

	b.subscribersMu.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.subscribersMu.Unlock()
}

func (b *Broadcast) Stop() {
	const OP = "Broadcast.Stop"

	b.subscribersMu.Lock()
	for ch := range b.subscribers {
		close(ch)
	}

	b.subscribersMu.Unlock()

	close(b.evnt)
}
