package main

import (
	"context"
	"flag"
	"log"
	tgClient "telego/clients/telegram"
	event_consumer "telego/consumer/event-consumer"
	"telego/events/telegram"
	"telego/storage/sqllite"
)

const (
	tgBotHost      = "api.telegram.org"
	sqlStoragePath = "data/sqllite/storage.db"
	batchSize      = 100
)

// 7121076448:AAEbwWDK-DyXxfPfhM3ilGD5xND5d_ESKrs
func main() {
	//s:= files.New(storagePath)
	s, err := sqllite.New(sqlStoragePath)
	if err != nil {
		log.Fatalf("can`t connect to database: %v", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can`t initialize database: %v", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("starting telegram bot")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot")

	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
