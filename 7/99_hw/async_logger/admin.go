package main

import (
	"log"
)

//type Notifier interface {
//	AddSubscriber(sub Admin_LoggingServer)
//	NewNotify(event *Event)
//}

type Admin struct {
	Broadcaster *Broadcast
	EventChan   chan *Event
	UnimplementedAdminServer
}

func NewAdmin() *Admin {
	const OP = "NewAdmin"
	log.Print(OP)

	eventChan := make(chan *Event)
	admin := &Admin{
		EventChan: eventChan,
	}
	admin.Broadcaster = NewBroadcast(eventChan)

	return admin
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	const OP = "Admin.Logging"
	log.Print(OP)

	eventChan := a.Broadcaster.Subscribe()
	defer a.Broadcaster.Unsubscribe(eventChan)

	for {
		select {
		case event := <-eventChan:
			log.Println(OP+": "+"прочитано: ", event)
			if err := server.Send(event); err != nil {
				log.Println(OP + ": ошибка отправки: " + err.Error())
				return err
			}
		case <-server.Context().Done():
			log.Println(OP + ": close connection")
			return server.Context().Err()
		}

	}
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	const OP = "Admin.Statistics"
	log.Print(OP)

	return nil
}
