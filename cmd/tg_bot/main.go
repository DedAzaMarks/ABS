package main

import (
	"log"
	"os"

	"github.com/DedAzaMarks/ABS/internal/server"
	"github.com/DedAzaMarks/ABS/internal/storage"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	log.SetPrefix("tg_bot: ")
	log.SetFlags(log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error on reading .env file: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	db, _ := storage.GetRepo("inmemory")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	s := server.NewServer(db, bot, redisClient)
	s.Start()
}
