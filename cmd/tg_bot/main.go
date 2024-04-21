package main

import (
	"github.com/DedAzaMarks/ABS/internal/client"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
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
	EMPTY   state = iota
	CONNECT       = iota
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

	socket, err := net.Listen("tcp", os.Getenv("CLIENT_SERVER_HOST")+":"+os.Getenv("CLIENT_SERVER_PORT"))
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
		// Check if we've gotten a message update.
		if update.Message != nil {
			botServer.mu.Lock()
			if _, ok := botServer.users[update.Message.From.ID]; !ok {
				botServer.users[update.Message.From.ID] = &TGUser{}
			}
			botServer.mu.Unlock()

			botServer.mu.Lock()
			if botServer.users[update.Message.From.ID].state == CONNECT {
				ip, err := net.ResolveTCPAddr("tcp", update.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error on parsing address")
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
						botServer.mu.Unlock()
						continue
					}
				}
				botServer.users[update.Message.From.ID].deviceHost = string(ip.IP)
				botServer.users[update.Message.From.ID].devicePort = strconv.Itoa(ip.Port)
				botServer.users[update.Message.From.ID].state = EMPTY
				botServer.mu.Unlock()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(ip.IP)+":"+strconv.Itoa(ip.Port)+" saved")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
				continue
			} else if !update.Message.IsCommand() {
				botServer.mu.Unlock()
				continue
			}
			botServer.mu.Unlock()
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "select action")
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case open:
				msg.Text = "select action"
				msg.ReplyMarkup = homeKeyboard
			case close_:
				msg.Text = "close keyboard"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			case abort:
				botServer.mu.Lock()
				botServer.users[update.Message.From.ID].state = EMPTY
				botServer.mu.Unlock()
				msg.Text = "action was aborted"
			case help:
				msg.Text = "Supported commands are: /open , /close , /help"
			}
			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				log.Println(err)
			}
		} else if update.CallbackQuery != nil {

			botServer.mu.Lock()
			if _, ok := botServer.users[update.CallbackQuery.From.ID]; !ok {
				botServer.users[update.CallbackQuery.From.ID] = &TGUser{}
			}
			botServer.mu.Unlock()

			botServer.mu.RLock()
			if botServer.users[update.CallbackQuery.From.ID].state == CONNECT {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "input host:port of your client device\nor use /abort to cancel")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
				botServer.mu.RUnlock()
				continue
			}
			botServer.mu.RUnlock()

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Print(err)
			}

			switch data := update.CallbackQuery.Data; data {
			case connect:
				botServer.mu.Lock()
				botServer.users[update.CallbackQuery.Message.Chat.ID].state = CONNECT
				botServer.mu.Unlock()
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "input host and port of your client device\nor use /abort to cancel")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}
			case ping:
				botServer.mu.RLock()
				host := botServer.users[update.CallbackQuery.From.ID].deviceHost
				port := botServer.users[update.CallbackQuery.From.ID].devicePort
				botServer.mu.RUnlock()
				if host == "" || port == "" {
					log.Println("host or port are empty")
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "host or port are empty")
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
					}
					continue
				}
				if host == "" {
					host = "localhost"
				}
				server.DoPing <- struct{}{}
				log.Printf("ping %s:%s", host, port)

			}
		}
	}
	close(server.DoPing)
}
