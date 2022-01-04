package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"MockPassive/Mockboi2/protobuf"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedMockServer
}

var client0, client1, client2 protobuf.MockClient
var value int32 = -1
var port string
var leader bool = false
var clients []protobuf.MockClient

func main() {
	log.Print("Welcome Server. You need to write 0, 1 or 2:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	port = strings.Replace(text, "\n", "", 1)

	go startServer(port)

	fmt.Println("Port is:", port)
	//Start client(s)
	if port == "0" {
		client1 = startClient("1")
		client2 = startClient("2")
		clients = append(clients, client1, client2)
	} else if port == "1" {
		client0 = startClient("0")
		client2 = startClient("2")
		clients = append(clients, client0, client2)
	} else if port == "2" {
		client0 = startClient("0")
		client1 = startClient("1")
		clients = append(clients, client0, client1)
		leader = true
		fmt.Println("Hello this is the leader speaking")
	}

	fmt.Println("Server is running")
	time.Sleep(1000 * time.Second)
}

func startClient(port string) protobuf.MockClient {
	conn, err := grpc.Dial(":808"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Printf("Connection to port 808%s was succesful", port)

	defer conn.Close()

	client := protobuf.NewMockClient(conn)
	return client
}

func startServer(port string) {
	//Start server
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
	fmt.Println("Increment got through to server")
	if leader {
		fmt.Println("Server is leader")
		value += 1
		fmt.Println("Value is:", value)
		response, err := client0.SetValue(context.Background(), &protobuf.SetValueRequest{Value: value})
		fmt.Println("response from 0", response.GetAck())
		if err != nil {
			fmt.Println("epic fail")
		}
		for _, cli := range clients {
			response, err := cli.SetValue(context.Background(), &protobuf.SetValueRequest{Value: value})
			fmt.Println(response.GetAck())
			if !response.GetAck() || err != nil {
				//one replica is down
				fmt.Println("A replica is down")
			}
		}

		/*var values []int32
		var message *protobuf.IncrementReply
		var err error
		for _, cli := range clients {
			message, err = cli.Increment(context.Background(), &protobuf.IncrementRequest{})
			if err == nil {
				values = append(values, message.NewValue)
			}
		}

		//Maybe this is not necessary?
		var highestValue int32
		for i := 0; i < len(values); i++ {
			if values[i] > highestValue {
				highestValue = values[i]
			}
		}

		//syncValues
		for _, cli := range clients {
			if err == nil && message.NewValue != highestValue {
				cli.SetValue(context.Background(), &protobuf.SetValueRequest{Value: highestValue})
			}
		}*/

	}
	return &protobuf.IncrementReply{NewValue: value}, nil
}

func (s *server) SetValue(ctx context.Context, in *protobuf.SetValueRequest) (*protobuf.SetValueReply, error) {
	log.Println("Server received set value request")
	fmt.Println("received setvalue")
	value = in.Value
	return &protobuf.SetValueReply{Ack: true}, nil
}
