package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"
	"log"
	"time"
)

type Notifier interface {
	AddSubscriber(sub Admin_LoggingServer)
	NewNotify(event *Event)
}

type Admin struct {
	LogNotifier *Logger

	UnimplementedAdminServer
}

func NewAdmin() *Admin {
	const OP = "NewAdmin"
	log.Print(OP)

	return &Admin{
		LogNotifier: NewLogger(),
	}
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	const OP = "Admin.Logging"
	log.Print(OP)

	a.LogNotifier.AddSubscriber(server)

	return nil
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {
	const OP = "Admin.Statistics"
	log.Print(OP)

	return nil
}

func (a *Admin) TapLogger(ctx context.Context, info *tap.Info) (context.Context, error) {
	const OP = "Admin.TapLogger"
	log.Print(OP)

	md, _ := metadata.FromIncomingContext(ctx)
	event := Event{}
	if val, ok := md["consumer"]; ok {
		event.Consumer = val[0]
	}

	event.Host = md[":authority"][0]
	event.Method = info.FullMethodName
	event.Timestamp = time.Now().Unix()

	fmt.Println(event)

	a.LogNotifier.NewNotify(&event)
	// TODO: Здесь же обработать для статистики

	return ctx, nil
}
