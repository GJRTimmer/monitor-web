package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/monitor-srv/proto/monitor"
	"github.com/micro/monitor-web/handler"
	"github.com/pborman/uuid"
)

func serve(r http.Handler) {
	l, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go http.Serve(l, r)

	host, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(host) == 0 || host == "[::]" || host == "::" {
		host, _ = os.Hostname()
	}

	service := &registry.Service{
		Name:    "go.micro.web.monitor",
		Version: "latest",
		Nodes: []*registry.Node{&registry.Node{
			Id:      uuid.NewUUID().String(),
			Address: host,
			Port:    port,
		}},
	}

	if err := registry.Register(service); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Listening on %v\n", l.Addr().String())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Printf("Received signal %s\n", <-ch)

	if err := registry.Deregister(service); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := l.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	cmd.Init()

	handler.MonitorClient = monitor.NewMonitorClient("go.micro.srv.monitor", client.DefaultClient)

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
	serve(r)
}
