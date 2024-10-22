package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "main/grpc"
	"os"
	"os/user"
)

var me proto.Client

func main() {
	//Following line is 'boilerplate', should be repeated mindlessly
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewChittyChatClient(conn)
	var name string
	if len(os.Args) <= 1 {
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		name = currentUser.Username
	} else {
		name = os.Args[1]
	}
	me = proto.Client{
		Name:  name,
		Count: 0,
	}

	stream, err := client.Join(context.Background(), &me)

	go listeningForMessage(stream)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "leave" {
			_, err = client.Leave(context.Background(), &me)
			if err != nil {
				panic(err)
			}
			break
		}
		increaseClock()
		message := proto.Response{
			Text:   text,
			Client: me.Name,
			Count:  me.Count,
		}
		_, err := client.PublishMessage(context.Background(), &message)
		if err != nil {
			panic(err)
		}
	}
}

func listeningForMessage(stream grpc.ServerStreamingClient[proto.Response]) {
	for {
		response, _ := stream.Recv()

		if response == nil {
			continue
		}
		if response.Err != "" {
			panic(response.Err)
		}

		if response.Count > me.Count {
			me.Count = response.Count
		}
		increaseClock()

		fmt.Printf("%d - %s:\n", response.Count, response.Client)
		fmt.Println(response.Text)
	}
}

func increaseClock() {
	me.Count++
}
