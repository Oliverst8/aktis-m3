package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	proto "main/grpc"
	"net"
	"sync"

	"google.golang.org/grpc"
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

	fmt.Printf("Now listening on localhost%s\n", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (s *ChittyChat) Join(client *proto.Client, stream proto.ChittyChat_JoinServer) error {
	s.updateClock(client.Count)
	s.increaseClock()
	log.Printf("%s is trying to join at Lamport time %d \n", client.Name, s.count)
	if s.streams[client.Name] != nil {
		response := proto.Response{
			Err: fmt.Sprintf("user with name \"%s\" already exists", client.Name),
		}
		err := s.SendMessage(&response, stream)
		if err != nil {
			return err
		}
		return errors.New("test error")
	}
	s.setStream(client.Name, stream)
	joinMessage := fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", client.Name, s.count)
	message := proto.Response{
		Text:   joinMessage,
		Client: client.Name,
		Count:  s.count,
	}
	s.increaseClock()
	err := s.Broadcast(&message)
	if err != nil {
		return err
	}
	log.Printf("%s has successfully joined the chat\n", client.Name)
	for s.getStream(client.Name) != nil {
	}
	return nil
}

func (s *ChittyChat) Leave(ctx context.Context, client *proto.Client) (*proto.Empty, error) {
	s.updateClock(client.Count)
	s.increaseClock()
	broadcastMessage := proto.Response{
		Text:   fmt.Sprintf("%s has left the chat at lamport time %d.", client.Name, s.count),
		Client: client.Name,
		Count:  s.count,
	}
	err := s.Broadcast(&broadcastMessage)
	if err != nil {
		return nil, err
	}
	s.streams[client.Name] = nil
	return &proto.Empty{}, nil
}

func (s *ChittyChat) PublishMessage(ctx context.Context, message *proto.Response) (*proto.Status, error) {
	s.updateClock(message.Count)
	s.increaseClock()
	if len(message.Text) > 128 {
		return &proto.Status{Success: false}, nil
	}
	log.Printf("User %s publishes message \"%s\" at client Lamport time: %d, received at server lamport time: %d", message.Client, message.Text, message.Count, s.count)
	err := s.Broadcast(message)
	if err != nil {
		return &proto.Status{Success: false}, err
	}
	return &proto.Status{Success: true}, nil
}

func (s *ChittyChat) Broadcast(message *proto.Response) error {
	log.Printf("Broadcasting message at Lamport time %d: %s", s.count, message.Text)
	lock.Lock()
	for _, stream := range s.streams {
		if stream == nil {
			continue
		}
		err := s.SendMessage(message, stream)
		if err != nil {
			panic(err)
		}
	}
	lock.Unlock()
	return nil
}

func (s *ChittyChat) SendMessage(message *proto.Response, stream proto.ChittyChat_JoinServer) error {
	s.increaseClock()
	message.Count = s.count
	log.Printf("sent message: %s at lamport time %d", message.Text, s.count)
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

func (s *ChittyChat) updateClock(reponseTime uint64) {
	if s.count < reponseTime {
		s.count = reponseTime
	}
}

func (s *ChittyChat) increaseClock() {
	s.count++
}
