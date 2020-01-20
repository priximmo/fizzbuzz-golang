package common

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"log"
	"fmt"
)

// Define the GRPC Server Struct
type GrpcServer struct {
	Ctx context.Context
	host string
	port int
	Server *grpc.Server
	State FizzbuzzServerState
	listener net.Listener
	services []ServiceInterface
}

func NewGrpcServer(ctx context.Context, host string, port int) *GrpcServer {
	return &GrpcServer{
		Ctx: ctx,
		host: host,
		port: port,
		Server: grpc.NewServer(),
		State: Init,
		listener: nil,
		services: make([]ServiceInterface, 0),
	};
}

func (o *GrpcServer) Listen() (err error) {
	uri := fmt.Sprintf("%s:%d", o.host, o.port)
	o.listener, err = net.Listen("tcp", uri)
	if err != nil {
		o.State = Error
		log.Fatalf("failed to listen: %v", err)
	}
	o.State = Listen
	log.Printf("[GRPC] services started, listen on %s\n", uri)
	return err
}

func (o *GrpcServer) AddService(service ServiceInterface) {
	o.services = append(o.services, service)
}

func (o *GrpcServer) startServices() {
	for _, service := range o.services {
		service.RegisterGrpc(o)
	}
}

func (o *GrpcServer) Start() {
	o.Listen()
	o.startServices()
	go o.Server.Serve(o.listener)
	o.State = Ready
}

func (o *GrpcServer) GracefulStop() {
	if o.isGracefulStopable() { 
		o.Server.GracefulStop()
	}
}


func (o *GrpcServer) isGracefulStopable() bool {
	switch (o.State) {
	case
		Ready,
		Listen:
		return true
	}
	return false
}