package main

func NewAdmin() *Admin {
	return &Admin{}
}

type Admin struct {
	UnimplementedAdminServer
}

func (a *Admin) Logging(nothing *Nothing, server Admin_LoggingServer) error {

	return nil
}

func (a *Admin) Statistics(interval *StatInterval, server Admin_StatisticsServer) error {

	return nil
}
