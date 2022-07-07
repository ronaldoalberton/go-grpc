package main

import (
	"context"
	"fmt"
	"github.com/ronaldoalberton/fc2-grpc/pb/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {

	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to grpc Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)

	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "joao",
		Email: "juca@bala.com",
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
		Name:  "joao",
		Email: "juca@bala.com",
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
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())

	}

}

func AddUsers(client pb.UserServiceClient) {

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Juca",
			Email: "bala@email.com",
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
		log.Fatalf("Error reiceve response: %v", err)
	}

	fmt.Println("Response: ", res)

}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Juca",
			Email: "bala@email.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Juca",
			Email: "bala@email.com",
		},
	}

	wait := make(chan int)

	go func() {

		for _, req := range reqs {

			fmt.Println("Sending user ", req.Name)

			stream.Send(req)
			time.Sleep(time.Second * 3)

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
				break
			}

			fmt.Println("Recebendo user %v com status: %v", res.GetUser(), res.Status)
		}
		close(wait)
	}()

	<-wait

}
