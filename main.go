package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ahfrd/grpc/auth-client/config"
	"github.com/ahfrd/grpc/auth-client/src/proto/auth"
	service "github.com/ahfrd/grpc/auth-client/src/service"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	// jwt := helpers.JwtWrapper{
	// 	SecretKey:       c.JWTSecretKey,
	// 	Issuer:          "micro-auth-grpc",
	// 	ExpirationHours: 24 * 365,
	// }
	fmt.Println(c.Port)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := service.ControllerAuth{}

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
