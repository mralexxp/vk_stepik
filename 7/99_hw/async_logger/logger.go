package main

import (
	"fmt"
	"log"
)

type Logger struct {
	Subscribers []Admin_LoggingServer
}

func NewLogger() *Logger {
	const OP = "NewLogger"
	log.Print(OP)

	logger := &Logger{}
	logger.Subscribers = make([]Admin_LoggingServer, 0)
	return logger
}

func (l *Logger) AddSubscriber(sub Admin_LoggingServer) {
	const OP = "Logger.AddSubscriber"
	log.Print(OP)

	l.Subscribers = append(l.Subscribers, sub)
}

func (l *Logger) NewNotify(event *Event) {
	const OP = "Logger.NewNotify"
	log.Print(OP)

	for _, sub := range l.Subscribers {
		err := sub.Send(event)
		if err != nil {
			fmt.Println("Error sending event to subscriber")
		}
	}
}
