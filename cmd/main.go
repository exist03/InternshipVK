package main

import (
	"fsm/internal/handlers"
	"fsm/pkg/DB"
	"github.com/joho/godotenv"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalln(err)
	}
	bot.Use(middleware.AutoRespond())
	bGroup := bot.Group()
	storage := memory.NewStorage()
	defer storage.Close()
	manager := fsm.NewManager(bGroup, memory.NewStorage())
	db, err := DB.OpenDB("quest:quest@/VK")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("db is open")
	defer db.Close()
	handlers.InitHandlers(bGroup, db, manager)
	log.Println("Handlers configured")
	bot.Start()
}
