package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/gen2brain/beeep"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/martinlindhe/inputbox"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var ClientID = uuid.New()

//go:embed embed_dowload_finished.sh
var downloadFinishedScrpitSrc []byte

func main() {
	log.SetPrefix("client: ")
	log.SetFlags(log.Lshortfile)
	pid := os.Getpid()
	log.Println("my pid:", pid)
	f, err := os.CreateTemp("", "*_download_finished.sh")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error on creating tmp script: %v", err)
		os.Exit(1)
	}
	defer func(name string) { _ = os.Remove(name) }(f.Name())
	if err := os.WriteFile(f.Name(), strconv.AppendInt(downloadFinishedScrpitSrc, int64(pid), 10), 0777); err != nil {
		log.Fatal("failed to remember pid")
	}

	ctx := context.Background()

	rootCertPool := x509.NewCertPool()
	redisCert, err := os.ReadFile("/Users/m.bordyugov/.redis/YandexInternalRootCA.crt")
	if err != nil {
		log.Fatal(err)
	}
	if !rootCertPool.AppendCertsFromPEM(redisCert) {
		log.Fatal("failed to append redis cert")
	}
	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		TLSConfig: &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		},
		Addrs:    []string{"c-c9q2pkmkffgciep0re8p.rw.mdb.yandexcloud.net:6380"},
		Password: "redisredis", // no password set
		DB:       0,            // use default DB
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	sessionKey, _ := inputbox.InputBox("Remote Download", "Введите Ключ", "")
	devidceName, _ := inputbox.InputBox("Remote Download", "Назовите устройство", "")
	registerReq := pb.RegisterNewClient{
		SessionKey: sessionKey,
		DeviceID:   ClientID.String(),
		DeviceName: devidceName,
	}
	buf, err := proto.Marshal(&registerReq)
	if err != nil {
		log.Fatal(err)
	}
	if cmd := redisClient.Publish(ctx, "register_new_client", buf); cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}
	log.Print("register message sent")
	pubSub := redisClient.Subscribe(ctx, sessionKey)
	defer pubSub.Close()
	var ctxDownload context.Context
	var cancelDownload context.CancelFunc
	for {
		select {
		default:
			log.Print("listening for messages")
			msg, err := pubSub.ReceiveMessage(ctxDownload)
			if err != nil {
				log.Print("Error receiving message: ", err)
				continue
			}
			log.Print("received message")
			handleMessage(&ctxDownload, &cancelDownload, msg.Payload)
		}
	}
}

func handleMessage(ctx *context.Context, cancel *context.CancelFunc, msg string) {
	message := pb.ServerToClientChannelMessage{}
	err := proto.Unmarshal([]byte(msg), &message)
	if err != nil {
		log.Print(err)
		return
	}
	switch action := message.Action.(type) {
	case *pb.ServerToClientChannelMessage_Start:
		if err := beeep.Alert("Remote Download Devices", "Новая загрузка", ""); err != nil {
			_ = beeep.Alert("Remote Download Devices", "Alert error", "")
			log.Print(err)
		}
		*ctx, *cancel = context.WithCancel(context.Background())
		go func() {
			torrentDownload(ctx, action.Start.Href)
			*ctx, *cancel = nil, nil
		}()
	case *pb.ServerToClientChannelMessage_Stop:
		if *ctx == nil && *cancel == nil {
			log.Print("Ничего не загружается")
			if err := beeep.Alert("Remote Download Devices", "Ничего не загружается", ""); err != nil {
				_ = beeep.Alert("Remote Download Devices", "Alert error", "")
				log.Print(err)
			}
			return
		}
		log.Print("stopping torrent")
		if err := beeep.Alert("Remote Download Devices", "Отмена загрузки", ""); err != nil {
			_ = beeep.Alert("Remote Download Devices", "Alert error", "")
			log.Print(err)
		}
		(*cancel)()
		log.Print("torrent stopped")
		*ctx, *cancel = nil, nil
	case *pb.ServerToClientChannelMessage_List:
	default:
		log.Print("unknown action")
		return
	}
}

func torrentDownload(ctx *context.Context, magnet string) {
	cmd := exec.Command(
		`transmission-cli`,
		magnet,
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
		case <-(*ctx).Done():
			fmt.Println("Загрузка прервана!")
			if err := cmd.Process.Kill(); err != nil {
				log.Println(err)
			}
			break loop
		case <-ticker.C:

			log.Print("идет загрузка...")
		}
	}

	// Здесь можно выполнить любую необходимую очистку
	fmt.Println("Загрузка завершена!")
}
