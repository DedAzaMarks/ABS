package main

import (
	"context"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix("client: ")
	log.SetFlags(log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error on reading .env file: %v", err)
	}

	host := os.Getenv("CLIENT_SERVER_HOST")
	port := os.Getenv("CLIENT_SERVER_PORT")

	conn, err := grpc.Dial(
		host+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	c := pb.NewClientClient(conn)
	pingStream, err := c.Ping(context.Background(), &pb.Init{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := pingStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error on recv: %v", err)
		}
		log.Println("ping")
	}
}
