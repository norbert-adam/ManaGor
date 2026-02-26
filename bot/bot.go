package bot

import (
	"fmt"

	"github.com/ManaGor/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(runCtx *models.RunCtx) error {
	var allowedUsers = map[int64]bool{
		7236307955: true,  // your Telegram user ID
		// 987654321: true,  // spouse's ID
	}
	bot, err := tgbotapi.NewBotAPI(runCtx.TgToken)
	if err != nil {
		return fmt.Errorf("error creating Telegram Bot API: %v", err)
	}
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

	return nil
}
