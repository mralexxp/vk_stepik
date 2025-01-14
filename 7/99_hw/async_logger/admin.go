package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"
)

func NewAdmin() *Admin {
	loggerChan := make(chan Event, 1)

	return &Admin{
		Logs: loggerChan,
	}
}

type Admin struct {
	Logs           chan Event
	RequestCounter map[string]int
	UnimplementedAdminServer
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {
	for {
		select {
		case <-server.Context().Done():
			return nil
		case e, ok := <-a.Logs:
			if !ok {
				return nil
			}

			err := server.Send(&e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {

	return nil
}

func (a *Admin) TapLogger(ctx context.Context, info *tap.Info) (context.Context, error) {
	go func(ctx context.Context, info *tap.Info) {
		md, _ := metadata.FromIncomingContext(ctx)
		msg := Event{}
		if val, ok := md["consumer"]; ok {
			msg.Consumer = val[0]
		}

		msg.Host = md[":authority"][0]
		msg.Method = info.FullMethodName

		a.Logs <- msg
	}(ctx, info)

	return ctx, nil
}
