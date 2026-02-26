package bot

import (
	"fmt"
	"strings"

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
		
		switch {
		case strings.Contains(update.Message.Command(), "todo"):
			err := TodoBot(runCtx, bot, update) 
			if err != nil {
				return err
			}
		case strings.Contains(update.Message.Command(), "cal"):
			reply := "Calendar command was called.\n"	
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
		}
	}

	return nil
}
