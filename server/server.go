package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	proto "main/grpc"
	"net"
	"sync"
)

var lock sync.Mutex

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
		panic(err)
	}

	proto.RegisterChittyChatServer(grpcServer, s)

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Now listening on localhost%s", port)
}

func (s *ChittyChat) Join(client *proto.Client, stream proto.ChittyChat_JoinServer) error {
	fmt.Printf("%s is trying to join\n", client.Name)
	if s.streams[client.Name] != nil {
		fmt.Printf("user with name %s already exists", client.Name)
		response := proto.Response{
			Err: fmt.Sprintf("user with name \"%s\" already exists", client.Name),
		}
		err := s.SendMessage(&response, stream)
		if err != nil {
			return err
		}
		return errors.New("Test error")
	}
	s.setStream(client.Name, stream)
	joinMessage := fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", client.Name, 0)
	message := proto.Response{
		Text:   joinMessage,
		Client: client.Name,
		Count:  client.Count,
	}
	err := s.Broadcast(&message)
	if err != nil {
		return err
	}
	fmt.Printf("%s has successfully joined the chat\n", client.Name)
	for s.getStream(client.Name) != nil {
	}
	return nil
}

func (s *ChittyChat) Leave(ctx context.Context, client *proto.Client) (*proto.Empty, error) {

	stream := s.streams[client.Name]

	youLeftMessage := proto.Response{
		Text:   "You left",
		Client: client.Name,
		Count:  0,
	}
	broadcastMessage := proto.Response{
		Text:   fmt.Sprintf("%s has left the chat.", client.Name),
		Client: client.Name,
		Count:  0,
	}
	err := s.SendMessage(&youLeftMessage, stream)
	if err != nil {
		return nil, err
	}
	err = s.Broadcast(&broadcastMessage)
	if err != nil {
		return nil, err
	}
	s.streams[client.Name] = nil
	return &proto.Empty{}, nil
}

func (s *ChittyChat) PublishMessage(ctx context.Context, message *proto.Response) (*proto.Status, error) {
	err := s.Broadcast(message)
	if err != nil {
		return &proto.Status{Success: false}, err
	}
	return &proto.Status{Success: true}, nil
}

func (s *ChittyChat) Broadcast(message *proto.Response) error {
	lock.Lock()
	for name, stream := range s.streams {
		if message.Client == name {
			continue
		}
		err := stream.Send(message)
		if err != nil {
			panic(err)
		}
	}
	lock.Unlock()
	return nil
}

func (s *ChittyChat) SendMessage(message *proto.Response, stream proto.ChittyChat_JoinServer) error {
	err := stream.Send(message)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *ChittyChat) getStream(name string) proto.ChittyChat_JoinServer {
	lock.Lock()
	defer lock.Unlock()
	return s.streams[name]
}

func (s *ChittyChat) setStream(name string, stream proto.ChittyChat_JoinServer) {
	lock.Lock()
	defer lock.Unlock()
	s.streams[name] = stream
}

func (me *ChittyChat) increaseClock() {
	me.count++
}
