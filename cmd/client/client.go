package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/DedAzaMarks/ABS/internal/proto"

	"github.com/gen2brain/beeep"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/martinlindhe/inputbox"
	"google.golang.org/protobuf/proto"
)

var ClientID = uuid.New()

func main() {
	log.SetPrefix("client: ")
	log.SetFlags(log.Lshortfile)

	redisAddr := flag.String("redis", "localhost:6379", "Redis address")
	flag.Parse()

	userID, ok := inputbox.InputBox("Remote Download", "Type user ID", "")
	if !ok {
		log.Fatal("host:port was not set")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer redisClient.Close()
	registerReq := pb.RegisterNewClient{
		UserID:   userID,
		ClientID: ClientID.String(),
	}
	buf, err := proto.Marshal(&registerReq)
	if err != nil {
		log.Fatal(err)
	}
	if cmd := redisClient.Publish(context.Background(), "register_new_client", buf); cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}
	log.Print("register message sent")
	pubSub := redisClient.Subscribe(context.Background(), userID)
	defer pubSub.Close()
	go func() {
		pingSub := redisClient.Subscribe(context.Background(), "ping:"+userID)
		for {
			log.Print("listening for ping")
			ping, err := pingSub.ReceiveMessage(context.Background())
			log.Print("ping received")
			if err != nil {
				log.Println(err)
				continue
			}
			log.Print(ping.Payload)
			if err := beeep.Alert("Remote Download Client", ping.Payload, ""); err != nil {
				log.Fatal(err)
			}
		}
	}()
	for {
		select {
		default:
			log.Print("listening for messages")
			msg, err := pubSub.ReceiveMessage(context.Background())
			if err != nil {
				log.Print("Error receiving message: ", err)
				continue
			}
			log.Print("received message")
			if err := beeep.Alert("Remote Download Client", msg.Payload, ""); err != nil {
				log.Fatal(err)
			}
		}
	}
}
