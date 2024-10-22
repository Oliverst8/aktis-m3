package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "main/grpc"
	"os"
)

func main() {
	//Following line is 'boilerplate', should be repeated mindlessly
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewChittyChatClient(conn)
	if len(os.Args) <= 1 {
		panic("You need to give a name for your client")
	}
	name := os.Args[1]
	me := proto.Client{
		Name:  name,
		Count: 0,
	}

	stream, err := client.Join(context.Background(), &me)

	/*var myErr *Errors.DuplicateNameError
	if errors.As(err, &myErr) {
		panic(err)
	}

	if errors.Is(err, &Errors.DuplicateNameError{}) {
		fmt.Println("DuplicanameError")
		panic(err)
		return
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Println("No error found when trying to join")
	}*/

	go listingForMessage(stream)
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
		message := proto.Response{
			Text:   text,
			Client: me.Name,
		}
		_, err := client.PublishMessage(context.Background(), &message)
		if err != nil {
			panic(err)
		}
	}
}

func listingForMessage(stream grpc.ServerStreamingClient[proto.Response]) {
	for {
		response, _ := stream.Recv()

		if response == nil {
			continue
		}
		if response.Err != "" {
			panic(response.Err)
		}

		fmt.Printf("%d - %s:\n", response.Count, response.Client)
		fmt.Println(response.Text)

	}
}
