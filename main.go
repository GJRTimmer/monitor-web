package main

import (
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"github.com/micro/monitor-web/handler"

	monitor "github.com/micro/monitor-srv/proto/monitor"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/healthchecks", handler.Healthchecks)
	r.HandleFunc("/healthchecks/{id}", handler.Healthcheck)
	r.HandleFunc("/stats", handler.Stats)
	r.HandleFunc("/stats/{service}", handler.ServiceStats)
	r.HandleFunc("/status", handler.Status)
	r.HandleFunc("/status/{service}", handler.ServiceStatus)
	r.HandleFunc("/services", handler.Index)
	r.HandleFunc("/services/{service}", handler.Service)

	service := web.NewService(
		web.Name("go.micro.web.monitor"),
		web.Handler(r),
	)

	service.Init()
	handler.MonitorClient = monitor.NewMonitorClient("go.micro.srv.monitor", client.DefaultClient)
	service.Run()
}
