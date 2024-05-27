package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/DedAzaMarks/ABS/internal/server/storage"
	"github.com/DedAzaMarks/ABS/internal/server/storage/cache"
	"github.com/DedAzaMarks/ABS/internal/server/storage/persistent"
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
	db, err := storage.NewCachedStorage(ctx, persistent.Postgres, cache.Redis, redisClient.(*redis.Client))
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(db, bot, redisClient.(*redis.Client))
	s.Start()
}
