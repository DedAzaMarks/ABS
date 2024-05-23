package main

import (
	"context"
	"github.com/DedAzaMarks/ABS/internal/server/storage"
	"log"
	"os"

	"github.com/DedAzaMarks/ABS/internal/server"
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

	ctx := context.Background()
	db, _ := storage.GetRepo(ctx, storage.InMemory)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(db, bot, redisClient)
	s.Start()
}
