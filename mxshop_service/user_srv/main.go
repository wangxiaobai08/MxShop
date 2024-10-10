package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"user_srv/handler"
	"user_srv/proto"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "IP地址")
	Port := flag.String("port", "50051", "端口号")
	flag.Parse()
	fmt.Println("ip:", *IP, "port:", *Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})
	lis, err := net.Listen("tcp", fmt.Sprint(*IP, ":", *Port))
	if err != nil {
		panic("failed to listen" + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to serve" + err.Error())
	}
}
