package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"

	"MockPassive/Mockboi2/protobuf"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedMockServer
}

var value int32 = -1

func main() {
	log.Print("Welcome Server. You need to write 0, 1 or 2:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	port := strings.Replace(text, "\n", "", 1)

	lis, err := net.Listen("tcp", ":808"+port)

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterMockServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Increment(ctx context.Context, in *protobuf.IncrementRequest) (*protobuf.IncrementReply, error) {
	log.Println("Server received increment")
	value += 1
	return &protobuf.IncrementReply{NewValue: value}, nil
}

func (s *server) SetValue(ctx context.Context, in *protobuf.SetValueRequest) (*protobuf.SetValueReply, error) {
	log.Println("Server received set value request")
	value = in.Value
	return &protobuf.SetValueReply{}, nil
}
