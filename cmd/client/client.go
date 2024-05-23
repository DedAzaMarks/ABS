package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/martinlindhe/inputbox"
	"google.golang.org/protobuf/proto"
)

var ClientID = uuid.New()

//go:embed embed_dowload_finished.sh
var downloadFinishedScrpitSrc []byte

func main() {
	log.SetPrefix("client: ")
	log.SetFlags(log.Lshortfile)

	log.Println(downloadFinishedScrpitSrc)
	pid := os.Getpid()
	log.Println("my pid:", pid)
	if err := os.WriteFile("/tmp/download_finished.sh", strconv.AppendInt(downloadFinishedScrpitSrc, int64(pid), 10), 0777); err != nil {
		log.Fatal("failed to remember pid")
	}

	redisAddr := flag.String("redis", "localhost:6379", "Redis address")
	flag.Parse()

	userID, ok := inputbox.InputBox("Remote Download", "Type user ID", "")
	if !ok {
		log.Fatal("host:port was not set")
	}

	ctx := context.Background()

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
	if cmd := redisClient.Publish(ctx, "register_new_client", buf); cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}
	log.Print("register message sent")
	pubSub := redisClient.Subscribe(ctx, userID)
	defer pubSub.Close()
	go func() {
		pingSub := redisClient.Subscribe(ctx, "ping:"+userID)
		for {
			log.Print("listening for ping")
			ping, err := pingSub.ReceiveMessage(ctx)
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
			msg, err := pubSub.ReceiveMessage(ctx)
			if err != nil {
				log.Print("Error receiving message: ", err)
				continue
			}
			log.Print("received message")
			if err := beeep.Alert("Remote Download Client", msg.Payload, ""); err != nil {
				log.Fatal(err)
			}
			cmd := exec.Command(
				`transmission-cli`,
				msg.Payload,
				`-w`, `/tmp`,
				`-f`, `/tmp/download_finished.sh`)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Print(err)
			}
			// Создаем канал для приема сигналов
			sigChan := make(chan os.Signal, 1)

			// Определяем, какие сигналы будем обрабатывать
			signal.Notify(sigChan, syscall.SIGUSR1)

			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			// Блокируем основной поток до тех пор, пока не получим сигнал

		loop:
			for {
				select {
				case sig := <-sigChan:
					switch sig {
					case syscall.SIGUSR1:
						fmt.Println("SIGUSR1 пойман за руку как дешевка!")
						if err := cmd.Process.Kill(); err != nil {
							log.Println(err)
						}
					}
					break loop
				case <-ticker.C:
					log.Print("идет загрузка...")
				}
			}

			// Здесь можно выполнить любую необходимую очистку
			fmt.Println("Загрузка завершена!")
		}
	}
}
