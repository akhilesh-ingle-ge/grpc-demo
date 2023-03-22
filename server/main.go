package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "grpc-prac/proto"

	"google.golang.org/grpc"
)

const port = ":8080"

type helloServer struct {
	pb.GreetServiceServer
}

// Unary RPC
func (s *helloServer) SayHello(ctx context.Context, req *pb.NoParam) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "hello, how are you?",
	}, nil
}

// Server-side streaming RPC
func (s *helloServer) SayHelloServerStreaming(req *pb.ListNames, stream pb.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("got requests with names: %v", req.Names)
	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello, " + name + " how are you?",
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

// Client-side streaming RPC
func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessageList{
				Messages: messages,
			})
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name %v", req.Name)
		res := "Hello " + req.Name + " How are you?"
		messages = append(messages, res)
	}
}

func (s *helloServer) SayHelloBidirectionalStreaming(stream pb.GreetService_SayHelloBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)
		strres := "Hello " + req.Name + " How are you?"
		res := &pb.HelloResponse{
			Message: strres,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	log.Printf("Server started at %v", lis.Addr())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
