package main

import (
	"fmt"
	"os"
	 "database/sql"
    _ "modernc.org/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Hell World!")

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found.")
		os.Exit(1)
    }
    
    token := os.Getenv("TELEGRAM_TOKEN")
    if token == "" {
		fmt.Println("No .env file found.")
		os.Exit(1)
    }


	db, err := sql.Open("sqlite", "./familyhub.db")
	if err != nil {
		fmt.Printf("DB Validation failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Error pinging DB: %v\n", err)
	}
	
	fmt.Println("DB is connected.")

	var allowedUsers = map[int64]bool{
		7236307955: true,  // your Telegram user ID
		// 987654321: true,  // spouse's ID
	}
	bot, _ := tgbotapi.NewBotAPI(token)
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	for update := range updates {
		if update.Message == nil { continue }

		if !allowedUsers[update.Message.From.ID] {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unauthorized"))
			continue
		}
		
		switch update.Message.Command() {
		case "todos":
			fmt.Println("todos command was called on Telegram")
			reply := "todos comand was received - this is the reply from the Go Application"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
		case "addtodo":
			text := update.Message.CommandArguments()
			fmt.Printf("addtodo was received with arguments: %s\n", text)

			reply := "addtodo comand was received - this is the reply from the Go Application"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
		}
	}
}
