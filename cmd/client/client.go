package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lrmiguel/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Leo",
			Email: "leo@leo.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Leo2",
			Email: "leo2@leo.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Leo3",
			Email: "leo3@leo.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Leo4",
			Email: "leo4@leo.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "Leo5",
			Email: "leo5@leo.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, _ := client.AddUserStreamBoth(context.Background())
	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Leo",
			Email: "leo@leo.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Leo2",
			Email: "leo2@leo.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Leo3",
			Email: "leo3@leo.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Leo4",
			Email: "leo4@leo.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "Leo5",
			Email: "leo5@leo.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
