package bot

import (
	"fmt"
	// "time"
	// "strings"

	"github.com/ManaGor/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


var botStr string = `
	<b><u>Welcome to the ManaGor platform!</u></b>

	This is an all-in-one life management platform that you can interact with from anywhere using Telegram.
	This is the main menu of the platform - you can access the different applications from here.

	The platform serves the following applications:
		- <b>ToDo:</b> add, list, search, filter, delete Todos. You can get Telegram notifications via setting the ReminderAt property of the Todo.
		- <b>Calendar:</b> 
		- <b>Warranty Tracker:</> save & track your warranties, get notifications when something expires, so that you could get rid of it.
		- <b>Recipes:</b> create a database of your favourite recipes.
		
	To access these applications, use the following commands:
		/todo
		/cal
	`

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

	go SendReminder(runCtx, bot)


	for update := range updates {
		if update.Message == nil { continue }

		if !allowedUsers[update.Message.From.ID] {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unauthorized"))
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, botStr)
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)

		case "todo":
			err := TodoBot(runCtx, bot, updates)
			if err != nil {
				return err
			}
		case "cal":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Calendar was called."))

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command - use /start to see available commands.")
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		}
	}

	return nil
}
