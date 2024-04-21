package main

import (
	"github.com/DedAzaMarks/ABS/internal/client"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ping = "ping"
	pong = "pong"

	open    = "open"
	connect = "connect"
	close_  = "close"
	abort   = "abort"
	help    = "help"
)

var homeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(connect, connect),
		tgbotapi.NewInlineKeyboardButtonData(ping, ping),
	),
)

type state = int

const (
	EMPTY state = iota
)

type TGUser struct {
	state      state
	deviceHost string
	devicePort string
}

type TGBotServer struct {
	mu    sync.RWMutex
	users map[int64]*TGUser
}

func main() {
	log.SetPrefix("tg_bot: ")
	log.SetFlags(log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error on reading .env file: %v", err)
	}

	socket, err := net.Listen("tcp", ":"+os.Getenv("CLIENT_SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	server := client.NewServer()
	pb.RegisterClientServer(grpcServer, server)
	log.Printf("Listening at %s", socket.Addr().String())
	go func() {
		log.Fatal(grpcServer.Serve(socket))
	}()

	botServer := TGBotServer{users: make(map[int64]*TGUser)}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.CallbackQuery != nil {

			botServer.mu.Lock()
			if _, ok := botServer.users[update.CallbackQuery.From.ID]; !ok {
				botServer.users[update.CallbackQuery.From.ID] = &TGUser{}
			}
			botServer.mu.Unlock()

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Print(err)
			}

			switch data := update.CallbackQuery.Data; data {
			case connect:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "deprecated or preprecated")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
			case ping:
				server.DoPing <- struct{}{}
			}
		}
	}
	close(server.DoPing)
}
