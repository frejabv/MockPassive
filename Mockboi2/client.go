package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"MockPassive/Mockboi2/protobuf"

	"google.golang.org/grpc"
)

func main() {
	LOG_FILE := "../log.txt"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	conn, err := grpc.Dial(":8070", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := protobuf.NewMockClient(conn)

	go TakeInput(client)
	fmt.Println("Press enter to increment, results will be written to log file")
	time.Sleep(1000 * time.Second)

}

func TakeInput(client protobuf.MockClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input != "\n" {
			log.Fatal("User typed something else")
		}

		message, error := client.Increment(context.Background(), &protobuf.IncrementRequest{})
		if error != nil {
			log.Fatalln("Couldn't increment - something went wrong")
		}

		log.Println("The new value is:", message.NewValue)
	}
}
