package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"MockPassive/Mockboi2/protobuf"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedMockServer
}

var client, client1, client2 protobuf.MockClient
var currentLeader int

func main() {
	LOG_FILE := "../log.txt"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	go startServer()

	//Start client(s)
	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}

	conn1, err1 := grpc.Dial(":8081", grpc.WithInsecure(), grpc.WithBlock())
	if err1 != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err1)
	}

	conn2, err2 := grpc.Dial(":8082", grpc.WithInsecure(), grpc.WithBlock())
	if err2 != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err2)
	}

	defer conn.Close()
	defer conn1.Close()
	defer conn2.Close()

	client = protobuf.NewMockClient(conn)
	client1 = protobuf.NewMockClient(conn1)
	client2 = protobuf.NewMockClient(conn2)

	currentLeader = 2

	fmt.Println("Frontend is running")
	//go TakeInput(client, client1, client2)
	time.Sleep(1000 * time.Second)
}

func startServer() {
	//Start server
	lis, err := net.Listen("tcp", ":8070")

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterMockServer(s, &server{})

	fmt.Println("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Increment(ctx context.Context, in *protobuf.IncrementRequest) (*protobuf.IncrementReply, error) {
	if currentLeader == 0 {
		fmt.Println("Currentleader is 0")
		message, error := client.Increment(context.Background(), &protobuf.IncrementRequest{})
		if error == nil {
			return &protobuf.IncrementReply{NewValue: message.NewValue}, nil
		}
	} else if currentLeader == 1 {
		fmt.Println("Currentleader is 1")
		message1, error1 := client1.Increment(context.Background(), &protobuf.IncrementRequest{})
		if error1 == nil {
			return &protobuf.IncrementReply{NewValue: message1.NewValue}, nil
		}
	} else if currentLeader == 2 {
		fmt.Println("Currentleader is 2")
		message2, error2 := client2.Increment(context.Background(), &protobuf.IncrementRequest{})
		if error2 == nil {
			return &protobuf.IncrementReply{NewValue: message2.NewValue}, nil
		}
	}
	return &protobuf.IncrementReply{NewValue: 0}, errors.New("No current leader")
}
