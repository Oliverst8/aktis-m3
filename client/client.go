package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	proto "main/grpc"
)

func main() {
	//Following line is 'boilerplate', should be repeated mindlessly
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working ;)")
	}
	client := proto.NewChittyChatClient(conn)
	me := proto.Client{
		Name:  "Me",
		Count: 0,
	}
	stream, err := client.Join(context.Background(), &me)

	for {
		message, _ := stream.Recv()
		if message == nil {
			continue
		}
		fmt.Println(message.Text)

	}

}
