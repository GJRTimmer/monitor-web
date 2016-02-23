package main

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"github.com/micro/monitor-web/handler"

	monitor "github.com/micro/monitor-srv/proto/monitor"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.monitor"),
		web.Handler(handler.Router()),
	)

	service.Init()

	handler.Init(
		"templates",
		monitor.NewMonitorClient("go.micro.srv.monitor", client.DefaultClient),
	)

	service.Run()
}
