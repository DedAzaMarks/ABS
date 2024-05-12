package server

import (
	"context"
	"errors"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/DedAzaMarks/ABS/internal/storage"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
)

const (
	ping = "ping"
	pong = "pong"

	start  = "start"
	open   = "open"
	close_ = "close"
	abort  = "abort"
	help   = "help"
	ID     = "ID"
)

var homeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ping, ping),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ID, ID),
	),
)

type Server struct {
	bot   *tgbotapi.BotAPI
	repo  storage.Repo
	redis *redis.Client
}

func NewServer(repo storage.Repo, bot *tgbotapi.BotAPI, redis *redis.Client) *Server {
	return &Server{
		repo:  repo,
		bot:   bot,
		redis: redis,
	}
}

func (s *Server) Start() {

	go s.registerNewClients()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)
	for update := range updates {
		log.Print("recv new message")
		if update.Message != nil {
			userID := strconv.FormatInt(update.Message.From.ID, 10)
			if update.Message.IsCommand() {
				log.Print("it is command")
				switch update.Message.Command() {
				case start:
					log.Print("start from", update.Message.From.ID)
					if err := s.repo.AddNewUser(userID); !errors.Is(err, storage.ErrorUserAlreadyExists) {
						msg := tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"New user registered. Your ID: "+
								userID+
								". Use it to connect your devices to system.")
						if _, err := s.bot.Send(msg); err != nil {
							log.Print(err)
						}
					}
					msg := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Choose action",
					)
					msg.ReplyMarkup = homeKeyboard
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
				}
			}
		} else if update.CallbackQuery != nil {
			usedID := strconv.FormatInt(update.CallbackQuery.From.ID, 10)
			callback := tgbotapi.NewCallback(usedID, update.CallbackQuery.Data)
			if _, err := s.bot.Request(callback); err != nil {
				log.Print(err)
			}
			switch data := update.CallbackQuery.Data; data {
			case ping:
				log.Print("publish: " + "ping:" + usedID)
				if cmd := s.redis.Publish(context.Background(), "ping:"+usedID, "ping"); cmd.Err() != nil {
					log.Print(cmd.Err())
				}
			case ID:
				log.Print("publish: " + usedID)
				if cmd := s.redis.Publish(context.Background(), usedID, usedID); cmd.Err() != nil {
					log.Print(cmd.Err())
				}
			}
		}
	}
}

func (s *Server) registerNewClients() {
	registerSub := s.redis.Subscribe(context.Background(), "register_new_client")
	log.Print("subscribed to register_new_client")
	defer registerSub.Close()
	for {
		register, err := registerSub.ReceiveMessage(context.Background())
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print("got register message")
		var registerReq pb.RegisterNewClient
		if err := proto.Unmarshal([]byte(register.Payload), &registerReq); err != nil {
			log.Print(err)
			continue
		}
		log.Print("register message unmarshalled")
		if err := s.repo.AddNewClient(registerReq.UserID, registerReq.ClientID); err != nil {
			log.Print(err)
			continue
		}
		log.Print("new client registered")
	}
}
