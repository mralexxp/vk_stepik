package main

import (
	"context"
	"time"
)

type Admin struct {
	Broadcaster *Broadcast
	EventChan   chan *Event

	UnimplementedAdminServer
}

func NewAdmin(ctx context.Context) *Admin {
	const OP = "NewAdmin"

	eventChan := make(chan *Event)

	admin := &Admin{EventChan: eventChan}

	admin.Broadcaster = NewBroadcast(ctx, eventChan)

	return admin
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	const OP = "Admin.Logging"

	eventChan := a.Broadcaster.Subscribe()
	defer a.Broadcaster.Unsubscribe(eventChan)

	for {
		select {
		case event := <-eventChan:
			if err := server.Send(event); err != nil {
				return err
			}
		case <-server.Context().Done():
			return server.Context().Err()
		}

	}
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	const OP = "Admin.Statistics"

	eventChan := a.Broadcaster.Subscribe()
	defer a.Broadcaster.Unsubscribe(eventChan)

	stat := Stat{
		ByMethod:   make(map[string]uint64),
		ByConsumer: make(map[string]uint64),
	}

	ticker := time.NewTicker(time.Duration(interval.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event := <-eventChan:
			stat.ByMethod[event.Method]++
			stat.ByConsumer[event.Consumer]++
		case <-ticker.C:
			err := server.Send(&stat)
			if err != nil {
				return err
			}
			stat.ByMethod = make(map[string]uint64)
			stat.ByConsumer = make(map[string]uint64)
		case <-server.Context().Done():
			return server.Context().Err()
		}
	}

}
