package main

import (
	"context"
	"github.com/ShatteredRealms/Accounts/pkg/accountspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func main() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := accountspb.NewUserServiceClient(conn)
	//client2 := accountspb.NewAuthenticationServiceClient(conn)
	ctx := context.Background()

	log.Printf("Sending message...")

	//_, err = client2.Register(ctx, &accountspb.RegisterAccountMessage{
	//    Email:     "wil@example.com",
	//    Password:  "password",
	//    Username:  "wil",
	//    FirstName: "Wil",
	//    LastName:  "Simpson",
	//})

	resp, err := client.GetAll(ctx, &emptypb.Empty{})

	log.Printf("Message received:\n")
	for _, v := range resp.Users {
		log.Printf("    username: %s\n", v.Username)
	}
}
