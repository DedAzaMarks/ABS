package main

import (
	"context"
	"flag"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {
	log.SetPrefix("client: ")
	log.SetFlags(log.Lshortfile)

	var host, port string
	flag.StringVar(&host, "host", "localhost", "set host")
	flag.StringVar(&port, "port", "42069", "set port")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error on reading .env file: %v", err)
	}

	conn, err := grpc.Dial(
		host+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	c := pb.NewClientClient(conn)
	log.Print("send init")
	pingStream, err := c.Ping(context.Background(), &pb.Init{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		log.Println("wait for ping")
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
