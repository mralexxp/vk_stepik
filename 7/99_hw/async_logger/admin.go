package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"
	"sync"
	"time"
)

type Log struct {
	LastEvent *Event
}

type Admin struct {
	M    sync.Mutex
	Logs Log
	UnimplementedAdminServer
}

func NewAdmin() *Admin {
	return &Admin{}
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {

	return nil
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {

	return nil
}

func (a *Admin) TapLogger(ctx context.Context, info *tap.Info) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	event := Event{}
	if val, ok := md["consumer"]; ok {
		event.Consumer = val[0]
	}

	event.Host = md[":authority"][0]
	event.Method = info.FullMethodName
	event.Timestamp = time.Now().Unix()

	a.M.Lock()
	a.Logs.LastEvent = &event
	a.M.Unlock()

	// TODO: Здесь же обработать для статистики

	return ctx, nil
}
