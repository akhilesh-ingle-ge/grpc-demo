package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "grpc-prac/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = ":8080"

// Unary RPC
func callSayHello(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("%v", res.Message)
}

// Server Streaming RPC
func callSayHelloServerStreaming(client pb.GreetServiceClient, names *pb.ListNames) {
	log.Println("==========================================")
	log.Println("Server Streaming started")
	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v\n", err)
		}
		log.Println(message)
	}
	log.Println("Server Streaming finished")
	log.Println("==========================================")
}

// Client Streaming RPC
func callSayHelloClientStreaming(client pb.GreetServiceClient, names *pb.ListNames) {
	log.Println("==========================================")
	log.Println("Client streaming started!!")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}
	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while sending: %v", err)
		}
		log.Printf("Sent request with name: %v", name)
		time.Sleep(time.Second * 2)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving: %v", err)
	}
	log.Println("Client streaming finished..")

	// M-1
	log.Printf("%v", res.Messages)
	log.Println("==========================================")

	// M-2
	// for _, message := range res.Messages {
	// 	log.Println(message)
	// }
}

func callSayHelloBidirectionalstreaming(client pb.GreetServiceClient, names *pb.ListNames) {
	log.Println("==========================================")
	log.Printf("Bidirectional Streaming started")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming finished")
	log.Println("==========================================")
}

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.ListNames{
		Names: []string{"Akhilesh", "Sourav", "Priyansh"},
	}

	// Unary RPC
	callSayHello(client)
	// Server-side streaming RPC
	callSayHelloServerStreaming(client, names)
	// Client-side streaming RPC
	callSayHelloClientStreaming(client, names)
	// Bi-directional streaming RPC
	callSayHelloBidirectionalstreaming(client, names)
}
