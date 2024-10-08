package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	proto "main/grpc"
	"net"
)

type ChittyChat struct {
	proto.UnimplementedChittyChatServer
	count   uint64
	streams map[string]proto.ChittyChat_JoinServer
}

func main() {
	server := &ChittyChat{
		count:   0,
		streams: make(map[string]proto.ChittyChat_JoinServer),
	}
	server.start_server()
}

func (s *ChittyChat) start_server() {
	grpcServer := grpc.NewServer()
	port := ":5050"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Could not listen on port: " + port)
	}

	proto.RegisterChittyChatServer(grpcServer, s)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Could not start the server..")
	}
	fmt.Printf("Now listening on localhost%s", port)
}

func (s *ChittyChat) Join(client *proto.Client, stream proto.ChittyChat_JoinServer) error {
	if s.streams[client.Name] != nil {
		return errors.New("Name already taken")
	}
	s.streams[client.Name] = stream
	joinMessage := fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", client.Name, 0)
	message := proto.Message{
		Text:   joinMessage,
		Client: client.Name,
		Count:  client.Count,
	}
	err := s.Broadcast(&message)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChittyChat) Leave(ctx context.Context, client *proto.Client) (*proto.Empty, error) {
	return &proto.Empty{}, errors.New("Can't leave")
}

func (s *ChittyChat) PublishMessage(ctx context.Context, message *proto.Message) (*proto.Status, error) {
	err := s.Broadcast(message)
	if err != nil {
		return &proto.Status{Success: false}, err
	}
	return &proto.Status{Success: true}, nil
}

func (s *ChittyChat) Broadcast(message *proto.Message) error {
	for _, stream := range s.streams {
		err := stream.Send(message)
		if err != nil {
			return err
		}
	}
	return nil
}
